package introductionpointservice

import (
	"encoding/json"
	"fmt"
	servicedescriptorrepository "hidden-service/repositories/service-descriptor"
	clientservice "hidden-service/services/client"
	"hidden-service/types"
	"math/rand"
	onionlib "onionLib/lib/lib-implementation"
	"os"
	"time"

	logger "github.com/apsdehal/go-logger"
	"github.com/pkg/errors"
)

type IntroductionService struct {
	cs           clientservice.ClientService
	descriptor   servicedescriptorrepository.ServiceDescriptorRepository
	onionLibrary onionlib.OnionLibrary
}

func NewIntroductionService(cs clientservice.ClientService,
	descriptor servicedescriptorrepository.ServiceDescriptorRepository,
	onionLibrary onionlib.OnionLibrary) IntroductionService {
	return IntroductionService{
		cs:           cs,
		descriptor:   descriptor,
		onionLibrary: onionLibrary,
	}
}

// PublishServiceDescriptor must be supplied the raw bytes of the pub key, marshalled version of types.PubKey
func (this *IntroductionService) PublishServiceDescriptor(publicKey []byte, logg *logger.Logger) error {

	pAdd, err := this.cs.GetPeersAddresses()
	if err != nil {
		return errors.Wrap(err, "failed to get peer addresses during PublishServiceDescriptor operation")
	}
	logg.Debugf("peers %v", len(pAdd))
	ips := []string{}
	//for i := 0; i < 2; i++ {
	ip1, i := this.choseRandomPeer(pAdd)
	pAdd = append(pAdd[:i], pAdd[i+1:]...)
	ips = append(ips, ip1)
	//}
	// save descriptor in db
	serviceDescriptor := types.ServiceDescriptor{
		ID:                 publicKey,
		IntroductionPoints: ips,
		KeyWords:           []string{"testing"},
	}
	err = this.descriptor.Save(serviceDescriptor)
	if err != nil {
		return errors.Wrap(err, "failed to save service descriptor during PublishServiceDescriptor")
	}
	descriptorBytes, err := json.Marshal(serviceDescriptor)
	if err != nil {
		return errors.Wrap(err, "failed to marshal serviceDescriptor to bytes  ")
	}
	// publish descriptor in api
	// publish taken peers
	url := "http://registry:4500/api/service"
	logg.Notice("ready to make request for " + url)

	resp, err := this.cs.Dial(url, descriptorBytes)
	if err != nil {
		return errors.Wrap(err, "failed to dial "+url)
	}
	body, err := this.cs.ParseResponse(resp)
	if err != nil {
		return errors.Wrap(err, "failed to parse response body ")
	}
	logg.Debugf("request body: %v \n", string(body))
	return nil
	// dial get peers

}
func (this *IntroductionService) BuildIPCircuit(publicKey []byte, privateKey types.PrivateKey, logg *logger.Logger) error {

	peerList, err := this.onionLibrary.Onionservice.GetPeers()
	if err != nil {
		return errors.Wrap(err, "failed to get peerList during BuildIPCircuit operation ")
	}
	for _, peerID := range peerList {
		logg.Debug(peerID)
	}
	descriptor, err := this.descriptor.Get(publicKey)
	if err != nil {
		return errors.Wrap(err, "failed to get descriptor during BuildIPCircuit operation  ")
	}
	destination := descriptor.IntroductionPoints[0]
	ip := descriptor.IntroductionPoints[0]
	peerList = append(peerList, ip)

	chainID, err := this.onionLibrary.Onionservice.CreateOnionChain(peerList, publicKey)
	if err != nil {
		return errors.Wrap(err, "failed to create onion chain during BuildIPCircuit operation")
	}
	privateKeyBytes, err := json.Marshal(privateKey)
	if err != nil {
		fmt.Println("error while Marshaling private key in client  ", err.Error())
		os.Exit(1)
	}
	logg.Debugf("chain id %v", chainID)
	if err := this.onionLibrary.Onionservice.HandshakeWithPeers(chainID, publicKey, privateKeyBytes); err != nil {
		return errors.Wrap(err, "failed to handshake with peers during BuildIPCircuit ")
	}
	if err := this.onionLibrary.Onionservice.GenerateSymetricKeys(chainID); err != nil {
		return errors.Wrap(err, "failed Generate SymmetricKeys during BuildIPCircuit operation")
	}

	hiddenServiceController := "hiddenservice:5000"
	if err = this.onionLibrary.Onionservice.BuildP2PCircuit([]byte(chainID), hiddenServiceController, destination); err != nil {
		return errors.Wrap(err, "failed to build p2p circuit during BuildIPCircuit operation")
	}
	return nil
}
func (this *IntroductionService) choseRandomPeer(peerList []string) (string, int) {
	rand.Seed(time.Now().Unix())
	if len(peerList) == 1 {
		return peerList[0], 0
	}
	n := rand.Int() % (len(peerList) - 1)
	return peerList[n], n
}
