package music

import (
	"context"
	"git.ronaksoftware.com/blip/server/internal/flusher"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"strings"
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

var songIndexer = flusher.New(100, 1, time.Millisecond, func(items []flusher.Entry) {
	b := songIndex.NewBatch()
	for _, item := range items {
		song := item.Key.(*Song)
		if d, _ := songIndex.Document(song.ID.Hex()); d == nil {
			_ = b.Index(song.ID.Hex(), song)
		}
	}
	err := songIndex.Batch(b)
	if err != nil {
		log.Warn("Error On Indexing Song", zap.Error(err))
	}
})

// updateLocalIndex updates the local index which will be used by search handlers
func updateLocalIndex(s *Song) {
	songIndexer.Enter(s, nil)
}

// SearchLocalIndex
func SearchLocalIndex(keyword string) ([]primitive.ObjectID, error) {
	qs := make([]query.Query, 0, 4)
	for _, t := range strings.Fields(keyword) {
		qs = append(qs, bleve.NewTermQuery(t))
	}
	searchRequest := bleve.NewSearchRequest(bleve.NewDisjunctionQuery(qs...))
	searchRequest.SortBy([]string{"_score"})
	res, err := songIndex.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	songIDs := make([]primitive.ObjectID, 0, len(res.Hits))
	for _, hit := range res.Hits {
		if songID, err := primitive.ObjectIDFromHex(hit.ID); err == nil {
			songIDs = append(songIDs, songID)
		}
	}
	return songIDs, nil
}

type searchCtx struct {
	cursorID   string
	ctx        context.Context
	cancelFunc context.CancelFunc
	resChan    <-chan *crawler.SearchResponse
	done       chan struct{}
	songChan   chan *Song
}

func (ctx *searchCtx) job() {
	log.Debug("SearchContext started", zap.String("CursorID", ctx.cursorID))
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
				)
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
					StoreID:        0,
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
				select {
				case ctx.songChan <- songX:
				default:
					log.Warn("Could not write to song channel")
				}
				continue
			}

			// If the song has not been downloaded from source yet, update the origin url
			if songX.StoreID == 0 {
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
				}
			}
		}
	}
	close(ctx.songChan)
	ctx.done <- struct{}{}
	log.Debug("SearchContext done", zap.String("CursorID", ctx.cursorID))
}

// StartSearch creates a new context and send the request to all the crawlers and waits for them to finish. If a context
// with the same id exists, we first cancel the old one and create a new context with new 'keyword'
func StartSearch(cursorID string, keyword string) <-chan *Song {
	ctx := getSearchCtx(cursorID)
	if ctx != nil {
		ctx.cancelFunc()

		// Wait for job to be finished
		<-ctx.done
	}

	ctx = &searchCtx{
		cursorID: cursorID,
		done:     make(chan struct{}, 1),
		songChan: make(chan *Song, 100),
	}
	ctx.ctx, ctx.cancelFunc = context.WithCancel(context.Background())
	ctx.resChan = crawler.Search(ctx.ctx, keyword)
	saveSearchCtx(ctx)
	go ctx.job()
	return ctx.songChan
}

// ResumeSearch checks if a context with 'cursorID' has been already exists and return the song channel, otherwise it returns
// nil.
func ResumeSearch(cursorID string) <-chan *Song {
	ctx := getSearchCtx(cursorID)
	if ctx == nil {
		return nil
	}

	return ctx.songChan
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
