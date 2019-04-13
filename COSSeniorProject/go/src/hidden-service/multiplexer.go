package main

import (
	"encoding/json"
	hiddenservicecontrollers "hidden-service/controllers"
	"io/ioutil"
	"net/http"
	"onionLib/types"

	"github.com/pkg/errors"
)

type HiddenServiceMultiplexer struct {
	connetionController hiddenservicecontrollers.ConnectionController
	PublicKey           []byte
	PrivateKey          types.PrivateKey
}

func NewHiddenServiceMultiplexer(connetionController hiddenservicecontrollers.ConnectionController,
	PublicKey []byte,
	PrivateKey types.PrivateKey) HiddenServiceMultiplexer {

	return HiddenServiceMultiplexer{
		connetionController: connetionController,
		PrivateKey:          PrivateKey,
		PublicKey:           PublicKey,
	}
}
func (this *HiddenServiceMultiplexer) parseIncomingRequest(r *http.Request) (string, []byte, error) {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to read request ")
	}
	payload := types.Request{}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		return "", nil, err
	}
	return payload.Action, payload.Data, nil
}
func (this *HiddenServiceMultiplexer) Multiplex(w http.ResponseWriter, r *http.Request) {

	_, data, err := this.parseIncomingRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := this.connetionController.TestMessage(this.PublicKey, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(response)

}
func (this *HiddenServiceMultiplexer) HandleTextFileDelivery(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("successfully contacted hidden service "))

}
