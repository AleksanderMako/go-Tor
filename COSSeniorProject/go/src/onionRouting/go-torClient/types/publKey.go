package types

import "crypto/rsa"

type PubKey struct {
	PubKey rsa.PublicKey `json:"pubkey"`
}
