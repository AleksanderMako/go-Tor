package storageserviceinterface

import "github.com/dgraph-io/badger"

type StorageService interface {
	GetDBVolume() (*badger.DB, error)
	Put(key []byte, val []byte) error
	Get(key []byte) ([]byte, error)
}
