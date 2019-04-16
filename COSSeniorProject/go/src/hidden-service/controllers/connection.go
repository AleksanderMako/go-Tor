package hiddenservicecontrollers

import (
	"encoding/json"
	"fmt"
	messagerepository "hidden-service/repositories/message"
	contentservice "hidden-service/services/content"
	serviceTypes "hidden-service/types"
	onionlib "onionLib/lib/lib-implementation"
	"onionLib/types"

	"github.com/pkg/errors"
)

type ConnectionController struct {
	onionLibrary   onionlib.OnionLibrary
	messageRepo    messagerepository.MessageRepository
	contentService contentservice.ContentService
	workDir        string
}

func NewConnectionController(onionLibrary onionlib.OnionLibrary, messageRepo messagerepository.MessageRepository, contentService contentservice.ContentService,
	workDir string) ConnectionController {

	return ConnectionController{
		onionLibrary:   onionLibrary,
		messageRepo:    messageRepo,
		contentService: contentService,
		workDir:        workDir,
	}
}
func (this *ConnectionController) parseMessage(data []byte) (types.CircuitPayload, error) {
	circuitPayload := types.CircuitPayload{}
	if err := json.Unmarshal(data, &circuitPayload); err != nil {
		return types.CircuitPayload{}, errors.Wrap(err, "failed to unmarshal incoming payload to circuit payload ")

	}

	fmt.Println("data is %v", string(data))

	return circuitPayload, nil
}
func (this *ConnectionController) TestMessage(publicKey []byte, data []byte) ([]byte, error) {

	circuitPayload, err := this.parseMessage(data)
	if err != nil {
		fmt.Println("an error has ocured ", err.Error())
		return nil, errors.Wrap(err, "failed to get message in hidden service ")
	}
	decrypted, err := this.onionLibrary.Onionservice.DeonionizeMessage(publicKey, circuitPayload.Payload)
	if err != nil {
		fmt.Println("hidden service :", err.Error())
		return nil, err
	}
	decryptedMessage, err := this.messageRepo.GetMessage(decrypted)
	if err != nil {
		fmt.Println("an error has ocured ", err.Error())

		return nil, errors.Wrap(err, "failed get message in hidden service ")
	}
	fmt.Println("hidden service says : ")
	fmt.Println(decryptedMessage)
	fmt.Println("sending circuit : ")
	fmt.Println(circuitPayload.Sender)
	hiddenResponse := serviceTypes.HiddenResponse{
		Data: []byte("successfully contacted hidden service"),
		ID:   circuitPayload.Sender,
	}

	var resp []byte
	switch decryptedMessage.Action {

	case "connect":
		resp, err = this.connect(publicKey, hiddenResponse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to hidden service ")
		}
	case "txt":
		resp, err = this.serveTextFile(publicKey, this.workDir, hiddenResponse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to hidden service ")
		}
	}
	return resp, nil
}
func (this *ConnectionController) connect(publicKey []byte, hiddenResponse serviceTypes.HiddenResponse) ([]byte, error) {

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
func (this *ConnectionController) serveTextFile(publicKey []byte, workDir string, hiddenResponse serviceTypes.HiddenResponse) ([]byte, error) {

	encryptedData, err := this.contentService.ServerTextFile(publicKey, workDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serve text file ")
	}

	hiddenResponse.Data = encryptedData
	respBytes, err := json.Marshal(hiddenResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal hidden response in connection controller ")
	}
	return respBytes, nil

}
