package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	onionlib "onionLib/lib/lib-implementation"
	messagerepository "onionRouting/go-torClient/repositories/message"
	storage "onionRouting/go-torClient/services/storage/storage-implementation"
	"onionRouting/go-torClient/types"
	"os"
)

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

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
	//torServer.listenHandle()
	httpAddr := flag.String("http", ":"+"8000", "Listen address")
	http.HandleFunc("/search", torServer.Search)
	http.HandleFunc("/connect", torServer.ConnectToServer)
	http.HandleFunc("/file", torServer.RequestTextFile)
	http.ListenAndServe(*httpAddr, nil)

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
	c_response := types.ConnectionResponse{
		Response: string(response),
	}
	cResponseBytes, err := json.Marshal(c_response)
	if err != nil {
		fmt.Println("failed to marshal connection response  " + err.Error())
		return nil, err
	}
	return cResponseBytes, nil
}
