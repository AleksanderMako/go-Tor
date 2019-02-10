package cryptoserviceinterface

import (
	"crypto/rsa"
	"onionRouting/go-torPeer/types"
)

type CryptoService interface {
	Sign(data []byte, privKey *rsa.PrivateKey) ([]byte, error)
	Verify(data []byte, signature []byte, publicKey types.PubKey) error
}
