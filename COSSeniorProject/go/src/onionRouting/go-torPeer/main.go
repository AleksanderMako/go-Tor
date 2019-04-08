package main

import (
	"flag"
	"fmt"
	"net/http"
	clientcapabilities "onionRouting/go-torPeer/client-capabilities"
	controller "onionRouting/go-torPeer/controllers"
	onionrepository "onionRouting/go-torPeer/repositories/onion"
	cryptoservice "onionRouting/go-torPeer/services/crypto/crypto-service"
	dfhservice "onionRouting/go-torPeer/services/diffie-hellman"
	peeronionprotocol "onionRouting/go-torPeer/services/onion"
	storage "onionRouting/go-torPeer/services/storage/storage-implementation"
	"os"
)

func main() {

	badgerDB := storage.NewStorage()
	onionRepo := onionrepository.NewOnionRepository(badgerDB)
	cryptoService := cryptoservice.NewCryptoService(badgerDB)
	dfh := dfhservice.NewDfhService(cryptoService, badgerDB)
	handShakeController := controller.NewTorHandshakeController(cryptoService, *dfh, badgerDB, onionRepo)
	onionService := peeronionprotocol.NewOnionService(onionRepo, cryptoService)
	onionServiceController := controller.NewOnionCOntroller(onionService)
	multiplexer := NewMultiplexer(handShakeController, onionServiceController)

	port := os.Getenv("PEER_PORT")

	fmt.Println("Peer started listening on port " + port)
	startUp(port)
	http.HandleFunc("/", multiplexer.MultiplexRequest)
	//http.HandleFunc("/contactIP")
	httpAddr := flag.String("http", ":"+port, "Listen address")

	http.ListenAndServe(*httpAddr, nil)
}
func startUp(port string) {
	err := clientcapabilities.RegisterPeer()
	if err != nil {
		fmt.Println("error during peer registration" + err.Error())
		os.Exit(1)
	}
	err = clientcapabilities.GetPeerAddresses("http://registry:4500/peer/peers")
	if err != nil {
		fmt.Println("error while getting peer addresses", err.Error())
		os.Exit(1)
	}

}
