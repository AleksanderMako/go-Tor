package main

import (
	"fmt"
	cryptoservice "onionRouting/go-torClient/services/crypto/crypto-service"
	diffiehellmanservice "onionRouting/go-torClient/services/diffie-hellman"
	handshakeprotocolservice "onionRouting/go-torClient/services/handshake"
	storage "onionRouting/go-torClient/services/storage/storage-implementation"

	onionprotocol "onionRouting/go-torClient/services/onion"
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

	onionService := onionprotocol.NewOnionService(badgeDB, nil, *hp)
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
	//cryService := cryptoservice.NewCryptoService(badgeDB)
	//dfhService := diffiehellmanservice.NewDiffieHellmanService(badgeDB)
	//hp := handshakeprotocolservice.NewHandshakeProtocol(*dfhService, cryService)
	// pkBytes, privateKey, err := hp.GenerateKeyPair()
	// if err != nil {
	// 	fmt.Println("error generating key pair ", err)
	// 	os.Exit(1)
	// }
	// keyExchangeReq := types.Request{
	// 	Action: "keyExchange",
	// 	Data:   pkBytes,
	// }
	// url := "http://127.0.0.1:9000/keyExchange"

	// res, err := request.Dial(url, keyExchangeReq)
	// HandleErr(err, "")
	// serverPublicKey := types.PubKey{}

	// serverPublicKeyBytes, err := request.ParseResponse(res)
	// HandleErr(err, "")

	// err = json.Unmarshal(serverPublicKeyBytes, &serverPublicKey)
	// HandleErr(err, "failed to unmarshal servers public key ")

	// //fmt.Println("umarshaled payload ", serverPublicKey.PubKey)

	// //generate diffie hellman koefs
	// dfhKoefficients, err := hp.StartDiffieHellman(privateKey)
	// HandleErr(err, "error in  starting diffie hellman")

	// //serialize handshake payload

	// dfhBytes, err := json.Marshal(dfhKoefficients)
	// HandleErr(err, "failed to marshal dfh keofficinets")
	// req := types.Request{
	// 	Action: "handleHandshake",
	// 	Data:   dfhBytes,
	// }
	// newUrl := "http://127.0.0.1:9000/handshake"
	// res, err = request.Dial(newUrl, req)
	// HandleErr(err, "failed to dial dfh endpoint")

	// peerPublicVariable := types.PublicVariable{}
	// peerPublicVariableBytes, err := request.ParseResponse(res)
	// fmt.Println(string(peerPublicVariableBytes))
	// err = json.Unmarshal(peerPublicVariableBytes, &peerPublicVariable)
	// HandleErr(err, "failed to unmarshal peer's public variable")

	// pPublicVar := new(big.Int)
	// pPublicVar.SetBytes(peerPublicVariable.Value)
	// // cPublicVar := new(big.Int)
	// // cPublicVar.SetBytes(dfhKoefficients.PublicVariable.Value)
	// // if pPublicVar == cPublicVar {
	// // 	fmt.Println("its the same")
	// // }
	// println()
	// println()
	// //	fmt.Println("peer's dfh public variable is ", pPublicVar)

	//hp.GenerateSharedSecret(pPublicVar, dfhKoefficients.N)
}
func HandleErr(err error, customErrMessage string) {
	if err != nil {
		fmt.Println(customErrMessage, err)
		os.Exit(1)
	}
	return
}
