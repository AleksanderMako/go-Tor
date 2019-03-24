package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/big"
	onionrepository "onionRouting/go-torPeer/repositories/onion"
	cryptoserviceinterface "onionRouting/go-torPeer/services/crypto/crypto-service-interface"
	dfhservice "onionRouting/go-torPeer/services/diffie-hellman"
	storageserviceinterface "onionRouting/go-torPeer/services/storage/storage-interface"
	"onionRouting/go-torPeer/types"

	"github.com/pkg/errors"
)

type HandShakeController interface {
	HandleHandshake(data []byte) ([]byte, error)
	HandleKeyExchange(data []byte) ([]byte, error)
}

type TorHandshakeController struct {
	cryptoService   cryptoserviceinterface.CryptoService
	dfh             dfhservice.DFHService
	peerPrivateKey  *rsa.PrivateKey
	storageService  storageserviceinterface.StorageService
	onionRepository onionrepository.OnionRepository
}

func NewTorHandshakeController(cryptoService cryptoserviceinterface.CryptoService,
	dfh dfhservice.DFHService,
	storageService storageserviceinterface.StorageService,
	onionRepository onionrepository.OnionRepository) HandShakeController {

	return &TorHandshakeController{
		cryptoService:   cryptoService,
		dfh:             dfh,
		storageService:  storageService,
		onionRepository: onionRepository,
	}
}
func (this *TorHandshakeController) HandleHandshake(data []byte) ([]byte, error) {

	if data == nil {
		return nil, errors.New("Handshake controller got empty payload")
	}
	onionPayload := types.OnionPayload{}
	if err := json.Unmarshal(data, &onionPayload); err != nil {

		return nil, errors.Wrap(err, "failed to unmarshal onion payload ")
	}
	clientsPayload := types.DFHCoefficients{}
	err := json.Unmarshal(onionPayload.Coefficients, &clientsPayload)
	fmt.Println(onionPayload)
	if err != nil {
		fmt.Println("failed to unmarshal client's payload in Handshake Controller", err)
		return nil, errors.Wrap(err, "failed to unmarshal client's payload in Handshake Controller")
	}

	clientsPubKeyBytes, err := this.storageService.Get("clientpubKey")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read client's pub key from database")
	}
	clientsPubKey := types.PubKey{}
	err = json.Unmarshal(clientsPubKeyBytes, &clientsPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal the clients's public key")
	}

	err = this.cryptoService.Verify(clientsPayload.PublicVariable.Value, clientsPayload.PublicVariable.Signature, clientsPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify signature during handshake ")
	}
	// generate initial private part of the key
	privateVariable, err := this.dfh.Genrate_Private_Variable()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate peer's private variable")
	}

	// generate  shared secret by using clients public var and private variable
	clientsPublicVariable := new(big.Int)
	clientsPublicVariable.SetBytes(clientsPayload.PublicVariable.Value)
	sharedSecret, err := this.dfh.GenerateSharedSecret(clientsPublicVariable, privateVariable, clientsPayload.N)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate shared secret for peer ")
	}

	if err := this.onionRepository.SaveCircuitLink(onionPayload.CircuitID, types.CircuitLinkParameters{
		SharedSecret: sharedSecret,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to persist shared secret for circuit %v", onionPayload.CircuitID)
	}
	// generate public variable for client
	publicVariable, err := this.dfh.GeneratePublicVariable(clientsPayload.G, privateVariable, clientsPayload.N, this.peerPrivateKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.Wrap(err, "failed to generate peers' public variable ")
	}

	publicVariableBytes, err := json.Marshal(publicVariable)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal peer's public variable")
	}
	return publicVariableBytes, nil
}

func (this *TorHandshakeController) HandleKeyExchange(data []byte) ([]byte, error) {

	if data == nil {
		return nil, errors.New("empty public key payload")
	}
	e := this.storageService.Put("clientpubKey", data)
	if e != nil {
		return nil, errors.Wrap(e, "failed to write client's pub key in file")
	}
	clientsPublicKey := types.PubKey{}

	err := json.Unmarshal(data, &clientsPublicKey)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal clinet's pub key ")
	}
	//	fmt.Println("client's public key :", clientsPublicKey.PubKey)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate key pair in handshake controller ")
	}
	this.peerPrivateKey = privateKey
	publicKey := &privateKey.PublicKey

	myKey := types.PubKey{
		PubKey: *publicKey,
	}
	keyBytes, err := json.Marshal(myKey)
	if err != nil {
		return nil, errors.Wrap(err, " failed to marshal peer's public key in handshake controller")

	}
	return keyBytes, nil
}
