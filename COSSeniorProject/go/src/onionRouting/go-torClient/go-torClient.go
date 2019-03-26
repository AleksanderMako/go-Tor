package main

import (
	"fmt"
	circuitrepository "onionRouting/go-torClient/repositories/circuit"
	peercredentialsrepository "onionRouting/go-torClient/repositories/credentials"
	cryptoservice "onionRouting/go-torClient/services/crypto/crypto-service"
	diffiehellmanservice "onionRouting/go-torClient/services/diffie-hellman"
	handshakeprotocolservice "onionRouting/go-torClient/services/handshake"
	onionprotocol "onionRouting/go-torClient/services/onion"
	storage "onionRouting/go-torClient/services/storage/storage-implementation"
	"os"
)

type CustomHandler struct{}

func main() {

	// generate key pair
	badgeDB := storage.NewStorage()
	//	databaseVolume, err := badgeDB.GetDBVolume()
	// if err != nil {
	// 	fmt.Println("error getting db volume", err.Error())
	// 	os.Exit(1)
	// }
	cryService := cryptoservice.NewCryptoService(badgeDB)

	dfhService := diffiehellmanservice.NewDiffieHellmanService(badgeDB, nil)

	hp := handshakeprotocolservice.NewHandshakeProtocol(*dfhService, cryService)
	circuitRepo := circuitrepository.NewPublicVariableRepository(badgeDB)
	peerCredentialsRepo := peercredentialsrepository.NewPeerCredentialsRepository(badgeDB)
	onionService := onionprotocol.NewOnionService(badgeDB, nil, *hp, circuitRepo, peerCredentialsRepo, cryService)
	peerList, err := onionService.GetPeers()
	if err != nil {
		fmt.Println("error getting peer list  ", err)
		os.Exit(1)
	}
	for _, peerID := range peerList {
		fmt.Println(peerID)
	}
	chainID, err := onionService.CreateOnionChain(peerList)
	if err != nil {
		fmt.Println("error while creating onion ring ", err.Error())
	}
	fmt.Println(chainID)

	if err = onionService.HandshakeWithPeers(chainID); err != nil {
		fmt.Println("error: ", err.Error())
		os.Exit(1)
	}

	// TODO:extract urls to env vars
	if err = onionService.GenerateSymetricKeys(chainID); err != nil {
		fmt.Println("error while exchanging symetric keys with peers " + err.Error())
		os.Exit(1)
	}
	destination := "registry:4500/peer/test"
	if err := onionService.BuildP2PCircuit([]byte(chainID), destination); err != nil {
		fmt.Println("error while building p2p circuit with peers " + err.Error())
		os.Exit(1)
	}
	if err = onionService.SendMessage([]byte(chainID), "hello server "); err != nil {
		fmt.Println("error while sending message " + err.Error())
		os.Exit(1)
	}
}

func HandleErr(err error, customErrMessage string) {
	if err != nil {
		fmt.Println(customErrMessage, err)
		os.Exit(1)
	}
	return
}
