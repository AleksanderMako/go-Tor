package messagerepository

import (
	"encoding/json"
	"fmt"
	"onionRouting/go-torClient/types"

	"github.com/pkg/errors"
)

type MessageRepository struct {
}

func NewMessageRepository() MessageRepository {

	message := new(MessageRepository)
	return *message
}
func (this *MessageRepository) CreateMessage(descriptorKey []byte, action string) ([]byte, error) {

	if action == "" {
		return nil, errors.New("Illegal empty action argument supplied to CreateMessage method ")
	}
	message := types.Message{
		Action:        action,
		Descriptorkey: descriptorKey,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message to bytes in CreateMessage ")
	}
	return messageBytes, nil
}
func (this *MessageRepository) GetMessage(messageBytes []byte) (types.Message, error) {

	message := types.Message{}
	fmt.Printf("message bytes are %v\n", messageBytes)
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		return types.Message{}, errors.Wrap(err, "failed to unmarshal message bytes  ")
	}
	return message, nil
}
