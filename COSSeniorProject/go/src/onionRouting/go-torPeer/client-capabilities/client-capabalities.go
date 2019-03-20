package clientcapabilities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"onionRouting/go-torPeer/client-capabilities/request"
	"onionRouting/go-torPeer/types"
	"os"

	"github.com/pkg/errors"
)

func RegisterPeer() error {

	peerID := GetPeerAddress()
	if peerID == "" {
		return errors.New("failed to get peer id ")
	}
	peerAdd := types.Register{
		PeerID: peerID,
	}
	peerAddBytes, err := json.Marshal(peerAdd)
	if err != nil {
		return errors.Wrap(err, "failed serialize peer id in registration ")
	}
	resp, err := request.Dial("http://registry:4500/peer/", peerAddBytes)
	if err != nil {
		return errors.Wrap(err, "failed to make request to registry")
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Status)
	res, err := request.ParseResponse(resp)
	if err != nil {
		return errors.Wrap(err, "failed to read registry's response ")
	}
	fmt.Println(string(res))
	return nil
}

func GetPeerAddress() string {
	dns := os.Getenv("DNS")
	port := os.Getenv("PEER_PORT")

	return dns + ":" + port

}
func GetPeerAddresses(url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to dial given url")
	}
	body, err := request.ParseResponse(resp)
	if err != nil {
		return errors.Wrap(err, "failed to read body for given url")
	}
	fmt.Println(string(body))
	return nil
}
