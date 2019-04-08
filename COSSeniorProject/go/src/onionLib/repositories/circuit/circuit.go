package circuitrepository

import (
	"encoding/json"
	storageserviceinterface "onionLib/services/storage/storage-interface"
	"onionLib/types"

	"github.com/dgraph-io/badger"

	"github.com/pkg/errors"
)

type CircuitRepository struct {
	db storageserviceinterface.StorageService
}

func NewPublicVariableRepository(db storageserviceinterface.StorageService) CircuitRepository {

	return CircuitRepository{
		db: db,
	}
}

func (this *CircuitRepository) Save(cID string, c types.Circuit, dbVolume *badger.DB) error {

	cBytes, e := json.Marshal(c)
	if e != nil {
		return errors.Wrap(e, "failed to marshal circuit to bytes ")
	}
	if e = this.db.Put(cID, cBytes); e != nil {
		return errors.Wrap(e, "failed to save circuit bytes in badger")
	}

	return nil
}
func (this *CircuitRepository) Get(cID string) (types.Circuit, error) {

	cBytes, e := this.db.Get(cID)
	if e != nil {
		return types.Circuit{}, errors.Wrap(e, "failed to get circuit in repo ")
	}

	if cBytes == nil {
		return types.Circuit{}, errors.New("failed to get circuit bytes for the given id " + cID)
	}
	c := types.Circuit{}
	if e = json.Unmarshal(cBytes, &c); e != nil {
		return types.Circuit{}, errors.Wrap(e, "failed to unmarshal bytes of circuit in repo")
	}
	return c, nil
}
