package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	diffiehellmanservice "onionRouting/go-torClient/services/diffie-hellman"
	handshakeprotocolservice "onionRouting/go-torClient/services/handshake"
	"onionRouting/go-torClient/services/request"
	"onionRouting/go-torClient/types"

	"os"
)

type CustomHandler struct{}

func main() {

	// generate key pair
	dfhService := diffiehellmanservice.NewDiffieHellmanService()
	hp := handshakeprotocolservice.NewHandshakeProtocol(*dfhService)
	pkBytes, _, err := hp.GenerateKeyPair()
	if err != nil {
		fmt.Println("error generating key pair ", err)
		os.Exit(1)
	}

	keyExchangeReq := types.Request{
		Action: "keyExchange",
		Data:   pkBytes,
	}
	url := "http://127.0.0.1:9000/keyExchange"

	res, err := request.Dial(url, keyExchangeReq)
	HandleErr(err, "")
	serverPublicKey := types.PubKey{}

	serverPublicKeyBytes, err := request.ParseResponse(res)
	HandleErr(err, "")

	err = json.Unmarshal(serverPublicKeyBytes, &serverPublicKey)
	HandleErr(err, "failed to unmarshal servers public key ")

	// generate diffie hellman koefs
	dfhKoefficients, err := hp.StartDiffieHellman()
	HandleErr(err, "error in  starting diffie hellman")

	//serialize handshake payload

	dfhBytes, err := json.Marshal(dfhKoefficients)
	HandleErr(err, "failed to marshal dfh keofficinets")
	request := types.Request{
		Action: "handleHandshake",
		Data:   dfhBytes,
	}

	// TODO fix peer's expected request payload
	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("error marshalling request payload ", err)
		os.Exit(1)
	}

	var buff bytes.Buffer
	buff.Write(requestBytes)

	resp, err := http.Post("http://127.0.0.1:9000/handshake", "application/json", &buff)
	if err != nil {
		fmt.Println("err making the request ", err)

		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	peerPubKey := types.PubKey{}

	err = json.Unmarshal(body, &peerPubKey)
	if err != nil {
		fmt.Println("error unmarshaling peer's pub key ")
		os.Exit(1)
	}
	fmt.Println("umarshaled payload ", peerPubKey.PubKey)
	//json.Unmarshal(body)

}
func HandleErr(err error, customErrMessage string) {
	if err != nil {
		fmt.Println(customErrMessage, err)
		os.Exit(1)
	}
	return
}
