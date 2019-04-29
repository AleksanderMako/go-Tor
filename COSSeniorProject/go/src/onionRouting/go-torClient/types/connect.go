package types

type Connect struct {
	Ip           string `json:"introductionPoint"`
	DescriptorID string `json:"descriptorID"`
	Keyword      string `json:"keyWord"`
}
