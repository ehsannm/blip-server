package music

import (
	"context"
	"git.ronaksoftware.com/blip/server/internal/flusher"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/pools"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"sort"
	"strings"
	"sync"
	"time"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type indexedSong struct {
	song  *Song
	score float64
}

var songIndexer = flusher.New(1000, 1, time.Millisecond, func(items []flusher.Entry) {
	b := songIndex.NewBatch()
	for _, item := range items {
		song := item.Key.(*Song)
		_ = b.Index(song.ID.Hex(), song)
	}
	err := songIndex.Batch(b)
	if err != nil {
		log.Warn("Error On Indexing Song", zap.Error(err))
	}
	for _, item := range items {
		item.Callback(nil)
	}
})

func updateLocalIndex(s *Song) {
	_ = songIndexer.EnterWithResult(s, nil)
}

func deleteFromLocalIndex(songID primitive.ObjectID) error {
	return songIndex.Delete(songID.Hex())
}

// SearchLocalIndex
func SearchLocalIndex(keyword string, result int) ([]indexedSong, error) {
	qs := make([]query.Query, 0, 4)

	terms := strings.Split(strings.ToLower(keyword), "+")
	for _, t := range terms {
		t := strings.Trim(t, "()")
		qs = append(qs, bleve.NewTermQuery(t))
	}
	searchRequest := bleve.NewSearchRequest(bleve.NewDisjunctionQuery(qs...))
	searchRequest.Explain = true
	searchRequest.Size = result
	searchRequest.SortBy([]string{"-_score"})
	res, err := songIndex.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	log.Debug("Local Index Search Finished",
		zap.Strings("Terms", terms),
		zap.Duration("Time", res.Took),
		zap.Uint64("Total", res.Total),
		zap.Float64("MaxScore", res.MaxScore),
	)
	log.Debug("Local Index Search Info",
		zap.Int("Total", res.Status.Total),
		zap.Int("Successful", res.Status.Successful),
		zap.Int("Failed", res.Status.Failed),
	)

	foundSongs := make([]indexedSong, 0, len(res.Hits))
	waitGroup := pools.AcquireWaitGroup()
	for _, hit := range res.Hits {
		if songID, err := primitive.ObjectIDFromHex(hit.ID); err == nil {
			waitGroup.Add(1)
			go func(songID primitive.ObjectID, score float64) {
				songX, _ := GetSongByID(songID)
				if songX != nil {
					foundSongs = append(foundSongs, indexedSong{
						song:  songX,
						score: score,
					})
				}
				waitGroup.Done()
			}(songID, hit.Score)
		}
	}
	waitGroup.Wait()
	sort.Slice(foundSongs, func(i, j int) bool {
		return foundSongs[i].score > foundSongs[j].score
	})
	return foundSongs, nil
}

type searchCtx struct {
	sync.RWMutex
	keyword    string
	cursorID   string
	ctx        context.Context
	cancelFunc context.CancelFunc
	resChan    <-chan *crawler.SearchResponse
	done       chan struct{}
	songChan   chan *Song
	sent       map[primitive.ObjectID]struct{}
}

func (ctx *searchCtx) job() {
	log.Debug("SearchCtx started", zap.String("CursorID", ctx.cursorID))
MainLoop:
	for r := range ctx.resChan {
		for _, foundSong := range r.Result {
			if ctx.ctx.Err() != nil {
				break MainLoop
			}
			if ce := log.Check(log.DebugLevel, "Crawler Found Song"); ce != nil {
				ce.Write(
					zap.String("Title", foundSong.Title),
					zap.String("Artist", foundSong.Artists),
					zap.String("Source", r.Source),
				)
			}
			foundSong.Artists = strings.TrimSpace(foundSong.Artists)
			foundSong.Title = strings.TrimSpace(foundSong.Title)
			if len(foundSong.Title) == 0 {
				continue
			}
			uniqueKey := GenerateUniqueKey(foundSong.Title, foundSong.Artists)
			songX, err := GetSongByUniqueKey(uniqueKey)
			if err != nil {
				songX = &Song{
					ID:             primitive.NilObjectID,
					UniqueKey:      uniqueKey,
					Title:          foundSong.Title,
					Genre:          foundSong.Genre,
					Lyrics:         foundSong.Lyrics,
					Artists:        foundSong.Artists,
					SongStoreID:    0,
					OriginCoverUrl: foundSong.CoverUrl,
					OriginSongUrl:  foundSong.SongUrl,
					Source:         r.Source,
				}
				_, err = SaveSong(songX)
				if err != nil {
					log.Warn("Error On Save Search Result",
						zap.Error(err),
						zap.String("Source", r.Source),
						zap.String("Title", foundSong.Title),
					)
					continue
				}
			} else if songX.SongStoreID == 0 {
				// If the song has not been downloaded from source yet, update the origin url
				songX.Artists = foundSong.Artists
				songX.Title = foundSong.Title
				songX.Genre = foundSong.Genre
				songX.Lyrics = foundSong.Lyrics
				songX.OriginSongUrl = foundSong.SongUrl
				songX.OriginCoverUrl = foundSong.CoverUrl
				songX.Source = r.Source
				_, err = SaveSong(songX)
				if err != nil {
					log.Warn("Error On Save Search Result",
						zap.Error(err),
						zap.String("Source", r.Source),
						zap.String("Title", foundSong.Title),
					)
					continue
				}
			}
			updateLocalIndex(songX)
			select {
			case ctx.songChan <- songX:
			default:
				log.Warn("SearchCtx Could not write to song channel")
			}
		}
	}
	close(ctx.songChan)
	ctx.done <- struct{}{}
	log.Debug("SearchCtx done", zap.String("CursorID", ctx.cursorID))
}

func (ctx *searchCtx) SongChan() <-chan *Song {
	return ctx.songChan
}

func (ctx *searchCtx) ShouldSend(songID primitive.ObjectID) bool {
	ctx.RLock()
	_, ok := ctx.sent[songID]
	ctx.RUnlock()
	if ok {
		return false
	}
	ctx.Lock()
	ctx.sent[songID] = struct{}{}
	ctx.Unlock()
	return true
}

// StartSearch creates a new context and send the request to all the crawlers and waits for them to finish. If a context
// with the same id exists, we first cancel the old one and create a new context with new 'keyword'
func StartSearch(cursorID string, keyword string) *searchCtx {
	ctx := getSearchCtx(cursorID)
	if ctx != nil {
		ctx.cancelFunc()

		// Wait for job to be finished
		<-ctx.done
	}

	ctx = &searchCtx{
		cursorID: cursorID,
		keyword:  keyword,
		done:     make(chan struct{}, 1),
		songChan: make(chan *Song, 100),
	}
	ctx.ctx, ctx.cancelFunc = context.WithCancel(context.Background())
	ctx.resChan = crawler.Search(ctx.ctx, keyword)
	saveSearchCtx(ctx)
	go ctx.job()
	return ctx
}

// ResumeSearch checks if a context with 'cursorID' has been already exists and return the song channel, otherwise it returns
// nil.
func ResumeSearch(cursorID string) *searchCtx {
	ctx := getSearchCtx(cursorID)
	if ctx == nil {
		return nil
	}

	return ctx
}

func saveSearchCtx(ctx *searchCtx) {
	searchContextsMtx.Lock()
	searchContexts[ctx.cursorID] = ctx
	searchContextsMtx.Unlock()
}
func getSearchCtx(cursorID string) *searchCtx {
	searchContextsMtx.RLock()
	ctx := searchContexts[cursorID]
	searchContextsMtx.RUnlock()
	if ctx != nil {
		select {
		case <-ctx.done:
			searchContextsMtx.Lock()
			delete(searchContexts, cursorID)
			searchContextsMtx.Unlock()
			return nil
		default:
		}
	}
	return ctx
}
