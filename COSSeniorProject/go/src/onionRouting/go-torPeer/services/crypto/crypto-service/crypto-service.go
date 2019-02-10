package cryptoservice

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	cryptoserviceinterface "onionRouting/go-torPeer/services/crypto/crypto-service-interface"
	"onionRouting/go-torPeer/types"

	"github.com/pkg/errors"
)

type CryptoService struct {
}

func NewCryptoService() cryptoserviceinterface.CryptoService {

	cryptoService := new(CryptoService)
	return cryptoService
}

func (this *CryptoService) Sign(data []byte, privKey *rsa.PrivateKey) ([]byte, error) {

	algorithm := crypto.SHA256
	newHash := algorithm.New()
	newHash.Write(data)
	hashed := newHash.Sum(nil)
	rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed)
	return sig, this.handleErr(err, "failed to generate signature")
}

func (this *CryptoService) handleErr(err error, customErrMsg string) error {
	if err != nil {
		return errors.Wrap(err, customErrMsg)
	} else {
		return nil
	}
}

func (this *CryptoService) Verify(data []byte, signature []byte, publicKey types.PubKey) error {
	algorithm := crypto.SHA256
	newHash := algorithm.New()
	newHash.Write(data)
	hashed := newHash.Sum(nil)
	err := rsa.VerifyPKCS1v15(&publicKey.PubKey, crypto.SHA256, hashed, signature)

	if err != nil {
		return errors.Wrap(err, "failed to verify signature")
	}
	return nil
}
