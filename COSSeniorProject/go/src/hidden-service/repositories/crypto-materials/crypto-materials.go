package cryptomaterialsrepository

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	storageserviceinterface "hidden-service/services/storage/storage-interface"
	"hidden-service/types"

	"github.com/pkg/errors"
)

type CryptoMaterialsRepository struct {
	db storageserviceinterface.StorageService
}

func NewCryptoMaterialsRepository(db storageserviceinterface.StorageService) CryptoMaterialsRepository {

	cryptoMaterials := new(CryptoMaterialsRepository)
	cryptoMaterials.db = db
	return *cryptoMaterials
}

//GenerateKeyPair  method returns bytes of public key wrapped
// with custom types.PubKey datatype
func (this *CryptoMaterialsRepository) GenerateKeyPair() ([]byte, *rsa.PrivateKey, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate key pair in handshake protocol ")
	}
	publicKey := &privateKey.PublicKey
	hpPublicKey := types.PubKey{
		PubKey: *publicKey,
	}
	pubKeyBytes, err := json.Marshal(hpPublicKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to serialize public key to bytes in handshake protocol ")
	}

	return pubKeyBytes, privateKey, nil
}
