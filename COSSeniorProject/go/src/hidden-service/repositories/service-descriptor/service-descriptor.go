package servicedescriptorrepository

import (
	"crypto/sha256"
	"encoding/json"
	storageserviceinterface "hidden-service/services/storage/storage-interface"
	"hidden-service/types"

	"github.com/pkg/errors"
)

type ServiceDescriptorRepository struct {
	db storageserviceinterface.StorageService
}

func NewServiceDescriptorRepository(db storageserviceinterface.StorageService) ServiceDescriptorRepository {

	return ServiceDescriptorRepository{
		db: db,
	}
}
func (this *ServiceDescriptorRepository) Save(serviceDescriptor types.ServiceDescriptor) error {

	serviceDescriptorBytes, e := json.Marshal(serviceDescriptor)
	if e != nil {
		return errors.Wrap(e, "failed to marshal serviceDescriptor to bytes ")
	}
	key := serviceDescriptor.ID
	hashedKey, e := this.createHash(key)
	if e != nil {
		return errors.Wrap(e, "failed to create hash ")
	}
	if e = this.db.Put(hashedKey, serviceDescriptorBytes); e != nil {

		return errors.Wrap(e, "failed to save service descriptor ")
	}
	return nil
}
func (this *ServiceDescriptorRepository) createHash(data []byte) ([]byte, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make hash ")
	}
	hashedData := hasher.Sum(nil)
	return hashedData, nil
}
func (this *ServiceDescriptorRepository) Get(key []byte) (types.ServiceDescriptor, error) {

	hashedKey, err := this.createHash(key)
	if err != nil {
		return types.ServiceDescriptor{}, errors.Wrap(err, "failed to recreate descriptor key during get opration")
	}
	savedDescriptorBytes, e := this.db.Get(hashedKey)
	if e != nil {
		return types.ServiceDescriptor{}, errors.Wrap(e, "failed to get savedDescriptorBytes")
	}
	savedDescriptor := types.ServiceDescriptor{}
	if e = json.Unmarshal(savedDescriptorBytes, &savedDescriptor); e != nil {
		return types.ServiceDescriptor{}, errors.Wrap(e, "failed to unmarshal savedDescriptorBytes ")
	}
	return savedDescriptor, nil
}
