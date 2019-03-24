package types

type P2PBuildCircuitRequest struct {
	ID       []byte `json:"id"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}
