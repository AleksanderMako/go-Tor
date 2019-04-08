package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Storage_ShouldSetGetDataCorrectly(t *testing.T) {

	//Arrange
	key := "ke1"
	data := "this is mock data"
	badgerDB := NewStorage()

	//Act
	err := badgerDB.Put(key, []byte(data))

	//assert

	assert.NoError(t, err, "should not throw error when all data is valid ")
	savedData, err1 := badgerDB.Get(key)
	assert.NoError(t, err1, "")
	assert.Equal(t, []byte(data), savedData, "data should be equal")

}
