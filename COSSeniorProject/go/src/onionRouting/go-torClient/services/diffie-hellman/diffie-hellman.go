package diffiehellmanservice

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"

	storage "onionRouting/go-torClient/services/storage/storage-interface"

	"github.com/kavehmz/prime"
	"github.com/pkg/errors"
)

type DiffiHellmanService struct {
	storageService storage.StorageService
}

func NewDiffieHellmanService(storageService storage.StorageService) *DiffiHellmanService {

	dfh := new(DiffiHellmanService)
	dfh.storageService = storageService
	return dfh
}
func (this *DiffiHellmanService) Generate_n() (*big.Int, error) {

	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(2000), nil).Sub(max, big.NewInt(1))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate n in diffie hellman service ")
	}
	return n, nil
}

func (this *DiffiHellmanService) Generate_g() uint64 {

	p := prime.SieveOfEratosthenes(300)
	index := len(p) - 1
	prime_g := p[index]
	return prime_g
}
func (this *DiffiHellmanService) Genrate_Private_Variable() (*big.Int, error) {

	privateVariable, err := rand.Int(rand.Reader, new(big.Int).SetUint64(2000))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private variable in diffie hellman service")
	}
	return privateVariable, nil
}
func (this *DiffiHellmanService) GenerateSharedSecret(publicVariable *big.Int, privateVariable *big.Int, modulo *big.Int) {

	shareSecret := new(big.Int)
	shareSecret.Exp(publicVariable, privateVariable, modulo)
	encoded := base64.StdEncoding.EncodeToString(shareSecret.Bytes())

	algorithm := crypto.SHA256
	newHash := algorithm.New()
	newHash.Write([]byte(encoded))
	hashed := newHash.Sum(nil)

	this.storageService.Put("testPeer", hashed)

	fmt.Println("shared secret is :", string(hashed))
}
