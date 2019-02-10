package storageserviceinterface

type storageService interface {
	Put(key string, val []byte) error
	Get(key string) ([]byte, error)
}
