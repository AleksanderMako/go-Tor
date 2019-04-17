package onionprotocol

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"math/rand"
	"net/http"
	circuitrepository "onionLib/repositories/circuit"
	cryptointerface "onionLib/services/crypto/crypto-interface"
	handshakeprotocolservice "onionLib/services/handshake"
	"onionLib/services/request"
	storageserviceinterface "onionLib/services/storage/storage-interface"
	"onionLib/types"
	"os"
	"time"

	peercredentialsrepository "onionLib/repositories/credentials"

	logger "github.com/apsdehal/go-logger"
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
	storage             storageserviceinterface.StorageService
	dbVolume            *badger.DB
	handshakeProtocol   handshakeprotocolservice.HandshakeProtocolService
	cr                  circuitrepository.CircuitRepository
	log                 *logger.Logger
	peerCredentialsRepo peercredentialsrepository.PeerCredentials
	cryptoService       cryptointerface.CryptoService
	PublicKey           []byte
}

func NewOnionService(storage storageserviceinterface.StorageService, db *badger.DB,
	handshakeProtocol handshakeprotocolservice.HandshakeProtocolService,
	cr circuitrepository.CircuitRepository,
	CredentialsRepo peercredentialsrepository.PeerCredentials,
	cryptoService cryptointerface.CryptoService,
	PublicKey []byte) OnionService {

	onionService := new(OnionService)
	onionService.storage = storage
	onionService.dbVolume = db
	onionService.handshakeProtocol = handshakeProtocol
	onionService.cr = cr
	onionService.peerCredentialsRepo = CredentialsRepo
	onionService.cryptoService = cryptoService
	onionService.log, _ = logger.New("OnionService", 1, os.Stdout)
	onionService.PublicKey = PublicKey
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
	var peers types.PeersDTO
	err = json.Unmarshal(body, &peers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal payload in onion service ")
	}
	return peers.Peers, nil
}
func (this *OnionService) GetServiceDescriptorsByKeyWords(keyword string) (types.ServiceDescriptorsDTO, error) {
	if keyword == "" {
		return types.ServiceDescriptorsDTO{}, errors.New("illegal empty keyword in GetServiceDescriptorsByKeyWords method  ")
	}
	ipURL := "http://registry:4500/api/hiddenservice/" + keyword
	resp, err := http.Get(ipURL)
	if err != nil {
		return types.ServiceDescriptorsDTO{}, errors.Wrap(err, "failed to dial given url")
	}
	body, err := request.ParseResponse(resp)
	if err != nil {
		return types.ServiceDescriptorsDTO{}, errors.Wrap(err, "failed to read body for given url")
	}
	this.log.Debugf("descriptor request body %v", string(body))
	var descriptors types.ServiceDescriptorsDTO
	err = json.Unmarshal(body, &descriptors)
	if err != nil {
		return types.ServiceDescriptorsDTO{}, errors.Wrap(err, "failed to unmarshal payload in onion service ")
	}
	return descriptors, nil
}
func (this *OnionService) CreateCryptoMaterials() ([]byte, types.PrivateKey, error) {

	publicKey, privateKey, err := this.handshakeProtocol.GenerateKeyPair()
	if err != nil {
		return nil, types.PrivateKey{}, errors.Wrap(err, "failed to generate cryptomaterials ")
	}
	privKey := types.PrivateKey{
		PrivateKey: *privateKey,
	}
	return publicKey, privKey, nil
}
func (this *OnionService) CreateOnionChain(peerList []string, publicKey []byte) (string, error) {
	circuit := types.Circuit{}

	ip := peerList[len(peerList)-1]
	this.log.Debugf("introduction pint is %v", ip)
	peerList = peerList[:len(peerList)-2]

	for i := 0; i < 2; i++ {
		peer, u := this.choseRandomPeer(peerList)
		circuit.PeerList = append(circuit.PeerList, peer)
		peerList = append(peerList[:u], peerList[u+1:]...)
	}
	encodedPubKey := base64.StdEncoding.EncodeToString(publicKey)
	circuit.PeerList = append(circuit.PeerList, ip)
	circuit.CID = []byte(encodedPubKey)
	this.log.Noticef("number of peers in chain is : %v", len(circuit.PeerList))
	err := this.cr.Save(encodedPubKey, circuit, this.dbVolume)
	if err != nil {
		return "", err
	}
	//TODO:b64 encode this
	return encodedPubKey, nil
}

