package types

type Connect struct {
	Ip           string `json:"introductionPoint"`
	DescriptorID []byte `json:"descriptorID"`
}
