package handshakeprotocolservice

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	diffiehellmanservice "onionRouting/go-torClient/services/diffie-hellman"
	"onionRouting/go-torClient/types"

	"github.com/pkg/errors"
)

type HandshakeProtocolService struct {
	dh diffiehellmanservice.DiffiHellmanService
}

func NewHandshakeProtocol(dfh diffiehellmanservice.DiffiHellmanService) *HandshakeProtocolService {

	hp := new(HandshakeProtocolService)
	hp.dh = dfh
	return hp
}
func (this *HandshakeProtocolService) GenerateKeyPair() ([]byte, *rsa.PrivateKey, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate key pair in handshake protocol ")
	}
	publicKey := &privateKey.PublicKey
	hpPublicKey := types.PubKey{
		PubKey: *publicKey,
	}
	pubKeyBytes, err := json.Marshal(hpPublicKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to serialize public key to bytes in handshake protocol ")
	}
	return pubKeyBytes, privateKey, nil
}
func (this *HandshakeProtocolService) StartDiffieHellman() (types.DFHCoefficients, error) {

	g := this.dh.Generate_g()
	n, err := this.dh.Generate_n()
	if err != nil {
		return types.DFHCoefficients{}, err
	}
	privateVariable, err := this.dh.Genrate_Private_Variable()
	if err != nil {
		return types.DFHCoefficients{}, err
	}
	dfhParams := this.generateDFHPublicKey(g, privateVariable, n)
	return dfhParams, nil
}
func (this *HandshakeProtocolService) generateDFHPublicKey(prime uint64, exponent *big.Int, modulo *big.Int) types.DFHCoefficients {

	g := new(big.Int).SetUint64(prime)
	dfhPublicKey := new(big.Int)
	dfhPublicKey.Exp(g, exponent, modulo)
	fmt.Println("dfh public key is ", dfhPublicKey)
	dfhPubKeyBytes := dfhPublicKey.Bytes()
	pubKeyEncoded := base64.StdEncoding.EncodeToString(dfhPubKeyBytes)

	//value := pubKeyEncoded

	dfhParams := types.DFHCoefficients{
		G:              prime,
		N:              modulo,
		PublicVariable: pubKeyEncoded,
	}
	return dfhParams
}