func (this *OnionService) choseRandomPeer(peerList []string) (string, int) {
	rand.Seed(time.Now().Unix())
	if len(peerList) == 1 {
		return peerList[0], 0
	}
	n := rand.Int() % (len(peerList) - 1)
	this.log.Noticef("chosen index is %v", 1)
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

func (this *OnionService) HandshakeWithPeers(cID string, publicKey []byte, privateKeyBytes []byte) error {

	this.log.Info("entering peer handshake ")
	if cID == "" {
		return errors.New("empty cID")
	}
	circuit, err := this.cr.Get(cID)
	if err != nil {
		return err
	}
	if err := this.storage.Put("client", privateKeyBytes); err != nil {
		return errors.Wrap(err, "failed to persist clients privateKey")
	}
	for _, peerID := range circuit.PeerList {

		if err := this.exchangePubKeyWithPeer(peerID, publicKey); err != nil {
			return errors.Wrap(err, "error during handshake with peer "+peerID)
		}
	}
	return nil
}

func (this *OnionService) exchangePubKeyWithPeer(peerAddress string, clientsPubKey []byte) error {

	keyExchangeReq := types.Request{
		Action: "keyExchange",
		Data:   clientsPubKey,
	}
	url := "http://" + peerAddress + "/"
	this.log.Notice("client making request to peer " + url)
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
	if err := this.peerCredentialsRepo.SavePeerCredentials(peerAddress, peerCredentials); err != nil {
		return errors.Wrap(err, "failed to exchange pubKeys in onion protocol ")
	}
	this.log.Noticef("peers public key %v \n ", serverPublicKey.PubKey)

	return nil
}
func (this *OnionService) GenerateSymetricKeys(cID string) error {

	var privateKey types.PrivateKey

	this.log.Notice("entering symmetric key generation with peers \n")
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

	for _, pID := range circuit.PeerList {

		this.log.Notice("symmetric key generation for peer id " + pID + "\n")
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
	this.log.Noticef("client making symmetric key request to peer %v \n", url)
	res, err := request.Dial(url, req)
	this.log.Debugf(" create shared secret response %v \n", res)

	if err != nil {
		return errors.Wrap(err, "failed to dial "+url)
	}
	peerPublicVariable := types.PublicVariable{}
	peerPublicVariableBytes, err := request.ParseResponse(res)
	this.log.Debugf("peers public variable %v \n", string(peerPublicVariableBytes))
	if err = json.Unmarshal(peerPublicVariableBytes, &peerPublicVariable); err != nil {
		return errors.Wrap(err, "failed to unmarshal peers public variable")
	}
	// setting peers public variable
	pPublicVar := new(big.Int)
	pPublicVar.SetBytes(peerPublicVariable.Value)
	peerSignature := peerPublicVariable.Signature
	// peerCredentialsBytes, err := this.storage.Get(peerID)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to read peers public key")
	// }
	// peerCredentials := types.PeerCredentials{}
	// if err := json.Unmarshal(peerCredentialsBytes, &peerCredentials); err != nil {
	// 	return errors.Wrap(err, "failed to unmarshal peers credentials from bytes ")
	// }
	peerCredentials, err := this.peerCredentialsRepo.GetPeerCredentials(peerID, this.PublicKey)
	if err != nil {
		return errors.Wrap(err, "failed to read peer credentials ")
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
	if err := this.peerCredentialsRepo.SavePeerCredentials(peerID, peerCredentials); err != nil {
		return err
	}

	return nil
}

func (this *OnionService) BuildP2PCircuit(cID []byte, source string, destination string) error {
	circuit, err := this.cr.Get(string(cID))
	if err != nil {
		return err
	}
	//client := "torclient:8000"

	connectionNodes := []string{source}
	connectionNodes = append(connectionNodes, circuit.PeerList[0:]...)
	//	connectionNodes = append(connectionNodes, destination)
	var next string
	for i := 1; i < len(connectionNodes); i++ {
		if i+1 >= (len(connectionNodes)) {
			next = ""
		} else {
			next = connectionNodes[i+1]
		}
		this.log.Debugf("next is %v \n", next)
		hop := types.P2PBuildCircuitRequest{
			Previous: connectionNodes[i-1],
			Next:     next,
			ID:       cID,
		}
		this.log.Noticef("sending p2p request for peer %v \n ", connectionNodes[i])
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

	url := "http://" + pAddress + "/"
	response, err := request.Dial(url, req)
	if err != nil {
		return errors.Wrap(err, "failed to dial peer "+pAddress)
	}
	this.log.Debugf("p2p request resp %v \n", response)
	if response.StatusCode != 200 {
		return errors.New("peer failed to build circuit " + response.Status)
	}
	body, err := request.ParseResponse(response)
	if err != nil {
		return err
	}
	this.log.Debugf("body is %v \n", string(body))
	return nil
}

//TODO: modify datatype later
func (this *OnionService) SendMessage(cID []byte, data []byte) ([]byte, error) {

	circuit, err := this.cr.Get(string(cID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get circuit in SendMessage ")
	}

	encrypted, err := this.onionizeMessage(circuit.PeerList, data)
	if err != nil {
		return nil, err
	}
	circuitPayload := types.CircuitPayload{
		ID:      cID,
		Payload: encrypted,
	}
	payloadBytes, err := json.Marshal(circuitPayload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal  payload to bytes")
	}
	hop1 := "http://" + circuit.PeerList[0]
	req := types.Request{
		Action: "relay",
		Data:   payloadBytes,
	}

	resp, err := request.Dial(hop1, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial ")
	}
	this.log.Debugf("hop1 response  %v \n", resp)

	body, _ := request.ParseResponse(resp)
	hiddenResponse := types.HiddenResponse{}
	if err = json.Unmarshal(body, &hiddenResponse); err != nil {
		return nil, errors.Wrap(err, "fialed to unmarshal hidden response ")
	}
	this.log.Debugf("hop1 body  %v \n", string(body))
	decryptedResponse, err := this.DeonionizeMessage(this.PublicKey, hiddenResponse.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt hidden service response ")
	}
	hiddenResponse.Data = decryptedResponse
	this.log.Noticef(" response is : %v", string(decryptedResponse))

	return decryptedResponse, nil
}
func (this *OnionService) onionizeMessage(peerList []string, data []byte) ([]byte, error) {

	for i := len(peerList) - 1; i >= 0; i-- {
		this.log.Debugf("appliying shared secret of peer %v \n", peerList[i])
		peerCredentials, err := this.peerCredentialsRepo.GetPeerCredentials(peerList[i], []byte(data))
		if err != nil {
			return nil, err
		}
		key := peerCredentials.SharedSecret
		dataBytes, err := this.cryptoService.Encrypt(data, key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encrypt message during onionizeMessage operation")
		}
		data = dataBytes
	}
	return []byte(data), nil
}
func (this *OnionService) DeonionizeMessage(publicKey []byte, data []byte) ([]byte, error) {

	encodedPubKey := base64.StdEncoding.EncodeToString(publicKey)
	circuit, err := this.cr.Get(encodedPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deonionize message ")
	}
	peerList := circuit.PeerList
	for i := 0; i < len(peerList); i++ {

		this.log.Debugf("DeonionizeMessage by applying shared secret of peer %v \n", peerList[i])
		peerCredentials, err := this.peerCredentialsRepo.GetPeerCredentials(peerList[i], publicKey)
		if err != nil {
			return nil, err
		}
		key := peerCredentials.SharedSecret
		this.log.Debugf("DeonionizeMessage encryption key casted to string %v", string(key))
		this.log.Debugf("DeonionizeMessage encryption key raw bytes  %v", key)

		dataBytes, err := this.cryptoService.Decrypt(data, key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt message during DeonionizeMessage operation")
		}
		data = dataBytes
	}
	return data, nil
}
func (this *OnionService) ApplyOnionLayers(publicKey []byte, data []byte) ([]byte, error) {
	encodedPubKey := base64.StdEncoding.EncodeToString(publicKey)
	circuit, err := this.cr.Get(encodedPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deonionize message ")
	}
	encryptedData, err := this.onionizeMessage(circuit.PeerList, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply onion layers in ")
	}
	return encryptedData, nil
}
