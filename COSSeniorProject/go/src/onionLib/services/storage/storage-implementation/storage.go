package storage

import (
	"fmt"
	storageserviceinterface "onionLib/services/storage/storage-interface"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
)

type Storage struct {
	options badger.Options
}

func NewStorage(options badger.Options) storageserviceinterface.StorageService {

	return &Storage{
		options: options,
	}
}

func (this *Storage) Put(key string, val []byte) error {

	db, err := badger.Open(this.options)
	if err != nil {
		return errors.Wrap(err, "failed to open badger in storage service ")
	}
	defer db.Close()

	txn := db.NewTransaction(true)
	defer txn.Discard()

	err = txn.Set([]byte(key), val)
	if err != nil {
		return errors.Wrap(err, "failed to save data in badger db ")
	}

	if err := txn.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return nil
}
func (this *Storage) Get(key string) ([]byte, error) {

	db, err := badger.Open(this.options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger in storage service ")
	}
	defer db.Close()
	var value []byte
	err = db.View(func(tx *badger.Txn) error {

		item, err := tx.Get([]byte(key))
		if err != nil {
			return err
		}
		itemBytes, err := item.ValueCopy(nil)
		if err != nil {
			return errors.Wrap(err, "failed to read the value from badger ")
		}
		fmt.Println("data", string(itemBytes))
		value = itemBytes
		return nil
	})

	if err != nil {
		return nil, err
	}
	return value, nil
}
