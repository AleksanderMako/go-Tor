package peeronionprotocol

import (
	"fmt"
	onionrepository "onionRouting/go-torPeer/repositories/onion"
	"onionRouting/go-torPeer/types"
)

type PeerOnionService struct {
	onionRepo onionrepository.OnionRepository
}

func NewOnionService(onionRepo onionrepository.OnionRepository) PeerOnionService {
	return PeerOnionService{
		onionRepo: onionRepo,
	}
}
func (this *PeerOnionService) SaveCircuit(circuit types.P2PBuildCircuitRequest) error {

	linkParamaeters := types.CircuitLinkParameters{
		Next:     circuit.Next,
		Previous: circuit.Previous,
	}
	fmt.Println(linkParamaeters.Next)
	if err := this.onionRepo.SaveCircuitLink(circuit.ID, linkParamaeters); err != nil {
		return err
	}

	return nil
}
func (this *PeerOnionService) GetSavedCircuit(cId []byte) (CircuitLinkGetDTO, error) {

	link, err := this.onionRepo.GetCircuitLinkParamaters(cId)
	if err != nil {
		return CircuitLinkGetDTO{}, err
	}
	savedLink := CircuitLinkGetDTO{
		ID:           cId,
		Next:         link.Next,
		Previous:     link.Previous,
		SharedSecret: link.SharedSecret,
	}
	return savedLink, nil
}
