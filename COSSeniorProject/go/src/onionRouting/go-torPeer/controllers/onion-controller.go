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
	err := this.onionService.PeelOnionLayer(circuitPayload)
	if err != nil {
		return err
	}
	return nil
}
