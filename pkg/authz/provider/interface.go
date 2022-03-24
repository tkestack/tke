package provider

import (
	"context"
	"fmt"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	authzv1 "tkestack.io/tke/api/authz/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	apiplatformv1 "tkestack.io/tke/api/platform/v1"
	platformv1 "tkestack.io/tke/pkg/platform/types/v1"
)

type Provider interface {
	Name() string
	OnFilter(ctx context.Context, annotations map[string]string) bool
	InitContext(param interface{}) context.Context
	GetTenantClusters(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, tenantID string) ([]string, error)
	GetSubject(ctx context.Context, userName string, cluster *platformv1.Cluster) (*rbacv1.Subject, error)
	DispatchMultiClusterRoleBinding(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, mcrb *authzv1.MultiClusterRoleBinding, rules []rbacv1.PolicyRule, clusterSubjects map[string]*rbacv1.Subject) error
	DeleteUnbindingResources(ctx context.Context, client platformversionedclient.PlatformV1Interface, mcrb *authzv1.MultiClusterRoleBinding, clusterIDs []string) error
	DeleteMultiClusterRoleBindingResources(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, mcrb *authzv1.MultiClusterRoleBinding) error
}

var _ Provider = &DelegateProvider{}

type DelegateProvider struct {
	ProviderName string
}

func (p *DelegateProvider) OnFilter(todo context.Context, annotations map[string]string) bool {
	return true
}

func (p *DelegateProvider) Name() string {
	if p.ProviderName == "" {
		return "unknown"
	}
	return p.ProviderName
}

func (p *DelegateProvider) InitContext(param interface{}) context.Context {
	return context.Background()
}

func (p *DelegateProvider) GetTenantClusters(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, tenantID string) ([]string, error) {
	var clusterIDs []string

	listOptions := metav1.ListOptions{
		ResourceVersion: "0",
		FieldSelector:   fmt.Sprintf("spec.tenantID=%s", tenantID),
	}
	clusters, err := platformClient.Clusters().List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	for _, cls := range clusters.Items {
		if cls.Spec.TenantID == tenantID && cls.Name != "global" {
			if cls.Status.Phase != apiplatformv1.ClusterInitializing && cls.Status.Phase != apiplatformv1.ClusterTerminating {
				clusterIDs = append(clusterIDs, cls.Name)
			}
		}
	}
	return clusterIDs, nil
}

func (p *DelegateProvider) GetSubject(ctx context.Context, platformUser string, cluster *platformv1.Cluster) (*rbacv1.Subject, error) {
	_, err := cluster.RESTConfig()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (p *DelegateProvider) DispatchMultiClusterRoleBinding(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, mcrb *authzv1.MultiClusterRoleBinding, rules []rbacv1.PolicyRule, clusterSubjects map[string]*rbacv1.Subject) error {
	return nil
}

func (p *DelegateProvider) DeleteUnbindingResources(ctx context.Context, client platformversionedclient.PlatformV1Interface, mcrb *authzv1.MultiClusterRoleBinding, clusterIDs []string) error {
	return nil
}

func (p *DelegateProvider) DeleteMultiClusterRoleBindingResources(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, mcrb *authzv1.MultiClusterRoleBinding) error {
	return nil
}
