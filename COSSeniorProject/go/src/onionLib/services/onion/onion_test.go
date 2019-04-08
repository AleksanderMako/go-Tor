package onionprotocol

import (
	"encoding/json"
	"fmt"
	storage "onionRouting/go-torClient/services/storage/storage-implementation"
	"onionRouting/go-torClient/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateOnionChain_ShouldPersistedCorrectly(t *testing.T) {

	// Arrange
	// err := os.Mkdir("/testDBpath", 0777)
	// assert.NoError(t, err, "asdsa")
	// println(err.Error())
	badgerDB := storage.MountDirectory("testDBFolder")
	peerList := []string{"peer1", "peer2", "peer3", "peer4", "peer5", "peer6"}

	//Act

	onionService := NewOnionService(badgerDB)

	circuitID, err := onionService.CreateOnionChain(peerList)

	//Assert
	assert.NoError(t, err, "should not throw error when all data is valid ")

	savedData, err := badgerDB.Get(circuitID)
	assert.NoError(t, err, "should not throw error when all data is valid ")
	assert.NotNil(t, savedData, "should be nil")
	var obj types.Circuit
	_ = json.Unmarshal(savedData, obj)
	fmt.Println(savedData)

}
