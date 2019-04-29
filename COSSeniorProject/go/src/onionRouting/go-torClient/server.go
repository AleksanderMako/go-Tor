package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	onionlib "onionLib/lib/lib-implementation"
	messagerepository "onionRouting/go-torClient/repositories/message"
	"onionRouting/go-torClient/types"

	"github.com/pkg/errors"
)

func (this *TorServer) listenHandle() {

	http.HandleFunc("/search", this.Search)
	http.HandleFunc("/connect", this.ConnectToServer)
	http.HandleFunc("/file", this.RequestTextFile)
}

func readBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request ")
	}
	return body, nil
}

type TorServer struct {
	torLib      onionlib.OnionLibrary
	PublicKey   []byte
	PrivateKey  []byte
	messageRepo messagerepository.MessageRepository
	chainId     string
}

func NewTorServer(torLib onionlib.OnionLibrary, PublicKey []byte, PrivateKey []byte, messageRepo messagerepository.MessageRepository) TorServer {
	return TorServer{
		torLib:      torLib,
		PublicKey:   PublicKey,
		PrivateKey:  PrivateKey,
		messageRepo: messageRepo,
	}
}

func (this *TorServer) Search(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)
	body, err := readBody(r)
	if err != nil {
		fmt.Printf("error in reading request for search %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query := types.Query{}

	if err := json.Unmarshal(body, &query); err != nil {
		fmt.Printf("error in reading request for search %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	descriptor, err := this.torLib.Onionservice.GetServiceDescriptorsByKeyWords(query.KeyWord)
	if err != nil {
		fmt.Printf("error in reading request for search %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	descriptorBytes, err := json.Marshal(descriptor)
	if err != nil {
		fmt.Printf("error in reading request for search %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(descriptorBytes)
}
func (this *TorServer) ConnectToServer(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)
	body, err := readBody(r)
	if err != nil {
		fmt.Printf("error in reading request for search %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	c := types.Connect{}

	if err := json.Unmarshal(body, &c); err != nil {
		fmt.Printf("error in reading request for ConnectToServer %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ip := c.Ip
	peerList, err := this.torLib.Onionservice.GetPeers()
	if err != nil {
		fmt.Printf("error in reading request for ConnectToServer %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	destination := "registry:4500/peer/test"
	peerList = append(peerList, ip)
	decodedID, err := base64.StdEncoding.DecodeString(c.DescriptorID)
	if err != nil {
		fmt.Printf("error in decoding decriptor %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	hiddenResponse, err := this.Connect(peerList, destination, decodedID)
	if err != nil {
		err = errors.Wrap(err, "failed to connect to server ")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(hiddenResponse)
}
func (this *TorServer) RequestTextFile(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)
	body, err := readBody(r)
	if err != nil {
		fmt.Printf("error in reading request for search %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	c := types.Connect{}

	if err := json.Unmarshal(body, &c); err != nil {
		fmt.Printf("error in reading request for RequestTextFile %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	decodedID, err := base64.StdEncoding.DecodeString(c.DescriptorID)
	if err != nil {
		fmt.Printf("error in decoding decriptor %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	messageBytes, err := this.messageRepo.CreateMessage(decodedID, c.Keyword)
	if err != nil {
		fmt.Println("error while building p2p circuit with peers " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Printf("constructed message is %v\n", string(messageBytes))
	hiddenResponse, err := this.torLib.Onionservice.SendMessage([]byte(this.chainId), messageBytes)
	if err != nil {
		err = errors.Wrap(err, "failed to connect to server ")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var fileEtension string
	var data string
	if c.Keyword == "text" {
		fileEtension = "txt"
		data = string(hiddenResponse)
	} else {
		fileEtension = "jpg"
		b64Blob := base64.StdEncoding.EncodeToString(hiddenResponse)
		data = b64Blob
	}

	file := types.FileResponse{
		Data:     data,
		FileType: fileEtension,
	}
	fileBytes, err := json.Marshal(file)
	if err != nil {
		fmt.Println("failed to marshal file " + err.Error())

	}
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}
