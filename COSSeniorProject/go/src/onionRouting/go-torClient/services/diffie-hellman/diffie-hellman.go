package diffiehellmanservice

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/kavehmz/prime"
	"github.com/pkg/errors"
)

type DiffiHellmanService struct {
}

func NewDiffieHellmanService() *DiffiHellmanService {

	dfh := new(DiffiHellmanService)
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

	algorithm := crypto.SHA256
	newHash := algorithm.New()
	newHash.Write(shareSecret.Bytes())
	hashed := newHash.Sum(nil)

	//	this.sharedSecret = hashed
	encoded := base64.StdEncoding.EncodeToString(hashed)

	fmt.Println("shared secret is :", encoded)
}
