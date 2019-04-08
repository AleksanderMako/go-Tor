package hiddenservicecontrollers

import (
	"fmt"
	onionlib "onionLib/lib/lib-implementation"
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

	decrypted, err := this.onionLibrary.Onionservice.DeonionizeMessage(publicKey, data)
	if err != nil {
		fmt.Println("hidden service :", err.Error())
		return err
	}
	fmt.Println("hidden service says : ")
	fmt.Println(string(decrypted))
	return nil
}
