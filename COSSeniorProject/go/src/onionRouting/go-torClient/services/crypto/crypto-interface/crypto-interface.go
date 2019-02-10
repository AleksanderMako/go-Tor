package cryptointerface

import (
	"crypto/rsa"
	"onionRouting/go-torClient/types"
)

type CryptoService interface {
	Sign(data []byte, privKey *rsa.PrivateKey) ([]byte, error)
	Verify(data []byte, signature []byte, publicKey types.PubKey) error
}
