package storageserviceinterface

import "github.com/dgraph-io/badger"

type StorageService interface {
	GetDBVolume() (*badger.DB, error)
	Put(key string, val []byte, db *badger.DB) error
	Get(key string) ([]byte, error)
}
