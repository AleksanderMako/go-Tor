package cryptoserviceinterface

import (
	"crypto/rsa"
	"onionRouting/go-torPeer/types"
)

type CryptoService interface {
	Sign(data []byte, privKey *rsa.PrivateKey) ([]byte, error)
	Verify(data []byte, signature []byte, publicKey types.PubKey) error
	Encrypt(data []byte, key string) ([]byte, error)
	Decrypt(data []byte, key string) ([]byte, error)
	GetEncryptionKey(key string) ([]byte, error)
}
