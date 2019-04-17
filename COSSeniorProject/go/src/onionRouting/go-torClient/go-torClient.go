package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	onionlib "onionLib/lib/lib-implementation"
	messagerepository "onionRouting/go-torClient/repositories/message"
	storage "onionRouting/go-torClient/services/storage/storage-implementation"
	"os"
)

func main() {

	//TODO:extract the entire onion protocol into a lib to use in server and client
	// generate key pair
	badgerOptions := storage.InitializeBadger()
	messageRepo := messagerepository.NewMessageRepository()
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
	privateKeyBytes, err := json.Marshal(privateKey)
	if err != nil {
		fmt.Println("error while Marshaling private key in client  ", err.Error())
		os.Exit(1)
	}

	torServer := NewTorServer(onionLib, publicKey, privateKeyBytes, messageRepo)
	torServer.listenHandle()
	httpAddr := flag.String("http", ":"+"8000", "Listen address")
	http.ListenAndServe(*httpAddr, nil)

	// get introduction point

	// descriptors, err := onionLib.Onionservice.GetServiceDescriptorsByKeyWords("testing")
	// if err != nil {
	// 	fmt.Println("error getting service descriptor  list  ", err.Error())
	// 	os.Exit(1)
	// }

	// ip := descriptors.ServiceDescriptors[0].IntroductionPoints[0]
	// destination := "registry:4500/peer/test"
	// peerList = append(peerList, ip)

	// privateKeyBytes, err := json.Marshal(privateKey)
	// if err != nil {
	// 	fmt.Println("error while Marshaling private key in client  ", err.Error())
	// 	os.Exit(1)
	// }
	// chainID, err := onionLib.Onionservice.CreateOnionChain(peerList, publicKey)
	// // 	chainID, err := onionService.CreateOnionChain(peerList)
	// if err != nil {
	// 	fmt.Println("error while creating onion ring ", err.Error())
	// 	os.Exit(1)
	// }
	// fmt.Println(chainID)

	// if err = onionLib.Onionservice.HandshakeWithPeers(chainID, publicKey, privateKeyBytes); err != nil {
	// 	fmt.Println("error: ", err.Error())
	// 	os.Exit(1)
	// }
	// if err = onionLib.Onionservice.GenerateSymetricKeys(chainID); err != nil {
	// 	fmt.Println("error while exchanging symmetric keys with peers " + err.Error())
	// 	os.Exit(1)
	// }
	// client := "torclient:8000"
	// if err = onionLib.Onionservice.BuildP2PCircuit([]byte(chainID), client, destination); err != nil {
	// 	fmt.Println("error while building p2p circuit with peers " + err.Error())
	// 	os.Exit(1)
	// }

	// messageBytes, err := messageRepo.CreateMessage(descriptors.ServiceDescriptors[0].ID, "txt")
	// if err != nil {
	// 	fmt.Println("error while building p2p circuit with peers " + err.Error())
	// 	os.Exit(1)
	// }

	// if err = onionLib.Onionservice.SendMessage([]byte(chainID), messageBytes); err != nil {
	// 	fmt.Println("error while sending message " + err.Error())
	// 	os.Exit(1)
	// }
}

func HandleErr(err error, customErrMessage string) {
	if err != nil {
		fmt.Println(customErrMessage, err)
		os.Exit(1)
	}
	return
}
func (this *TorServer) Connect(peerList []string, destination string, descriptorID []byte) ([]byte, error) {
	chainID, err := this.torLib.Onionservice.CreateOnionChain(peerList, this.PublicKey)
	// 	chainID, err := onionService.CreateOnionChain(peerList)
	if err != nil {
		fmt.Println("error while creating onion ring ", err.Error())
		return nil, err
	}
	this.chainId = chainID
	fmt.Println(chainID)

	if err = this.torLib.Onionservice.HandshakeWithPeers(chainID, this.PublicKey, this.PrivateKey); err != nil {
		fmt.Println("error: ", err.Error())
		return nil, err
	}
	if err = this.torLib.Onionservice.GenerateSymetricKeys(chainID); err != nil {
		fmt.Println("error while exchanging symmetric keys with peers " + err.Error())
		return nil, err
	}
	client := "torclient:8000"
	if err = this.torLib.Onionservice.BuildP2PCircuit([]byte(chainID), client, destination); err != nil {
		fmt.Println("error while building p2p circuit with peers " + err.Error())
		return nil, err
	}
	messageBytes, err := this.messageRepo.CreateMessage(descriptorID, "connect")
	if err != nil {
		fmt.Println("error while building p2p circuit with peers " + err.Error())
		return nil, err
	}
	response, err := this.torLib.Onionservice.SendMessage([]byte(chainID), messageBytes)
	if err != nil {
		fmt.Println("error while sending message " + err.Error())
		return nil, err
	}
	return response, nil
}
