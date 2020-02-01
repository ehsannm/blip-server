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

func UpdateLocalIndex(s *Song) error {
	return songIndex.Index(s.ID.Hex(), s)
}
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
