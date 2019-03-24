package peeronionprotocol

type CircuitLinkGetDTO struct {
	Next         string `json:"next"`
	Previous     string `json:"previous"`
	SharedSecret []byte `json:"shareSecret"`
	ID           []byte `json:"id"`
}
