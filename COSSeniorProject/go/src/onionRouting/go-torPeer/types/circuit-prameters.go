package types

type CircuitLinkParameters struct {
	SharedSecret []byte `json:"shareSecret"`
	Next         string `json:"next"`
	Previous     string `json:"previous"`
}
