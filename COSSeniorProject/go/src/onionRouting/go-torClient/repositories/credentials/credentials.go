package peercredentialsrepository

import (
	"encoding/json"
	storageserviceinterface "onionRouting/go-torClient/services/storage/storage-interface"
	"onionRouting/go-torClient/types"

	"github.com/dgraph-io/badger"

	"github.com/pkg/errors"
)

type PeerCredentials struct {
	db storageserviceinterface.StorageService
}

func NewPeerCredentialsRepository(db storageserviceinterface.StorageService) PeerCredentials {

	return PeerCredentials{
		db: db,
	}
}
func (this *PeerCredentials) GetPeerCredentials(peerID string) (types.PeerCredentials, error) {

	credentialBytes, err := this.db.Get(peerID)
	if err != nil {
		return types.PeerCredentials{}, errors.Wrap(err, "failed get Peer credentials bytes from the database ")
	}
	peerCredentials := types.PeerCredentials{}
	if err := json.Unmarshal(credentialBytes, &peerCredentials); err != nil {

		return types.PeerCredentials{}, errors.Wrap(err, "failed to unmarshal credentialBytes")
	}
	return peerCredentials, nil
}
func (this *PeerCredentials) SavePeerCredentials(peerID string, credentials types.PeerCredentials, dbVolume *badger.DB) error {

	savedCredentialBytes, err := this.db.Get(peerID)
	if err != nil && err != badger.ErrKeyNotFound {

		return errors.Wrap(err, "failed to get savedCredentialBytes ")
	}
	if savedCredentialBytes == nil {
		credentialBytes, err := json.Marshal(credentials)
		if err != nil {
			return errors.Wrap(err, "failed to marshal credentials")
		}
		if err := this.db.Put(peerID, credentialBytes, dbVolume); err != nil {
			return errors.Wrap(err, "failed to save peer credentials ")
		}
		return nil
	}
	peerCredentials := types.PeerCredentials{}
	if err := json.Unmarshal(savedCredentialBytes, &peerCredentials); err != nil {

		return errors.Wrap(err, "failed to unmarshal credentialBytes")
	}
	peerCredentials.SharedSecret = credentials.SharedSecret
	newPeerCredentialBytes, err := json.Marshal(peerCredentials)
	if err != nil {
		return errors.Wrap(err, "failed to to marshal peerCredentials")
	}
	if err = this.db.Put(peerID, newPeerCredentialBytes, dbVolume); err != nil {
		return errors.Wrap(err, "failed to save newPeerCredentialBytes in the database ")
	}

	return nil
}
