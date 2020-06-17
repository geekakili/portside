package driver

import (
	"github.com/zippoxer/bow"
)

// DB ...
type DB struct {
	Badger *bow.DB
}

// BadgerCon
var badgerCon = &DB{}

// ConnectBadger creates a badgerDB connection
func ConnectBadger(path string) (*DB, error) {
	db, err := bow.Open(path)
	badgerCon.Badger = db

	return badgerCon, err
}
