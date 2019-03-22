package onionprotocol

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	handshakeprotocolservice "onionRouting/go-torClient/services/handshake"
	"onionRouting/go-torClient/services/request"
	storageserviceinterface "onionRouting/go-torClient/services/storage/storage-interface"
	"onionRouting/go-torClient/types"
	"time"

	"github.com/dgraph-io/badger"

	"github.com/pkg/errors"
)

/*
  1 query registry for 2 peers
  2 perform handshake with each peer
  3 build circuit chain
	3.1 generate id for the entire chain
	3.2 send to first peer : client address | his address | next
*/

type OnionService struct {
	storage           storageserviceinterface.StorageService
	dbVolume          *badger.DB
	handshakeProtocol handshakeprotocolservice.HandshakeProtocolService
}

func NewOnionService(storage storageserviceinterface.StorageService, db *badger.DB,
	handshakeProtocol handshakeprotocolservice.HandshakeProtocolService) OnionService {

	onionService := new(OnionService)
	onionService.storage = storage
	onionService.dbVolume = db
	onionService.handshakeProtocol = handshakeProtocol
	return *onionService
}

func (this *OnionService) GetPeers() ([]string, error) {

	url := "http://registry:4500/peer/peers"
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial given url")
	}
	body, err := request.ParseResponse(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read body for given url")
	}
	//	fmt.Println(string(body))
	var peers types.PeersDTO
	err = json.Unmarshal(body, &peers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal payload in onion service ")
	}
	return peers.Peers, nil
}

func (this *OnionService) CreateOnionChain(peerList []string) (string, error) {
	circuit := types.Circuit{}

	for i := 0; i < 3; i++ {
		//fmt.Println(peerList[i])
		peer, u := this.choseRandomPeer(peerList)
		circuit.PeerList = append(circuit.PeerList, peer)
		peerList = append(peerList[:u], peerList[u+1:]...)
	}
	timeStamp := time.Now().Format(time.RFC3339)
	hash, err := this.createHash([]byte(timeStamp))
	if err != nil {
		return "", errors.Wrap(err, "failed to hash timestamp ")
	}
	circuit.CID = hash

	circuitBytes, err := json.Marshal(circuit)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal circuit ")
	}

	err = this.storage.Put(string(hash), circuitBytes, this.dbVolume)
	if err != nil {
		return "", errors.Wrap(err, "failed to persist circuit ")
	}

	//TODO:b64 encode this
	return string(hash), nil
}

func (this *OnionService) choseRandomPeer(peerList []string) (string, int) {
	rand.Seed(time.Now().Unix())
	if len(peerList) == 1 {
		return peerList[0], 0
	}
	n := rand.Int() % (len(peerList) - 1)
	return peerList[n], n
}

func (this *OnionService) createHash(data []byte) ([]byte, error) {

	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make hash ")
	}
	hashedData := hasher.Sum(nil)
	return hashedData, nil
}

func (this *OnionService) HandshakeWithPeers(cID string) error {

	circuitBytes, err := this.storage.Get(cID)
	if err != nil {
		return errors.Wrap(err, "failed to get circuit for given id ")
	}

	circuit := types.Circuit{}
	if err := json.Unmarshal(circuitBytes, &circuit); err != nil {
		return errors.Wrap(err, "failed to unmarshal circuit bytes")
	}

	publicKey, _, err := this.handshakeProtocol.GenerateKeyPair()
	if err != nil {
		return errors.Wrap(err, "error generating key pair ")
	}
	for _, peerID := range circuit.PeerList {

		if err := this.exchangePubKeyWithPeer(peerID, publicKey); err != nil {
			return errors.Wrap(err, "error during handshake with peer "+peerID)
		}
	}
	// Exchange public key with each peer
	return nil
}

func (this *OnionService) exchangePubKeyWithPeer(perrAddress string, clientsPubKey []byte) error {

	keyExchangeReq := types.Request{
		Action: "keyExchange",
		Data:   clientsPubKey,
	}
	url := "http://" + perrAddress + "/keyExchange"
	res, err := request.Dial(url, keyExchangeReq)
	if err != nil {
		return errors.Wrap(err, "failed to dial peer with address "+perrAddress)
	}
	serverPublicKey := types.PubKey{}
	serverPublicKeyBytes, err := request.ParseResponse(res)
	if err != nil {
		return errors.Wrap(err, "failed to parse peers public key from body of response ")
	}

	if serverPublicKeyBytes == nil {
		return errors.New("peer's public key is empty ")
	}
	if err = json.Unmarshal(serverPublicKeyBytes, &serverPublicKey); err != nil {
		return errors.Wrap(err, "failed to unmarshal peers' public key  ")

	}
	if err := this.storage.Put(perrAddress, serverPublicKeyBytes, this.dbVolume); err != nil {
		return errors.Wrap(err, "failed to exchange pubKeys in onion protocol ")
	}
	fmt.Println(serverPublicKey.PubKey)

	return nil
}
