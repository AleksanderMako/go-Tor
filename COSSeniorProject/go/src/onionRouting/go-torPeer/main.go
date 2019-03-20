package main

import (
	"fmt"
	"net/http"
	clientcapabilities "onionRouting/go-torPeer/client-capabilities"
	controller "onionRouting/go-torPeer/controllers"
	cryptoservice "onionRouting/go-torPeer/services/crypto/crypto-service"
	dfhservice "onionRouting/go-torPeer/services/diffie-hellman"
	storage "onionRouting/go-torPeer/services/storage/storage-implementation"
	"os"
)

func main() {
	badgerDB := storage.NewStorage()
	cryptoService := cryptoservice.NewCryptoService(badgerDB)
	dfh := dfhservice.NewDfhService(cryptoService, badgerDB)
	handShakeController := controller.NewTorHandshakeController(cryptoService, *dfh)
	multiplexer := NewMultiplexer(handShakeController)

	port := os.Getenv("PEER_PORT")
	_ = http.Server{
		Addr: "127.0.0.1:" + port,
	}
	fmt.Println("Peer started listening on port " + port)
	os.Setenv("PEER_ADD", "127.0.0.1:"+port)

	err := clientcapabilities.RegisterPeer()
	if err != nil {
		fmt.Println("error during peer registration" + err.Error())
		os.Exit(1)
	}
	http.HandleFunc("/", multiplexer.MultiplexRequest)
	//	server.ListenAndServe()
}
