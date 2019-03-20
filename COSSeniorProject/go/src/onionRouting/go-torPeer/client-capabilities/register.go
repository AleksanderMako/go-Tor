package clientcapabilities

import (
	"encoding/json"
	"fmt"
	"onionRouting/go-torPeer/client-capabilities/request"
	"onionRouting/go-torPeer/types"
	"os"

	"github.com/pkg/errors"
)

func RegisterPeer() error {

	peerID := GetPeerAddress()
	peerAdd := types.Register{
		PeerID: peerID,
	}
	peerAddBytes, err := json.Marshal(peerAdd)
	if err != nil {
		return errors.Wrap(err, "failed serialize peer id in registration ")
	}
	resp, err := request.Dial("http://127.0.0.1:4500/peer/", peerAddBytes)
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
	return os.Getenv("PEER_ADD")

}
