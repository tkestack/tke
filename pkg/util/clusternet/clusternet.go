package clusternet

import (
	"context"
	"errors"
	"fmt"

	appsv1alpha1 "github.com/clusternet/apis/apps/v1alpha1"
	clustersv1beta1 "github.com/clusternet/apis/clusters/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/rest"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(clustersv1beta1.AddToScheme(scheme))
	utilruntime.Must(appsv1alpha1.AddToScheme(scheme))
}

func GetHubClient(config *rest.Config) (client.Client, error) {
	var err error

	if config == nil {
		return nil, errors.New("empty  hub restconfig file")
	}

	config.ContentConfig.ContentType = "application/json"

	if err != nil {
		return nil, fmt.Errorf("fail to get hub cluster rest config ,err is %v", err)
	}

	clusternetClient, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to build a clusternet clien, error is %v", err)
	}

	return clusternetClient, nil
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
	key, err := client.ObjectKeyFromObject(sub)
	if err != nil {
		return nil, fmt.Errorf("get subscription key failed: %v", err)
	}
	err = clientSet.Get(context.TODO(), key, sub)
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
	key, err := client.ObjectKeyFromObject(hr)
	if err != nil {
		return nil, fmt.Errorf("get hemrelease key failed: %v", err)
	}
	err = clientSet.Get(context.TODO(), key, hr)
	if err != nil {
		return nil, fmt.Errorf("get helmrelease %s in %s failed: %v", name, namespace, err)
	}
	return hr, nil
}
