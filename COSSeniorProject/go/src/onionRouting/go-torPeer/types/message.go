package types

type Message struct {
	Descriptorkey []byte `json:"descriptor"`
	Action        string `json:"action"`
}
