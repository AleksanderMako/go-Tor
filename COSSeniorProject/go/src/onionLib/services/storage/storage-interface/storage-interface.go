package storageserviceinterface

type StorageService interface {
	Put(key string, val []byte) error
	Get(key string) ([]byte, error)
}
