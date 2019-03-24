package onionprotocol

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	circuitrepository "onionRouting/go-torClient/repositories/circuit"
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
	cr                circuitrepository.CircuitRepository
}

func NewOnionService(storage storageserviceinterface.StorageService, db *badger.DB,
	handshakeProtocol handshakeprotocolservice.HandshakeProtocolService,
	cr circuitrepository.CircuitRepository) OnionService {

	onionService := new(OnionService)
	onionService.storage = storage
	onionService.dbVolume = db
	onionService.handshakeProtocol = handshakeProtocol
	onionService.cr = cr
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

	err = this.cr.Save(string(hash), circuit, this.dbVolume)
	if err != nil {
		return "", err
	}
	// circuitBytes, err := json.Marshal(circuit)
	// if err != nil {
	// 	return "", errors.Wrap(err, "failed to marshal circuit ")
	// }

	// err = this.storage.Put(string(hash), circuitBytes, this.dbVolume)
	// if err != nil {
	// 	return "", errors.Wrap(err, "failed to persist circuit ")
	// }

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

	fmt.Println("entering peer handshake ")
	if cID == "" {
		return errors.New("empty cID")
	}
	circuit, err := this.cr.Get(cID)
	if err != nil {
		return err
	}
	// circuitBytes, err := this.storage.Get(cID)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to get circuit for given id ")
	// }
	// if circuitBytes == nil {
	// 	return errors.New("No peers found for the given circuit id ")
	// }
	// circuit := types.Circuit{}
	// if err := json.Unmarshal(circuitBytes, &circuit); err != nil {
	// 	return errors.Wrap(err, "failed to unmarshal circuit bytes")
	// }

	publicKey, privateKey, err := this.handshakeProtocol.GenerateKeyPair()
	if err != nil {
		return errors.Wrap(err, "error generating key pair ")
	}
	privateKeyBytes, err := json.Marshal(types.PrivateKey{
		PrivateKey: *privateKey,
	})

	if err != nil {
		return errors.Wrap(err, "failed to marshal private key to bytes during handshake")
	}
	if err := this.storage.Put("client", privateKeyBytes, this.dbVolume); err != nil {
		return errors.Wrap(err, "failed to persist clients privateKey")
	}
	for _, peerID := range circuit.PeerList {

		if err := this.exchangePubKeyWithPeer(peerID, publicKey); err != nil {
			return errors.Wrap(err, "error during handshake with peer "+peerID)
		}
	}
	// Exchange public key with each peer
	return nil
}

