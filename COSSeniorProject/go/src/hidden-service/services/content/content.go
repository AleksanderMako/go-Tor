package contentservice

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	onionlib "onionLib/lib/lib-implementation"
	"os"

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
func (this *ContentService) ServeImage(publicKey []byte, workingDir string) ([]byte, error) {

	path := workingDir + "servables/onion.jpg"
	image, err := os.Open(path)
	if err != nil {
		fmt.Println("error in ServeImage " + err.Error())
		return nil, errors.Wrap(err, "failed to open jpg file ")
	}
	defer image.Close()

	imgData, err := jpeg.Decode(image)
	if err != nil {
		fmt.Println("error in ServeImage " + err.Error())
		return nil, errors.Wrap(err, "failed to decode jpeg image")
	}

	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, imgData, nil)
	if err != nil {
		fmt.Println("error in ServeImage " + err.Error())
		return nil, errors.Wrap(err, "failed to make image buffer")
	}
	encrypted, err := this.onionLibrary.Onionservice.ApplyOnionLayers(publicKey, buffer.Bytes())
	if err != nil {
		fmt.Println("error in ServeImage " + err.Error())

		return nil, errors.Wrap(err, "failed to encrypt text file contentents")
	}
	return encrypted, nil
}
