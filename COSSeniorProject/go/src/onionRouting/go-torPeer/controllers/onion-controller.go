package controller

import (
	"encoding/json"
	"fmt"
	peeronionprotocol "onionRouting/go-torPeer/services/onion"
	"onionRouting/go-torPeer/types"

	"github.com/pkg/errors"
)

type OnionController struct {
	onionService peeronionprotocol.PeerOnionService
}

func NewOnionCOntroller(onionService peeronionprotocol.PeerOnionService) OnionController {
	return OnionController{
		onionService: onionService,
	}
}
func (this *OnionController) SaveCircuit(data []byte) error {

	circuit := types.P2PBuildCircuitRequest{}
	if err := json.Unmarshal(data, &circuit); err != nil {
		return errors.Wrap(err, "failed to unmarshal circuit in onion controller ")
	}

	if err := this.onionService.SaveCircuit(circuit); err != nil {
		return err
	}
	return nil
}
func (this *OnionController) RelayMessage(data []byte) error {

	circuitPayload := types.CircuitPayload{}
	if err := json.Unmarshal(data, &circuitPayload); err != nil {
		return errors.Wrap(err, "failed to unmarshal circuitPayload during RelayMessage operation ")
	}
	fmt.Println("Relay message activated ")

	peeledData, next, err := this.onionService.PeelOnionLayer(circuitPayload)
	if err != nil {
		return err
	}
	forwardType := "relay"
	hasNext, err := this.onionService.Forward(peeledData, circuitPayload.ID, next, forwardType)
	if err != nil {
		return err
	}
	if !hasNext {
		// decrypted data should be the bytes of  types.PubKey
		forwardType := "backPropagate"
		circuitID, link, err := this.onionService.BackTrack(peeledData)
		if err != nil {
			return errors.Wrap(err, "failed to backtrack ")
		}
		data, previous, err := this.onionService.AddOnionLayer(peeledData, link)
		if err != nil {
			return err
		}
		_, err = this.onionService.Forward(data, circuitID, previous, forwardType)
		if err != nil {
			return errors.Wrap(err, "failed to backpropagate ")
		}

	}
	return nil
}
func (this *OnionController) BackPropagate(data []byte) error {

	circuitPayload := types.CircuitPayload{}
	if err := json.Unmarshal(data, &circuitPayload); err != nil {
		return errors.Wrap(err, "failed to unmarshal circuitPayload during RelayMessage operation ")
	}

	linkDTO, err := this.onionService.GetSavedCircuit(circuitPayload.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get saved circuit during back propagation ")
	}
	link := types.CircuitLinkParameters{
		Next:         linkDTO.Next,
		Previous:     linkDTO.Previous,
		SharedSecret: linkDTO.SharedSecret,
	}
	data, previous, err := this.onionService.AddOnionLayer(circuitPayload.Payload, link)
	if err != nil {
		return err
	}
	forwardType := "backPropagate"
	_, err = this.onionService.Forward(data, circuitPayload.ID, previous, forwardType)
	if err != nil {
		return errors.Wrap(err, "failed to send payload during back propagation in peer ")
	}
	return nil

}
