package controller

import (
	"encoding/json"
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
	return nil
	if err := this.onionService.SaveCircuit(circuit); err != nil {
		return err
	}
	return nil
}
