package music

import (
	"github.com/blevesearch/bleve"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func SearchIndex(keyword string) ([]primitive.ObjectID, error) {
	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(keyword))
	res, err := songIndex.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	songIDs := make([]primitive.ObjectID, 0, len(res.Hits))
	for _, hit := range res.Hits {
		songID, err := primitive.ObjectIDFromHex(hit.ID)
		if err == nil {
			songIDs = append(songIDs, songID)
		}
	}
	return songIDs, nil
}
