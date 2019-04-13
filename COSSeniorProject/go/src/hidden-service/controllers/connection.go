package hiddenservicecontrollers

import (
	"encoding/json"
	"fmt"
	serviceTypes "hidden-service/types"
	onionlib "onionLib/lib/lib-implementation"
	"onionLib/types"

	"github.com/pkg/errors"
)

type ConnectionController struct {
	onionLibrary onionlib.OnionLibrary
}

func NewConnectionController(onionLibrary onionlib.OnionLibrary) ConnectionController {

	return ConnectionController{
		onionLibrary: onionLibrary,
	}
}
func (this *ConnectionController) TestMessage(publicKey []byte, data []byte) ([]byte, error) {

	circuitPayload := types.CircuitPayload{}
	if err := json.Unmarshal(data, &circuitPayload); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal incoming payload to circuit payload ")
	}
	fmt.Println("data is %v", string(data))
	decrypted, err := this.onionLibrary.Onionservice.DeonionizeMessage(publicKey, circuitPayload.Payload)
	if err != nil {
		fmt.Println("hidden service :", err.Error())
		return nil, err
	}
	fmt.Println("hidden service says : ")
	fmt.Println(string(decrypted))
	fmt.Println("sending circuit : ")
	fmt.Println(circuitPayload.Sender)
	hiddenResponse := serviceTypes.HiddenResponse{
		Data: []byte("successfully contacted hidden service"),
		ID:   circuitPayload.Sender,
	}
	encryptedData, err := this.onionLibrary.Onionservice.ApplyOnionLayers(publicKey, hiddenResponse.Data)
	if err != nil {
		return nil, err
	}
	hiddenResponse.Data = encryptedData
	respBytes, err := json.Marshal(hiddenResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal hidden response in connection controller ")
	}
	return respBytes, nil
}
