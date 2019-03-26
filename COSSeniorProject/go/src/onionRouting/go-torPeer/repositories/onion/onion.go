package onionrepository

import (
	"encoding/json"
	"onionRouting/go-torPeer/client-capabilities/request"
	storageserviceinterface "onionRouting/go-torPeer/services/storage/storage-interface"
	"onionRouting/go-torPeer/types"

	logger "github.com/apsdehal/go-logger"
	"github.com/dgraph-io/badger"

	"github.com/pkg/errors"
)

type OnionRepository struct {
	db storageserviceinterface.StorageService
}

func NewOnionRepository(db storageserviceinterface.StorageService) OnionRepository {

	onionRepo := new(OnionRepository)
	onionRepo.db = db
	return *onionRepo
}

func (this *OnionRepository) SaveCircuitLink(cID []byte, link types.CircuitLinkParameters) error {

	savedLinkBytes, err := this.db.Get(string(cID))
	if err != nil && err != badger.ErrKeyNotFound {
		return errors.Wrap(err, "failed to lookup savedLinkBytes ")
	}
	linkBytes, e := json.Marshal(link)
	if e != nil {
		return errors.Wrap(e, "failed to marshal link to bytes in SaveCircuitLink method ")
	}
	if savedLinkBytes == nil {
		if err := this.db.Put(string(cID), linkBytes); err != nil {
			return errors.Wrap(err, "failed to save link in SaveCircuitLink method")
		}
		return nil
	}

	savedLink := types.CircuitLinkParameters{}
	if e := json.Unmarshal(savedLinkBytes, &savedLink); e != nil {
		return errors.Wrap(e, "failed to unmarshal saved bytes in SaveCircuitLink")
	}
	savedLink.Previous = link.Previous
	savedLink.Next = link.Next
	newSavedLinkBytes, e := json.Marshal(savedLink)
	if e != nil {
		return errors.Wrap(e, "failed to marshal newSavedLink")
	}
	if e = this.db.Put(string(cID), newSavedLinkBytes); e != nil {
		return errors.Wrap(e, "failed to save link in badger ")
	}
	return nil
}
func (this *OnionRepository) GetCircuitLinkParamaters(cID []byte, log *logger.Logger) (types.CircuitLinkParameters, error) {

	log.Debug("GetCircuitLinkParamaters activated")
	savedLinkBytes, e := this.db.Get(string(cID))
	if e != nil {
		return types.CircuitLinkParameters{}, errors.Wrap(e, "failed to get savedLinkBytes from badger")
	}
	savedLink := types.CircuitLinkParameters{}
	if e = json.Unmarshal(savedLinkBytes, &savedLink); e != nil {
		return types.CircuitLinkParameters{}, errors.Wrap(e, "failed to get saved link in onion repository ")
	}
	log.Debug("exited all ops ")
	return savedLink, nil
}
func (this *OnionRepository) DialNext(cID []byte, next string, peeledData []byte, log *logger.Logger) error {

	log.Debug("entered dial next ")
	circuitPayload := types.CircuitPayload{
		ID:      cID,
		Payload: peeledData,
	}
	circuitPayloadBytes, e := json.Marshal(circuitPayload)
	if e != nil {
		return errors.Wrap(e, "failed to marshal circuitPayload during DialNext operation ")
	}
	req := types.Request{
		Action: "relay",
		Data:   circuitPayloadBytes,
	}
	reqBytes, e := json.Marshal(req)
	if e != nil {
		return errors.Wrap(e, "failed to marshal reqBytes during DialNext")
	}
	resp, e := request.Dial(next, reqBytes)
	if e != nil {
		return errors.Wrap(e, "failed to dial next during DialNext operation ")
	}
	log.Debugf("peer requests result %v \n", resp)
	if resp.StatusCode != 200 {
		return errors.New(resp.Status + "during DialNext where next is " + next)
	}
	body, _ := request.ParseResponse(resp)
	log.Debugf("peer reqeust body %v \n", string(body))
	return nil
}
