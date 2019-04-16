package contentservice

import (
	"io/ioutil"
	onionlib "onionLib/lib/lib-implementation"

	"github.com/pkg/errors"
)

type ContentService struct {
	onionLibrary onionlib.OnionLibrary
}

func NewContentService(onionLibrary onionlib.OnionLibrary) ContentService {

	return ContentService{
		onionLibrary: onionLibrary,
	}
}
func (this *ContentService) ServerTextFile(publicKey []byte, workingDir string) ([]byte, error) {

	path := workingDir + "servables/sample.txt"
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read data of file ")
	}
	encrypted, err := this.onionLibrary.Onionservice.ApplyOnionLayers(publicKey, fileData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt text file contentents")
	}
	return encrypted, nil
}
