package types

type Circuit struct {
	CID      []byte   `json:"cid"`
	PeerList []string `json:"peerList"`
}
