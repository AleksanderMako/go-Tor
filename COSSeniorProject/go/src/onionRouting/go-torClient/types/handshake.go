package types

type HandshakePayload struct {
	DFH       DFHCoefficients `json:"dfh"`
	PublicKey []byte          `json:"pubKey"`
}
