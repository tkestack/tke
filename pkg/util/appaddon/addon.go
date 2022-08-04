package appaddon

import (
	"errors"
	"fmt"
	appsv1alpha1 "github.com/clusternet/apis/apps/v1alpha1"
	clustersv1beta1 "github.com/clusternet/apis/clusters/v1beta1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	application "tkestack.io/tke/api/application/v1"

	"k8s.io/apimachinery/pkg/runtime"
	runtimeutil "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/rest"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	runtimeutil.Must(clientgoscheme.AddToScheme(scheme))
	runtimeutil.Must(application.AddToScheme(scheme))
	runtimeutil.Must(clustersv1beta1.AddToScheme(scheme))
	runtimeutil.Must(appsv1alpha1.AddToScheme(scheme))
}

func GetExtenderClient(config *rest.Config) (client.Client, error) {
	var err error

	if config == nil {
		return nil, errors.New("kube restconfig file is empty")
	}

	config.ContentConfig.ContentType = "application/json"

	if err != nil {
		return nil, fmt.Errorf("failed to get cluster rest config with err: %v", err)
	}

	extenderClient, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster client with error: %v", err)
	}

	return extenderClient, nil
}
