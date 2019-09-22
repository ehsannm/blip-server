package goes

import (
	"github.com/jinzhu/gorm"
)

// Store is a type alias to avoid user to import gorm (and thus avoid version problems)
type Store = *gorm.DB

// Tx is a type alias
type Tx = *gorm.DB

// DB should be the only one way to access the DB in your application
var DB Store

// Init initialize the db package
func Init(databaseURL string) error {
	db, err := gorm.Open("postgres", databaseURL)
	if err != nil {
		return err
	}

	DB = db
	return DB.DB().Ping()
}

// IsRecordNotFoundError returns true if error is or contains a RecordNotFound error
func IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
