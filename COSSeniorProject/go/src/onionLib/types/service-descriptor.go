package types

type ServiceDescriptor struct {
	IntroductionPoints []string `json:"ips"`
	ID                 []byte   `json:"id"`
	KeyWords           []string `json:"keyWords"`
}
