package main

import (
	"encoding/json"
	"fmt"
	onionlib "onionLib/lib/lib-implementation"
	storage "onionRouting/go-torClient/services/storage/storage-implementation"
	"os"
)

func main() {

	//TODO:extract the entire onion protocol into a lib to use in server and client
	// generate key pair
	badgerOptions := storage.InitializeBadger()
	publicKey, privateKey, err := onionlib.CreateCryptoMaterials()
	if err != nil {
		fmt.Println("error while CreateCryptoMaterials  ", err.Error())
		os.Exit(1)

	}
	onionLib := onionlib.NewOnionLib(badgerOptions, publicKey)
	peerList, err := onionLib.Onionservice.GetPeers()
	if err != nil {
		fmt.Println("error getting peer list  ", err)
		os.Exit(1)
	}
	for _, peerID := range peerList {
		fmt.Println(peerID)
	}
	// get introduction point

	descriptors, err := onionLib.Onionservice.GetServiceDescriptorsByKeyWords("testing")
	if err != nil {
		fmt.Println("error getting service descriptor  list  ", err.Error())
		os.Exit(1)
	}
	ip := descriptors.ServiceDescriptors[0].IntroductionPoints[0]
	destination := "registry:4500/peer/test"
	peerList = append(peerList, ip)

	privateKeyBytes, err := json.Marshal(privateKey)
	if err != nil {
		fmt.Println("error while Marshaling private key in client  ", err.Error())
		os.Exit(1)
	}
	chainID, err := onionLib.Onionservice.CreateOnionChain(peerList, publicKey)
	// 	chainID, err := onionService.CreateOnionChain(peerList)
	if err != nil {
		fmt.Println("error while creating onion ring ", err.Error())
		os.Exit(1)
	}
	fmt.Println(chainID)

	if err = onionLib.Onionservice.HandshakeWithPeers(chainID, publicKey, privateKeyBytes); err != nil {
		fmt.Println("error: ", err.Error())
		os.Exit(1)
	}
	if err = onionLib.Onionservice.GenerateSymetricKeys(chainID); err != nil {
		fmt.Println("error while exchanging symmetric keys with peers " + err.Error())
		os.Exit(1)
	}
	client := "torclient:8000"
	//	destination := "registry:4500/peer/test"
	if err = onionLib.Onionservice.BuildP2PCircuit([]byte(chainID), client, destination); err != nil {
		fmt.Println("error while building p2p circuit with peers " + err.Error())
		os.Exit(1)
	}
	// 	// message
	if err = onionLib.Onionservice.SendMessage([]byte(chainID), string(descriptors.ServiceDescriptors[0].ID)); err != nil {
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
