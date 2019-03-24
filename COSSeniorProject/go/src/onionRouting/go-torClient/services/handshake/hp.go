package handshakeprotocolservice

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"math/big"
	cryptointerface "onionRouting/go-torClient/services/crypto/crypto-interface"
	diffiehellmanservice "onionRouting/go-torClient/services/diffie-hellman"
	"onionRouting/go-torClient/types"

	"github.com/pkg/errors"
)

type HandshakeProtocolService struct {
	dh              diffiehellmanservice.DiffiHellmanService
	cryptoService   cryptointerface.CryptoService
	privateVariable *big.Int
	sharedSecret    []byte
}

func NewHandshakeProtocol(dfh diffiehellmanservice.DiffiHellmanService,
	cryptoService cryptointerface.CryptoService) *HandshakeProtocolService {

	hp := new(HandshakeProtocolService)
	hp.dh = dfh
	hp.cryptoService = cryptoService
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
func (this *HandshakeProtocolService) StartDiffieHellman(privKey *rsa.PrivateKey) (types.DFHCoefficients, error) {

	g := this.dh.Generate_g()
	n, err := this.dh.Generate_n()
	if err != nil {
		return types.DFHCoefficients{}, err
	}
	privateVariable, err := this.dh.Genrate_Private_Variable()
	if err != nil {
		return types.DFHCoefficients{}, err
	}
	this.privateVariable = privateVariable
	dfhParams, err := this.generateDFHPublicKey(g, privateVariable, n, privKey)
	if err != nil {
		return types.DFHCoefficients{}, err
	}
	return dfhParams, nil
}
func (this *HandshakeProtocolService) generateDFHPublicKey(prime uint64, exponent *big.Int, modulo *big.Int, privKey *rsa.PrivateKey) (types.DFHCoefficients, error) {

	g := new(big.Int).SetUint64(prime)
	dfhPublicKey := new(big.Int)
	dfhPublicKey.Exp(g, exponent, modulo)
	//	fmt.Println("dfh public key is ", dfhPublicKey)
	dfhPubKeyBytes := dfhPublicKey.Bytes()
	sig, err := this.cryptoService.Sign(dfhPubKeyBytes, privKey)
	if err != nil {
		return types.DFHCoefficients{}, err
	}
	dfhParams := types.DFHCoefficients{
		G: prime,
		N: modulo,
		PublicVariable: types.PublicVariable{
			Signature: sig,
			Value:     dfhPubKeyBytes,
		},
	}
	return dfhParams, nil
}
func (this *HandshakeProtocolService) GenerateSharedSecret(publicVariable *big.Int, modulo *big.Int, signature []byte, publicKey types.PubKey) ([]byte, error) {

	if err := this.cryptoService.Verify(publicVariable.Bytes(), signature, publicKey); err != nil {

		return nil, errors.Wrap(err, "failed to verify peers public variable  ")
	}
	return this.dh.GenerateSharedSecret(publicVariable, this.privateVariable, modulo), nil
}
