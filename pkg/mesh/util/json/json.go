package json

import (
	"bytes"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var (
	json                = jsoniter.ConfigCompatibleWithStandardLibrary
	MarshalToString     = json.MarshalToString
	Marshal             = json.Marshal
	MarshalIndent       = json.MarshalIndent
	Unmarshal           = json.Unmarshal
	UnmarshalFromString = json.UnmarshalFromString
	NewDecoder          = json.NewDecoder
	NewEncoder          = json.NewEncoder
	Get                 = json.Get
)

func NewJsonRequest(req *http.Request) (*JsonRequest, error) {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read the request body")
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	return &JsonRequest{
		req,
		buf,
	}, nil
}

type JsonRequest struct {
	req *http.Request
	raw []byte
}

func (j *JsonRequest) FindObject(jsonpath ...interface{}) jsoniter.Any {
	return json.Get(j.raw, jsonpath...)
}
