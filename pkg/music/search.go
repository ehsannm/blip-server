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
	return songIndex.Index(s.ID.Hex(), s)
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
	for r := range ctx.resChan {
		uniqueKey := GenerateUniqueKey(r.Result.Title, r.Result.Artists)
		songX, err := GetSongByUniqueKey(uniqueKey)
		if err != nil {
			songX = &Song{
				ID:             primitive.NilObjectID,
				UniqueKey:      uniqueKey,
				Title:          r.Result.Title,
				Genre:          r.Result.Genre,
				Lyrics:         r.Result.Lyrics,
				Artists:        r.Result.Artists,
				Cdn:            "",
				OriginCoverUrl: r.Result.CoverUrl,
				OriginSongUrl:  r.Result.SongUrl,
				Source:         r.Source,
			}
			_, err = SaveSong(songX)
			if err != nil {
				log.Warn("Error On Save Search Result",
					zap.Error(err),
					zap.String("Source", r.Source),
					zap.String("Title", r.Result.Title),
				)
				continue
			}
			select {
			case ctx.songChan <- songX:
			default:
			}
			continue
		}
		songX.Artists = r.Result.Artists
		songX.Title = r.Result.Title
		songX.Genre = r.Result.Genre
		songX.Lyrics = r.Result.Lyrics
		songX.OriginSongUrl = r.Result.SongUrl
		songX.OriginCoverUrl = r.Result.CoverUrl
		songX.Source = r.Source
		_, err = SaveSong(songX)
		if err != nil {
			log.Warn("Error On Save Search Result",
				zap.Error(err),
				zap.String("Source", r.Source),
				zap.String("Title", r.Result.Title),
			)
		}
	}
	close(ctx.songChan)
	ctx.done <- struct{}{}
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
		done:     make(chan struct{}),
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
	v := searchContexts[cursorID]
	searchContextsMtx.RUnlock()
	return v
}
