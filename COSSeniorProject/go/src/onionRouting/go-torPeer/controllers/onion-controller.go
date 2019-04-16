package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	messagerepository "onionRouting/go-torClient/repositories/message"
	peeronionprotocol "onionRouting/go-torPeer/services/onion"
	"onionRouting/go-torPeer/types"

	"github.com/pkg/errors"
)

type OnionController struct {
	onionService peeronionprotocol.PeerOnionService
	messageRepo  messagerepository.MessageRepository
}

func NewOnionCOntroller(onionService peeronionprotocol.PeerOnionService, messageRepo messagerepository.MessageRepository) OnionController {
	return OnionController{
		onionService: onionService,
		messageRepo:  messageRepo,
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
func (this *OnionController) RelayMessage(data []byte) ([]byte, error) {

	circuitPayload := types.CircuitPayload{}
	if err := json.Unmarshal(data, &circuitPayload); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal circuitPayload during RelayMessage operation ")
	}
	fmt.Println("Relay message activated ")

	peeledData, next, err := this.onionService.PeelOnionLayer(circuitPayload)
	if err != nil {
		return nil, err
	}
	forwardType := "relay"
	hasNext, body, err := this.onionService.Forward(peeledData, circuitPayload.ID, next, forwardType, nil)
	if err != nil {
		return nil, err
	}
	sendingCircuit := circuitPayload.ID
	if !hasNext {
		fmt.Println("introduction point met !!!!!!! ")
		// decrypted data should be the bytes of  types.PubKey
		forwardType := "backPropagate"

		fmt.Printf("introduction point data %v \n", string(peeledData))
		message, err := this.messageRepo.GetMessage(peeledData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get message in introduction point ")
		}

		encodedPubKey := base64.StdEncoding.EncodeToString(message.Descriptorkey)

		fmt.Printf("chaind id in bakctrack is %v\n ", encodedPubKey)
		_, link, err := this.onionService.BackTrack([]byte(encodedPubKey))
		if err != nil {
			return nil, errors.Wrap(err, "failed to backtrack ")
		}
		data, previous, err := this.onionService.AddOnionLayer(peeledData, link)
		if err != nil {
			return nil, err
		}
		_, body, err = this.onionService.Forward(data, []byte(encodedPubKey), previous, forwardType, sendingCircuit)
		if err != nil {
			return nil, errors.Wrap(err, "failed to backpropagate ")
		}
		body, err = this.ResponseDecryptor(body, link.SharedSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt response in introduction point before switching circuits ")
		}
		body, err = this.HandleIPResponse(body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}
	body, err = this.HandleIPResponse(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to handle response ")
	}
	return body, nil
}
func (this *OnionController) BackPropagate(data []byte) ([]byte, error) {

	circuitPayload := types.CircuitPayload{}
	if err := json.Unmarshal(data, &circuitPayload); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal circuitPayload during RelayMessage operation ")
	}

	linkDTO, err := this.onionService.GetSavedCircuit(circuitPayload.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get saved circuit during back propagation ")
	}
	link := types.CircuitLinkParameters{
		Next:         linkDTO.Next,
		Previous:     linkDTO.Previous,
		SharedSecret: linkDTO.SharedSecret,
	}
	data, previous, err := this.onionService.AddOnionLayer(circuitPayload.Payload, link)
	if err != nil {
		return nil, err
	}
	forwardType := "backPropagate"
	_, body, err := this.onionService.Forward(data, circuitPayload.ID, previous, forwardType, circuitPayload.SenderPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send payload during back propagation in peer ")
	}
	body, err = this.ResponseDecryptor(body, linkDTO.SharedSecret)
	if err != nil {
		return nil, err
	}
	return body, nil

}
func (this *OnionController) createHash(data []byte) ([]byte, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make hash ")
	}
	hashedData := hasher.Sum(nil)
	return hashedData, nil
}

func (this *OnionController) HandleIPResponse(response []byte) ([]byte, error) {

	fmt.Printf("HandleIPResponse received resp %v\n", string(response))
	hiddenResponse := types.HiddenResponse{}
	if err := json.Unmarshal(response, &hiddenResponse); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response in HandleIPResponse")
	}
	linkDTO, err := this.onionService.GetSavedCircuit(hiddenResponse.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get saved circuit during back propagation ")
	}
	link := types.CircuitLinkParameters{
		Next:         linkDTO.Next,
		Previous:     linkDTO.Previous,
		SharedSecret: linkDTO.SharedSecret,
	}
	data, _, err := this.onionService.AddOnionLayer(hiddenResponse.Data, link)
	if err != nil {
		return nil, err
	}
	hiddenResponse.Data = data
	hiddenResponseBytes, err := json.Marshal(hiddenResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marahshal hidden response to bytes ")
	}
	return hiddenResponseBytes, nil
}
func (this *OnionController) ResponseDecryptor(response []byte, key []byte) ([]byte, error) {
	hiddenResponse := types.HiddenResponse{}
	if err := json.Unmarshal(response, &hiddenResponse); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response in HandleIPResponse")
	}
	decrypted, err := this.onionService.DecryptData(hiddenResponse.Data, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt data in ResponseDecryptor")
	}
	hiddenResponse.Data = decrypted
	hiddenResponseBts, err := json.Marshal(hiddenResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal hidden response to bytes in ResponseDecryptor")
	}
	return hiddenResponseBts, nil
}
