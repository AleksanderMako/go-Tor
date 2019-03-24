package types

type OnionPayload struct {
	Coefficients []byte `json:"coefficients"`
	CircuitID    []byte `json:"cId"`
}
