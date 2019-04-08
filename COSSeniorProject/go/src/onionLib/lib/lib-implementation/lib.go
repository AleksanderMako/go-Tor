package onionlib

import (
	circuitrepository "onionLib/repositories/circuit"
	peercredentialsrepository "onionLib/repositories/credentials"
	cryptoservice "onionLib/services/crypto/crypto-service"
	diffiehellmanservice "onionLib/services/diffie-hellman"
	handshakeprotocolservice "onionLib/services/handshake"
	onionprotocol "onionLib/services/onion"
	storage "onionLib/services/storage/storage-implementation"

	"github.com/dgraph-io/badger"
)

type OnionLibrary struct {
	Onionservice onionprotocol.OnionService
}

func NewOnionLib(options badger.Options) OnionLibrary {
	badgeDB := storage.NewStorage(options)
	cryService := cryptoservice.NewCryptoService(badgeDB)
	dfhService := diffiehellmanservice.NewDiffieHellmanService(badgeDB, nil)
	hp := handshakeprotocolservice.NewHandshakeProtocol(*dfhService, cryService)
	circuitRepo := circuitrepository.NewPublicVariableRepository(badgeDB)
	peerCredentialsRepo := peercredentialsrepository.NewPeerCredentialsRepository(badgeDB)
	onionService := onionprotocol.NewOnionService(badgeDB, nil, *hp, circuitRepo, peerCredentialsRepo, cryService)
	return OnionLibrary{
		Onionservice: onionService,
	}
}
