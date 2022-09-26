package page

import (
	"context"
	"encoding/base64"
	"encoding/json"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
)

type ContinueStruct struct {
	ClsName  string `json:"clsName,omitempty"`
	Resource string `json:"resource,omitempty"`
	Name     string `json:"name,omitempty"`
	Start    int64  `json:"start,omitempty"`
	Limit    int64  `json:"limit,omitempty"`
}

func EncodeContinue(ctx context.Context, resource, name string, start, limit int64) (string, error) {
	clusterName := filter.ClusterFrom(ctx)
	continueStruct := &ContinueStruct{
		ClsName:  clusterName,
		Resource: resource,
		Name:     name,
		Start:    start,
		Limit:    limit,
	}
	continueJSONByte, err := json.Marshal(continueStruct)
	if err != nil {
		return "", err
	}
	continueBase64Str := base64.StdEncoding.EncodeToString(continueJSONByte)

	return continueBase64Str, nil
}

func DecodeContinue(ctx context.Context, resource, name, continueStr string) (int64, int64, error) {
	continueJSONByte, err := base64.StdEncoding.DecodeString(continueStr)
	if err != nil {
		return 0, 0, err
	}
	continueStruct := &ContinueStruct{}
	err = json.Unmarshal(continueJSONByte, continueStruct)
	if err != nil {
		return 0, 0, err
	}
	clusterName := filter.ClusterFrom(ctx)
	if continueStruct.ClsName != clusterName || continueStruct.Name != name || continueStruct.Resource != resource {
		return 0, 0, k8serror.NewNotFound(schema.GroupResource{Group: "", Resource: "cache"}, "")
	}
	return continueStruct.Start, continueStruct.Limit, nil
}
