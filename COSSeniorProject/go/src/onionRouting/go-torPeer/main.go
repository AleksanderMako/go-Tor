package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	controller "onionRouting/go-torPeer/controllers"
	"onionRouting/go-torPeer/types"
	"os"
)

type PubKey struct {
	PubKey rsa.PublicKey `json"pubkey"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
func HandShakeHandler(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		os.Exit(1)
	}

	clientsPayload := types.HandshakePayload{}
	err = json.Unmarshal(body, &clientsPayload)
	if err != nil {
		fmt.Println("failed to unmarshal clients pub key ", err)
		os.Exit(1)
	}
	fmt.Println("clients pub key is   ", clientsPayload.PublicKey)
	fmt.Println("clients public g is   ", clientsPayload.DFH.G)
	fmt.Println("clients public  n   ", clientsPayload.DFH.N)
	fmt.Println("clients public variable is    ", clientsPayload.DFH.PublicVariable)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		w.Write([]byte("error happened during key "))
	}
	publicKey := &privateKey.PublicKey

	myKey := PubKey{
		PubKey: *publicKey,
	}
	keyBytes, err := json.Marshal(myKey)
	if err != nil {
		w.Write([]byte("error happened during key "))

	}
	w.Write(keyBytes)

}

func main() {

	handShakeController := controller.NewTorHandshakeController()
	multiplexer := NewMultiplexer(handShakeController)

	server := http.Server{
		Addr: "127.0.0.1:9000",
	}
	http.HandleFunc("/", multiplexer.MultiplexRequest)
	server.ListenAndServe()
}
