package types

type Request struct {
	Action string `json:"action"`
	Data   []byte `json:"data"`
}
