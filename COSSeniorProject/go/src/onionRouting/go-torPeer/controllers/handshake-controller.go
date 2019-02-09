package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"onionRouting/go-torPeer/types"

	"github.com/pkg/errors"
)

type HandShakeController interface {
	HandleHandshake(data []byte) ([]byte, error)
	HandleKeyExchange(data []byte) ([]byte, error)
}

type TorHandshakeController struct {
}

func NewTorHandshakeController() HandShakeController {

	return &TorHandshakeController{}
}
func (this *TorHandshakeController) HandleHandshake(data []byte) ([]byte, error) {

	if data == nil {
		return nil, errors.New("Handshake controller got empty payload")
	}
	clientsPayload := types.HandshakePayload{}
	err := json.Unmarshal(data, &clientsPayload)
	if err != nil {
		fmt.Println("failed to unmarshal clients pub key ", err)
		return nil, errors.Wrap(err, "failed to unmarshal client's payload in Handshake Controller")
	}
	clientsPublicKey := types.PubKey{}

	err = json.Unmarshal(clientsPayload.PublicKey, &clientsPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal client's public key in Handshake Controller")
	}
	fmt.Println("clients pub key is   ", clientsPublicKey)
	fmt.Println("clients public g is   ", clientsPayload.DFH.G)
	fmt.Println("clients public  n   ", clientsPayload.DFH.N)
	fmt.Println("clients public variable is    ", clientsPayload.DFH.PublicVariable)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate key pair in handshake controller ")
	}
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

func (this *TorHandshakeController) HandleKeyExchange(data []byte) ([]byte, error) {

	clientsPublicKey := types.PubKey{}

	err := json.Unmarshal(data, &clientsPublicKey)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal clinet's pub key ")
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate key pair in handshake controller ")
	}
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
