package clientservice

import (
	"bytes"
	"encoding/json"
	"hidden-service/types"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type ClientService struct {
	client *http.Client
}

func NewClientService(client *http.Client) ClientService {
	return ClientService{
		client: client,
	}
}
func (this *ClientService) Dial(url string, req []byte) (*http.Response, error) {

	var buff bytes.Buffer
	buff.Write(req)

	resp, err := this.client.Post(url, "application/json", &buff)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	return resp, nil
}
func (this *ClientService) ParseResponse(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse server response ")
	}

	return body, nil
}
func (this *ClientService) GetPeersAddresses() ([]string, error) {

	//TODO: extract registry url to env variable
	url := "http://registry:4500/peer/peers"
	resp, err := this.client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make get request for given url: "+url+" during GetPeersAddresses operation in client service")
	}
	body, err := this.ParseResponse(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get peers in client service ")
	}
	peersDTO := types.PeersDTO{}
	if err = json.Unmarshal(body, &peersDTO); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal peers to peersDTO")
	}
	return peersDTO.Peers, nil
}
