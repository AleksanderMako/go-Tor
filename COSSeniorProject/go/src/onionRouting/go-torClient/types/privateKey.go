package types

import "crypto/rsa"

type PrivateKey struct {
	PrivateKey rsa.PrivateKey `json:"privateKey"`
}
