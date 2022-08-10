package deletion

import (
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiauthzv1 "tkestack.io/tke/api/authz/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/authz/provider"
	"tkestack.io/tke/pkg/util/log"
)

type MultiClusterRoleBindingDeleter interface {
	Delete(ctx context.Context, mcrb *apiauthzv1.MultiClusterRoleBinding, provider provider.Provider) error
}

func New(client clientset.Interface, platformClient platformversionedclient.PlatformV1Interface) MultiClusterRoleBindingDeleter {
	return &MultiClusterRoleBindingResourcesDeleter{
		client:         client,
		platformClient: platformClient,
	}
}

type MultiClusterRoleBindingResourcesDeleter struct {
	client         clientset.Interface
	platformClient platformversionedclient.PlatformV1Interface
}

func (c *MultiClusterRoleBindingResourcesDeleter) Delete(ctx context.Context, mcrb *apiauthzv1.MultiClusterRoleBinding, provider provider.Provider) error {
	// 删除集群中对应的资源
	if err := provider.DeleteMultiClusterRoleBindingResources(ctx, c.platformClient, mcrb); err != nil {
		log.Warnf("Unable to finalize MultiClusterRoleBinding '%s/%s', err: %v", mcrb.Namespace, mcrb.Name, err)
		return err
	}
	policyFinalize := apiauthzv1.MultiClusterRoleBinding{}
	policyFinalize.ObjectMeta = mcrb.ObjectMeta
	policyFinalize.Finalizers = []string{}
	if err := c.client.AuthzV1().RESTClient().Put().Resource("multiclusterrolebindings").
		Namespace(mcrb.Namespace).
		Name(mcrb.Name).
		SubResource("finalize").
		Body(&policyFinalize).
		Do(context.Background()).
		Into(&policyFinalize); err != nil {
		log.Warnf("Unable to finalize multiclusterrolebinding '%s/%s', err: %v", mcrb.Namespace, mcrb.Name, err)
		return err
	}
	if err := c.client.AuthzV1().MultiClusterRoleBindings(mcrb.Namespace).Delete(ctx, mcrb.Name, metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			log.Warnf("Unable to delete multiclusterrolebinding '%s/%s', err: %v", mcrb.Namespace, mcrb.Name, err)
			return err
		}
	}
	return nil
}
