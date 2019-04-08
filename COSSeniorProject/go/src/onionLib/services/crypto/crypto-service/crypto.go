package cryptoservice

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	cryptointerface "onionLib/services/crypto/crypto-interface"
	storageserviceinterface "onionLib/services/storage/storage-interface"
	"onionLib/types"

	"github.com/pkg/errors"
)

type CryptoService struct {
	storageService storageserviceinterface.StorageService
}

func NewCryptoService(storageService storageserviceinterface.StorageService) cryptointerface.CryptoService {

	cryptoService := new(CryptoService)
	cryptoService.storageService = storageService
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
func (this *CryptoService) Encrypt(data []byte, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	fmt.Println("new cypher started with key : " + string(key) + "\n")
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate aes cypher in crypto service")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate gcm in crypto service ")
	}
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	cypherText := gcm.Seal(nonce, nonce, data, nil)
	return cypherText, nil
}
func (this *CryptoService) Decrypt(data []byte, key []byte) ([]byte, error) {

	//secret is b64 first and the sha256 hashed
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate aes cypher in crypto service")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate gcm in crypto service ")
	}
	nonceSize := gcm.NonceSize()

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt data ")
	}
	return plainText, nil

}
func (this *CryptoService) GetEncryptionKey(key string) ([]byte, error) {

	data, err := this.storageService.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get encryption key in crypto service of clinet ")
	}
	return data, nil
}
