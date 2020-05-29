package driver

import (
	badger "github.com/dgraph-io/badger/v2"
)

// DB ...
type DB struct {
	Badger *badger.DB
}

// BadgerCon
var badgerCon = &DB{}

// ConnectBadger creates a badgerDB connection
func ConnectBadger(path string) (*DB, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	badgerCon.Badger = db

	return badgerCon, err
}
