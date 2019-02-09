package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"onionRouting/go-torClient/types"

	"github.com/pkg/errors"
)

func Dial(url string, req types.Request) (*http.Response, error) {

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request ")
	}

	var buff bytes.Buffer
	buff.Write(reqBytes)

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
