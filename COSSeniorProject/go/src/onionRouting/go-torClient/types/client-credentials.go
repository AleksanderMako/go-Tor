package types

type PeerCredentials struct {
	PublicKey    []byte `json:"pubkey"`
	SharedSecret []byte `json:"shareSecret`
}
