package dfhservice

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	cryptoserviceinterface "onionRouting/go-torPeer/services/crypto/crypto-service-interface"
	storageserviceinterface "onionRouting/go-torPeer/services/storage/storage-interface"
	"onionRouting/go-torPeer/types"

	"github.com/pkg/errors"
)

type DFHService struct {
	cs              cryptoserviceinterface.CryptoService
	privateVariable *big.Int
	sharedSecret    []byte
	storageService  storageserviceinterface.StorageService
}

func NewDfhService(cs cryptoserviceinterface.CryptoService, storageService storageserviceinterface.StorageService) *DFHService {

	dfhService := new(DFHService)
	dfhService.cs = cs
	dfhService.storageService = storageService
	return dfhService
}

func (this *DFHService) Genrate_Private_Variable() (*big.Int, error) {

	privateVariable, err := rand.Int(rand.Reader, new(big.Int).SetUint64(2000))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private variable in diffie hellman service")
	}
	this.privateVariable = privateVariable
	return privateVariable, nil
}
func (this *DFHService) GeneratePublicVariable(prime uint64, exponent *big.Int, modulo *big.Int, privKey *rsa.PrivateKey) (types.PublicVariable, error) {
	g := new(big.Int).SetUint64(prime)
	dfhPublicKey := new(big.Int)
	dfhPublicKey.Exp(g, exponent, modulo)
	//fmt.Println("dfh public key is ", dfhPublicKey)
	dfhPubKeyBytes := dfhPublicKey.Bytes()
	sig, err := this.cs.Sign(dfhPubKeyBytes, privKey)
	if err != nil {
		return types.PublicVariable{}, errors.Wrap(err, "failed to sign peers public variable")
	}

	peerDfhPublicKey := types.PublicVariable{
		Signature: sig,
		Value:     dfhPubKeyBytes,
	}
	return peerDfhPublicKey, nil
}

func (this *DFHService) GenerateSharedSecret(publicVariable *big.Int, privateVariable *big.Int, modulo *big.Int) error {

	shareSecret := new(big.Int)
	shareSecret.Exp(publicVariable, privateVariable, modulo)
	encoded := base64.StdEncoding.EncodeToString(shareSecret.Bytes())

	algorithm := crypto.SHA256
	newHash := algorithm.New()
	newHash.Write([]byte(encoded))
	hashed := newHash.Sum(nil)

	this.sharedSecret = hashed

	err := this.storageService.Put("clientSecret", this.sharedSecret)
	if err != nil {
		return errors.Wrap(err, "failed to persist share secret in storage")
	}
	fmt.Println("shared secret is :", this.sharedSecret)
	return nil
}
