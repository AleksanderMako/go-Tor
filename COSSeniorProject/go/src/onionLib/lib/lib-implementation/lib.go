package onionlib

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	circuitrepository "onionLib/repositories/circuit"
	peercredentialsrepository "onionLib/repositories/credentials"
	cryptoservice "onionLib/services/crypto/crypto-service"
	diffiehellmanservice "onionLib/services/diffie-hellman"
	handshakeprotocolservice "onionLib/services/handshake"
	onionprotocol "onionLib/services/onion"
	storage "onionLib/services/storage/storage-implementation"
	"onionLib/types"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
)

type OnionLibrary struct {
	Onionservice onionprotocol.OnionService
}

func NewOnionLib(options badger.Options, publicKey []byte) OnionLibrary {
	badgeDB := storage.NewStorage(options)
	cryService := cryptoservice.NewCryptoService(badgeDB)
	dfhService := diffiehellmanservice.NewDiffieHellmanService(badgeDB, nil)
	hp := handshakeprotocolservice.NewHandshakeProtocol(*dfhService, cryService)
	circuitRepo := circuitrepository.NewPublicVariableRepository(badgeDB)
	peerCredentialsRepo := peercredentialsrepository.NewPeerCredentialsRepository(badgeDB)
	onionService := onionprotocol.NewOnionService(badgeDB, nil, *hp, circuitRepo, peerCredentialsRepo, cryService, publicKey)
	return OnionLibrary{
		Onionservice: onionService,
	}
}
func CreateCryptoMaterials() ([]byte, types.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, types.PrivateKey{}, errors.Wrap(err, "failed to generate key pair in handshake protocol ")
	}
	publicKey := &privateKey.PublicKey
	hpPublicKey := types.PubKey{
		PubKey: *publicKey,
	}
	privKey := types.PrivateKey{
		PrivateKey: *privateKey,
	}
	pubKeyBytes, err := json.Marshal(hpPublicKey)
	if err != nil {
		return nil, types.PrivateKey{}, errors.Wrap(err, "failed to serialize public key to bytes in handshake protocol ")
	}
	return pubKeyBytes, privKey, nil
}
