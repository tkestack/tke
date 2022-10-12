package extenderapi

import (
	"context"
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

func GetManagedCluster(clientSet client.Client, name string) (*clustersv1beta1.ManagedCluster, error) {

	mcSet := clustersv1beta1.ManagedClusterList{}

	err := clientSet.List(context.TODO(), &mcSet, client.MatchingFields{"metadata.name": name})

	if err != nil {
		return nil, fmt.Errorf("fail to get managed cluster object which name is %s, err is %v", name, err)
	}

	if len(mcSet.Items) == 0 {
		return nil, fmt.Errorf("not find a managed cluster named %s", name)
	}

	// mcs should have  only one item
	mc := mcSet.Items[0]

	return &mc, nil
}

func GetSubscription(clientSet client.Client, name, namespace string) (*appsv1alpha1.Subscription, error) {
	sub := new(appsv1alpha1.Subscription)
	sub.Name = name
	sub.Namespace = namespace
	key := client.ObjectKeyFromObject(sub)
	err := clientSet.Get(context.TODO(), key, sub)
	if err != nil {
		return nil, fmt.Errorf("get subscription failed: %v", err)
	}
	return sub, nil
}

func GenerateHelmReleaseName(subName string, feed appsv1alpha1.Feed) string {
	return fmt.Sprintf("%s-helm-%s-%s", subName, feed.Namespace, feed.Name)
}

func GetHelmRelease(clientSet client.Client, name, namespace string) (*appsv1alpha1.HelmRelease, error) {
	hr := new(appsv1alpha1.HelmRelease)
	hr.Name = name
	hr.Namespace = namespace
	key := client.ObjectKeyFromObject(hr)
	err := clientSet.Get(context.TODO(), key, hr)
	if err != nil {
		return nil, err
	}
	return hr, nil
}
