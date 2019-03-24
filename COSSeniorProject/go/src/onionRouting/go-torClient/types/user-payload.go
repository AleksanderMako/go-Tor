package types

type UserPayload struct {
	Destination string `json:"destination"`
	Data        []byte `json:"data"`
}
