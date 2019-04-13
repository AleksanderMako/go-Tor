package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"onionLib/types"
	"time"

	"github.com/pkg/errors"
)

func NewHttpClient() *http.Client {
	client := &http.Client{
		Timeout: time.Second * 300,
	}
	return client
}
func Dial(url string, req types.Request) (*http.Response, error) {

	c := NewHttpClient()
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request ")
	}

	var buff bytes.Buffer
	buff.Write(reqBytes)

	resp, err := c.Post(url, "application/json", &buff)
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
