package types

type CircuitPayload struct {
	ID              []byte `json:"id"`
	Payload         []byte `json:"payload"`
	SenderPublicKey []byte `json:"sender"`
}
