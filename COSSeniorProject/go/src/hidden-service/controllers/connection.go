package hiddenservicecontrollers

import (
	"encoding/json"
	"fmt"
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
func (this *ConnectionController) TestMessage(publicKey []byte, data []byte) error {

	circuitPayload := types.CircuitPayload{}
	if err := json.Unmarshal(data, &circuitPayload); err != nil {
		return errors.Wrap(err, "failed to unmarshal incoming payload to circuit payload ")
	}
	fmt.Println("data is %v", string(data))
	decrypted, err := this.onionLibrary.Onionservice.DeonionizeMessage(publicKey, circuitPayload.Payload)
	if err != nil {
		fmt.Println("hidden service :", err.Error())
		return err
	}
	fmt.Println("hidden service says : ")
	fmt.Println(string(decrypted))
	return nil
}
