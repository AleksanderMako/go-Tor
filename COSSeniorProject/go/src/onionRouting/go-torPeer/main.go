package main

import (
	"net/http"
	controller "onionRouting/go-torPeer/controllers"
	cryptoservice "onionRouting/go-torPeer/services/crypto/crypto-service"
	dfhservice "onionRouting/go-torPeer/services/diffie-hellman"
)

func main() {
	cryptoService := cryptoservice.NewCryptoService()
	dfh := dfhservice.NewDfhService(cryptoService)
	handShakeController := controller.NewTorHandshakeController(cryptoService, *dfh)
	multiplexer := NewMultiplexer(handShakeController)

	server := http.Server{
		Addr: "127.0.0.1:9000",
	}
	http.HandleFunc("/", multiplexer.MultiplexRequest)
	server.ListenAndServe()
}
