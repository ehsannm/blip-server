package music

import (
	"context"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	"github.com/blevesearch/bleve"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// UpdateLocalIndex updates the local index which will be used by search handlers
func UpdateLocalIndex(s *Song) error {
	if d, _ := songIndex.Document(s.ID.Hex()); d == nil {
		return songIndex.Index(s.ID.Hex(), s)
	}
	return nil
}

// SearchLocalIndex
func SearchLocalIndex(keyword string) ([]primitive.ObjectID, error) {
	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(keyword))
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
	for r := range ctx.resChan {
		for _, foundSong := range r.Result {
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
				}
				continue
			}
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
