package request

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func Dial(url string, req []byte) (*http.Response, error) {

	var buff bytes.Buffer
	buff.Write(req)

	resp, err := http.Post(url, "application/json", &buff)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	return resp, nil
}
func ParseResponse(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse server response ")
	}

	return body, nil
}