func (this *OnionService) exchangePubKeyWithPeer(peerAddress string, clientsPubKey []byte) error {

	keyExchangeReq := types.Request{
		Action: "keyExchange",
		Data:   clientsPubKey,
	}
	url := "http://" + peerAddress + "/"
	fmt.Println("client making request to peer " + url)
	res, err := request.Dial(url, keyExchangeReq)
	if err != nil {
		return errors.Wrap(err, "failed to dial peer with address "+peerAddress)
	}
	if res.StatusCode != 200 {
		return errors.New("request error " + res.Status)
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
	peerCredentials := types.PeerCredentials{
		PublicKey: serverPublicKeyBytes,
	}
	credentialsBytes, err := json.Marshal(peerCredentials)
	if err != nil {
		return errors.Wrap(err, "failed to marshal peer credentials in exchangePubKeyWithPeer ")
	}
	if err := this.storage.Put(peerAddress, credentialsBytes, this.dbVolume); err != nil {
		return errors.Wrap(err, "failed to exchange pubKeys in onion protocol ")
	}
	fmt.Println(serverPublicKey.PubKey)

	return nil
}
func (this *OnionService) GenerateSymetricKeys(cID string) error {

	var privateKey types.PrivateKey

	fmt.Println("entering symmetric key generation with peers ")
	privateKeyBytes, err := this.storage.Get("client")
	if err != nil {
		return errors.Wrap(err, "failed to read clients private key from persistance")
	}
	if err := json.Unmarshal(privateKeyBytes, &privateKey); err != nil {
		return errors.Wrap(err, "failed to unmarshal clients private during GenrateSymmetricKey ")
	}

	circuit, err := this.cr.Get(cID)
	if err != nil {
		return errors.Wrap(err, "failed to read circuit bytes from persistance")
	}
	// circuitBytes, err := this.storage.Get(cID)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to read circuit bytes from persistance")
	// }
	// circuit := types.Circuit{}
	// if err := json.Unmarshal(circuitBytes, &circuit); err != nil {

	// 	return errors.Wrap(err, "failed to unmarshal circuit bytes during GenerateSymetricKeys operations ")
	// }

	for _, pID := range circuit.PeerList {

		//this.storage.Get(pID)
		fmt.Println("symetric key generation peer id " + pID)
		dfCoefficients, err := this.handshakeProtocol.StartDiffieHellman(&privateKey.PrivateKey)
		if err != nil {
			return errors.Wrap(err, "failed to generate dfCoefficients for exchange with peer: "+pID)
		}
		if err := this.createShareSecret(dfCoefficients, pID, cID); err != nil {
			return errors.Wrap(err, "failed to create sharesecret with peer ")
		}

	}
	return nil

}
func (this *OnionService) createShareSecret(coefficients types.DFHCoefficients, peerID string, cID string) error {

	// serialize dfh coefficients
	coefficientsBytes, err := json.Marshal(coefficients)
	if err != nil {
		return errors.Wrap(err, "failed to marshal coefficients to bytes during createShareSecret operation")
	}
	onionPayload := types.OnionPayload{
		Coefficients: coefficientsBytes,
		CircuitID:    []byte(cID),
	}
	onionPayloadBytes, err := json.Marshal(onionPayload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal onion payload")
	}
	// formulate shared secret generation request
	req := types.Request{
		Action: "handleHandshake",
		Data:   onionPayloadBytes,
	}

	url := "http://" + peerID + "/"
	fmt.Println("client making symmetric key request to peer " + url)
	res, err := request.Dial(url, req)
	fmt.Println(res)

	if err != nil {
		return errors.Wrap(err, "failed to dial "+url)
	}
	peerPublicVariable := types.PublicVariable{}
	peerPublicVariableBytes, err := request.ParseResponse(res)
	fmt.Println(string(peerPublicVariableBytes))
	if err = json.Unmarshal(peerPublicVariableBytes, &peerPublicVariable); err != nil {
		return errors.Wrap(err, "failed to unmarshal peers public variable")
	}
	// setting peers public variable
	pPublicVar := new(big.Int)
	pPublicVar.SetBytes(peerPublicVariable.Value)
	peerSignature := peerPublicVariable.Signature
	peerCredentialsBytes, err := this.storage.Get(peerID)
	if err != nil {
		return errors.Wrap(err, "failed to read peers public key")
	}
	peerCredentials := types.PeerCredentials{}
	if err := json.Unmarshal(peerCredentialsBytes, &peerCredentials); err != nil {
		return errors.Wrap(err, "failed to unmarshal peers credentials from bytes ")
	}
	peerPublicKey := types.PubKey{}
	if err := json.Unmarshal(peerCredentials.PublicKey, &peerPublicKey); err != nil {
		return errors.Wrap(err, "failed to unmarshal peers public key from bytes ")
	}

	sharedSecretBytes, err := this.handshakeProtocol.GenerateSharedSecret(pPublicVar, coefficients.N, peerSignature, peerPublicKey)
	if err != nil {
		return err
	}
	peerCredentials.SharedSecret = sharedSecretBytes
	newpeerCredentialsBytes, err := json.Marshal(peerCredentials)
	if err != nil {
		return errors.Wrap(err, "failed to marshal peer credential to bytes ")
	}
	if err := this.storage.Put(peerID, newpeerCredentialsBytes, this.dbVolume); err != nil {
		return errors.Wrap(err, "failed to persist updated peer credentials ")
	}
	return nil
}

func (this *OnionService) BuildP2PCircuit(cID []byte, destination string) error {
	//TODO:build every peers payload for now
	circuit, err := this.cr.Get(string(cID))
	if err != nil {
		return err
	}
	client := "http://torclient:8000/circuit"
	connectionNodes := []string{client}
	connectionNodes = append(connectionNodes, circuit.PeerList[0:]...)
	connectionNodes = append(connectionNodes, destination)
	var next string
	for i := 1; i < len(connectionNodes); i++ {
		if i+1 >= (len(connectionNodes)) {
			next = ""
		} else {
			next = connectionNodes[i+1]
		}
		hop := types.P2PBuildCircuitRequest{
			Previous: connectionNodes[i-1],
			Next:     next,
			ID:       cID,
		}
		err := this.sendP2PRequest(cID, hop, connectionNodes[i])
		if err != nil {
			return err
		}
	}

	return nil
}
func (this *OnionService) sendP2PRequest(cID []byte, hop types.P2PBuildCircuitRequest, pAddress string) error {

	hopBytes, err := json.Marshal(hop)
	if err != nil {
		return errors.Wrap(err, "failed to marshal hop to bytes during sendP2PRequest")
	}
	req := types.Request{
		Action: "buildCircuit",
		Data:   hopBytes,
	}

	url := "http://" + pAddress + "/circuit"
	response, err := request.Dial(url, req)
	if err != nil {
		return errors.Wrap(err, "failed to dial peer "+pAddress)
	}
	if response.StatusCode != 200 {
		return errors.New("peer failed to build circuit " + response.Status)
	}
	return nil
}

// func (this *OnionService) ForwardMessage(cID []byte, userPayload types.UserPayload) error {

// 	circuit, err := this.cr.Get(string(cID))
// 	if err != nil {
// 		return err
// 	}circuit, err := this.cr.Get(scircuit, err := this.cr.Get(string(cID))
// 	if err != nil {circuit, err := this.cr.Get(strcircuit, err := this.cr.Get(string(cID))
// 	if err != nil {
// 		return err
// 	}ing(cID))
// 	if err != nil {
// 		return err
// 	}
// 		return err
// 	}circuit, err := this.cr.Get(string(cID))
// 	if err != nil {
// 		return err
// 	}circuit, err := this.cr.Get(string(cID))
// 	if err != nil {
// 		return err
// 	}circuit, err := this.cr.Get(string(cID))
// 	if err != nil {
// 		return err
// 	}tring(cID))
// 	if err != nil {
// 		return err
// 	}
// 	pAddress := circuit.PeerList[0]
// 	startHop := "http://" + pAddress + "/forward"
// 	lastPeerHop := types.CircuitPayload{
// 		ID:          cID,
// 		PeerAddress: userPayload.Destination,
// 		Payload:     userPayload.Data,
// 	}

// 	middleHop := types.CircuitPayload{
// 		ID:cID,
// 		PeerAddress:circuit.PeerList[2],
// 		Payload
// 	}
// 	entryHop:= types.CircuitPayload{
// 		ID:cID,
// 		PeerAddress:circuit.PeerList[1],
// 	}
// 	return nil
// }
