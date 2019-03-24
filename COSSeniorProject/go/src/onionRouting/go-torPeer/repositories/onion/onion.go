package onionrepository

import (
	"encoding/json"
	storageserviceinterface "onionRouting/go-torPeer/services/storage/storage-interface"
	"onionRouting/go-torPeer/types"

	"github.com/dgraph-io/badger"

	"github.com/pkg/errors"
)

type OnionRepository struct {
	db storageserviceinterface.StorageService
}

func NewOnionRepository(db storageserviceinterface.StorageService) OnionRepository {

	onionRepo := new(OnionRepository)
	onionRepo.db = db
	return *onionRepo
}

func (this *OnionRepository) SaveCircuitLink(cID []byte, link types.CircuitLinkParameters) error {

	savedLinkBytes, err := this.db.Get(string(cID))
	if err != nil && err != badger.ErrKeyNotFound {
		return errors.Wrap(err, "failed to lookup savedLinkBytes ")
	}
	linkBytes, e := json.Marshal(link)
	if e != nil {
		return errors.Wrap(e, "failed to marshal link to bytes in SaveCircuitLink method ")
	}
	if savedLinkBytes == nil {
		if err := this.db.Put(string(cID), linkBytes); err != nil {
			return errors.Wrap(err, "failed to save link in SaveCircuitLink method")
		}
		return nil
	}

	savedLink := types.CircuitLinkParameters{}
	if e := json.Unmarshal(savedLinkBytes, &savedLink); e != nil {
		return errors.Wrap(e, "failed to unmarshal saved bytes in SaveCircuitLink")
	}
	savedLink.Previous = link.Previous
	savedLink.Next = link.Next
	newSavedLinkBytes, e := json.Marshal(savedLink)
	if e != nil {
		return errors.Wrap(e, "failed to marshal newSavedLink")
	}
	if e = this.db.Put(string(cID), newSavedLinkBytes); e != nil {
		return errors.Wrap(e, "failed to save link in badger ")
	}
	return nil
}
func (this *OnionRepository) GetCircuitLinkParamaters(cID []byte) (types.CircuitLinkParameters, error) {

	savedLinkBytes, e := this.db.Get(string(cID))
	if e != nil {
		return types.CircuitLinkParameters{}, errors.Wrap(e, "failed to get savedLinkBytes from badger")
	}
	savedLink := types.CircuitLinkParameters{}
	if e = json.Unmarshal(savedLinkBytes, &savedLink); e != nil {
		return types.CircuitLinkParameters{}, errors.Wrap(e, "failed to get saved link in onion repository ")
	}
	return savedLink, nil
}
