package main

import (
	"encoding/json"

	"io/ioutil"
	"net/http"
	controller "onionRouting/go-torPeer/controllers"
	"onionRouting/go-torPeer/types"

	"github.com/pkg/errors"
)

type Multiplexer struct {
	handShakeController controller.HandShakeController
	onionController     controller.OnionController
}

func NewMultiplexer(handShakeController controller.HandShakeController,
	OnionRoutingController controller.OnionController) Multiplexer {

	multiplexer := new(Multiplexer)
	multiplexer.handShakeController = handShakeController
	multiplexer.onionController = OnionRoutingController
	return *multiplexer
}
func (this *Multiplexer) setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (this *Multiplexer) parseIncomingRequest(r *http.Request) (string, []byte, error) {

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
func (this *Multiplexer) MultiplexRequest(w http.ResponseWriter, r *http.Request) {

	this.setupResponse(&w, r)

	action, data, err := this.parseIncomingRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var handlerErr error
	var resp []byte
	switch action {

	case "handleHandshake":
		resp, handlerErr = this.handShakeController.HandleHandshake(data)
		w.Write(resp)

	case "keyExchange":
		resp, handlerErr = this.handShakeController.HandleKeyExchange(data)
		w.Write(resp)

	case "buildCircuit":
		handlerErr = this.onionController.SaveCircuit(data)
		w.Write([]byte("added circuit link"))

	case "relay":
		handlerErr = this.onionController.RelayMessage(data)
		w.Write([]byte("message received "))
	case "backPropagate":
		handlerErr = this.onionController.BackPropagate(data)
		w.Write([]byte("successfully back propagated message "))

	}

	if handlerErr != nil {
		http.Error(w, handlerErr.Error(), http.StatusInternalServerError)
		return
	}
}
