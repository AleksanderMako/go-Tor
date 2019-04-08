package main

import (
	"flag"
	"fmt"
	hiddenservicecontrollers "hidden-service/controllers"
	servicedescriptorrepository "hidden-service/repositories/service-descriptor"
	clientservice "hidden-service/services/client"
	introductionpointservice "hidden-service/services/ip"
	storage "hidden-service/services/storage/storage-implementation"
	"hidden-service/types"
	"net/http"
	onionlib "onionLib/lib/lib-implementation"
	"os"
	"time"

	logger "github.com/apsdehal/go-logger"
)

type TestHandler struct {
	PublicKey []byte

	PrivateKey types.PrivateKey
}

func main() {

	badgerDB := storage.NewStorage()
	badgerOpts := storage.InitializeBadger()
	onionLib := onionlib.NewOnionLib(badgerOpts)
	publicKey, privateKey, err := onionLib.Onionservice.CreateCryptoMaterials()
	if err != nil {
		fmt.Println("HIDDEN SERVICE ERROR ", err.Error())
		os.Exit(1)
	}
	connectionController := hiddenservicecontrollers.NewConnectionController(onionLib)
	multiPlexer := NewHiddenServiceMultiplexer(connectionController, publicKey, privateKey)
	client := NewHttpClient()
	clientService := clientservice.NewClientService(client)
	serviceDescriptorRepository := servicedescriptorrepository.NewServiceDescriptorRepository(badgerDB)
	introductionProtocol := introductionpointservice.NewIntroductionService(clientService, serviceDescriptorRepository, onionLib)
	log, _ := logger.New("HiddenService", 1, os.Stdout)
	err = introductionProtocol.PublishServiceDescriptor(publicKey, log)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	hiddenServicePrivateKey := types.PrivateKey{
		PrivateKey: privateKey.PrivateKey,
	}
	if err = introductionProtocol.BuildIPCircuit(publicKey, hiddenServicePrivateKey, log); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	http.HandleFunc("/", multiPlexer.Multiplex)
	httpAddr := flag.String("http", ":"+"5000", "Listen address")
	http.ListenAndServe(*httpAddr, nil)
}
func NewHttpClient() *http.Client {
	client := &http.Client{
		Timeout: time.Second * 120,
	}
	return client
}
