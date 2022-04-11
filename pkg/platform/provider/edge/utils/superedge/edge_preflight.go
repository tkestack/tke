package superedge

//
//import (
//	"context"
//	"encoding/base64"
//	"fmt"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/util/version"
//	"k8s.io/client-go/kubernetes"
//	"net/url"
//	"github.com/pkg/errors"
//	rbacv1 "k8s.io/api/rbac/v1"
//	v1 "k8s.io/api/core/v1"
//	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
//	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
//	kubeadmapiv1beta2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
//	"path/filepath"
//	"tkestack.io/tke/pkg/util/apiclient"
//	"github.com/pkg/errors"
//	v1 "k8s.io/api/core/v1"
//	rbacv1 "k8s.io/api/rbac/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/util/version"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/klog/v2"
//	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
//	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
//	kubeadmapiv1beta2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
//	"k8s.io/kubernetes/cmd/kubeadm/app/componentconfigs"
//	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
//	clusterinfophase "k8s.io/kubernetes/cmd/kubeadm/app/phases/bootstraptoken/clusterinfo"
//	nodebootstraptokenphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/bootstraptoken/node"
//	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
//	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
//	configutil "k8s.io/kubernetes/cmd/kubeadm/app/util/config"
//	clusterinfophase "k8s.io/kubernetes/cmd/kubeadm/app/phases/bootstraptoken/clusterinfo"
//	nodebootstraptokenphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/bootstraptoken/node"
//	//kubeclient "tkestack.io/tke/pkg/util/apiclient"
//	"github.com/superedge/superedge/pkg/util/kubeclient"
//)
//
//// DeployEdgePreflight is a preflight step for the addon.
//func DeployEdgePreflight(clientSet kubernetes.Interface, manifestsDir string, masterPublicAddr string, configPath string) error {
//	if err := ensureKubePublicNamespace(clientSet); err != nil {
//		return err
//	}
//
//	if err := ensureBootstrapTokenRBAC(clientSet, masterPublicAddr, manifestsDir, configPath); err != nil {
//		return err
//	}
//	clusterConfiguration, err := ensureKubeadmConfigConfigMap(clientSet, configPath)
//	if err != nil {
//		return err
//	}
//
//	if err := ensureKubeletConfigMap(clientSet, clusterConfiguration); err != nil {
//		return err
//	}
//
//	if err := ensureKubeadmRBAC(clientSet); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// ensureKubePublicNamespace ensure the namespace kube-public exits.
//func ensureKubePublicNamespace(client kubernetes.Interface) error {
//	if err := kubeclient.CreateOrUpdateNamespace(context.Background(), client, &v1.Namespace{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: NamespaceKubePublic,
//		},
//	}); err != nil {
//		return err
//	}
//	return nil
//}
//
//// ensureBootstrapTokenRBAC ensure the bootstrap token RBAC rules exist.
//func ensureBootstrapTokenRBAC(clientSet kubernetes.Interface, masterPublicAddr string, manifestsDir string, configPath string) error {
//
//	// Create RBAC rules that makes the bootstrap tokens able to get nodes
//	if err := nodebootstraptokenphase.AllowBoostrapTokensToGetNodes(clientSet); err != nil {
//		return errors.Wrap(err, "error allowing bootstrap tokens to get Nodes")
//	}
//	// Create RBAC rules that makes the bootstrap tokens able to post CSRs
//	if err := nodebootstraptokenphase.AllowBootstrapTokensToPostCSRs(clientSet); err != nil {
//		return errors.Wrap(err, "error allowing bootstrap tokens to post CSRs")
//	}
//	// Create RBAC rules that makes the bootstrap tokens able to get their CSRs approved automatically
//	if err := nodebootstraptokenphase.AutoApproveNodeBootstrapTokens(clientSet); err != nil {
//		return errors.Wrap(err, "error auto-approving node bootstrap tokens")
//	}
//
//	// Create/update RBAC rules that makes the nodes to rotate certificates and get their CSRs approved automatically
//	if err := nodebootstraptokenphase.AutoApproveNodeCertificateRotation(clientSet); err != nil {
//		return err
//	}
//
//	// Create the cluster-info ConfigMap with the associated RBAC rules
//	if err := createBootstrapConfigMapIfNotExists(clientSet, masterPublicAddr, manifestsDir, configPath); err != nil {
//		return errors.Wrap(err, "error creating bootstrap ConfigMap")
//	}
//
//	if err := clusterinfophase.CreateClusterInfoRBACRules(clientSet); err != nil {
//		return errors.Wrap(err, "error creating clusterinfo RBAC rules")
//	}
//	return nil
//}
//
//// createBootstrapConfigMapIfNotExists ensure the bootstrap configmap exists.
//func createBootstrapConfigMapIfNotExists(clientSet kubernetes.Interface, masterPublicAddr string, manifestsDir string, configPath string) error {
//	clusterInfoKubeConfig := filepath.Join(manifestsDir, manifests.ClusterInfoKubeConfig)
//	cluster, err := kubeclient.GetClusterInfo(configPath)
//	if err != nil {
//		return err
//	}
//	server := cluster.Server
//	if masterPublicAddr != "" {
//		server = fmt.Sprintf("https://%s", masterPublicAddr)
//	}
//	yamlKubeConfig, err := kubeclient.ParseString(
//		ReadYaml(clusterInfoKubeConfig, manifests.ClusterInfoKubeConfigYaml),
//		map[string]interface{}{
//			"Server": server,
//			"CAData": base64.StdEncoding.EncodeToString(cluster.CertificateAuthorityData),
//		})
//	if err != nil {
//		return err
//	}
//
//	configMap := &v1.ConfigMap{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "cluster-info",
//			Namespace: constant.NamespaceKubePublic,
//		},
//		Data: map[string]string{
//			"kubeconfig": string(yamlKubeConfig),
//		},
//	}
//
//	return apiclient.CreateOrRetainConfigMap(clientSet, configMap, "cluster-info")
//}
//
//// ensureKubeadmRBAC ensure the RBAC rules for the kubeadm.
//func ensureKubeadmRBAC(clientSet kubernetes.Interface) error {
//	// kubeadm:nodes-kubeadm-config
//	role := rbacv1.Role{
//		TypeMeta: metav1.TypeMeta{},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      NodesKubeadmConfigClusterRoleName,
//			Namespace: constant.NamespaceKubeSystem,
//		},
//		Rules: nil,
//	}
//
//	role.Rules = append(role.Rules, rbacv1.PolicyRule{
//		APIGroups:     []string{"*"},
//		Resources:     []string{"configmaps"},
//		ResourceNames: []string{kubeadmconstants.KubeadmConfigConfigMap},
//		Verbs:         []string{"get", "list", "watch"},
//	})
//
//	roleBinding := rbacv1.RoleBinding{
//		TypeMeta: metav1.TypeMeta{},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      NodesKubeadmConfigClusterRoleName,
//			Namespace: constant.NamespaceKubeSystem,
//		},
//		RoleRef: rbacv1.RoleRef{
//			Name:     NodesKubeadmConfigClusterRoleName,
//			Kind:     "Role",
//			APIGroup: rbacv1.GroupName,
//		},
//		Subjects: nil,
//	}
//
//	roleBinding.Subjects = append(roleBinding.Subjects, rbacv1.Subject{
//		APIGroup: rbacv1.GroupName,
//		Kind:     rbacv1.GroupKind,
//		Name:     kubeadmconstants.NodeBootstrapTokenAuthGroup,
//	}, rbacv1.Subject{
//		APIGroup: rbacv1.GroupName,
//		Kind:     rbacv1.GroupKind,
//		Name:     kubeadmconstants.NodesGroup,
//	})
//
//	//kubeadm:get-nodes
//	getNodeClusterRole := rbacv1.ClusterRole{
//		TypeMeta: metav1.TypeMeta{},
//		ObjectMeta: metav1.ObjectMeta{
//			Name: GetNodesClusterRoleName,
//		},
//		Rules: nil,
//	}
//
//	getNodeClusterRole.Rules = append(getNodeClusterRole.Rules, rbacv1.PolicyRule{
//		APIGroups: []string{""},
//		Resources: []string{"nodes"},
//		Verbs:     []string{"get"},
//	})
//
//	getNodeClusterRoleBinding := rbacv1.ClusterRoleBinding{
//		TypeMeta: metav1.TypeMeta{},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      GetNodesClusterRoleName,
//			Namespace: NamespaceKubeSystem,
//		},
//		RoleRef: rbacv1.RoleRef{
//			Name:     GetNodesClusterRoleName,
//			Kind:     "ClusterRole",
//			APIGroup: rbacv1.GroupName,
//		},
//		Subjects: nil,
//	}
//
//	getNodeClusterRoleBinding.Subjects = append(getNodeClusterRoleBinding.Subjects, rbacv1.Subject{
//		APIGroup: rbacv1.GroupName,
//		Kind:     rbacv1.GroupKind,
//		Name:     kubeadmconstants.NodeBootstrapTokenAuthGroup,
//	})
//
//	//kubeadm:kubelet-bootstrap
//	kubeletBootstrapClusterRole := rbacv1.ClusterRole{
//		TypeMeta: metav1.TypeMeta{},
//		ObjectMeta: metav1.ObjectMeta{
//			Name: NodeBootstrapperClusterRoleName,
//		},
//		Rules: nil,
//	}
//
//	kubeletBootstrapClusterRole.Rules = append(kubeletBootstrapClusterRole.Rules, rbacv1.PolicyRule{
//		APIGroups: []string{"certificates.k8s.io"},
//		Resources: []string{"certificatesigningrequests"},
//		Verbs:     []string{"get", "list", "watch", "create"},
//	})
//
//	kubeletBootstrapClusterRoleBinding := rbacv1.ClusterRoleBinding{
//		TypeMeta: metav1.TypeMeta{},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      KubeletBootstrapClusterRoleName,
//			Namespace: constant.NamespaceKubeSystem,
//		},
//		RoleRef: rbacv1.RoleRef{
//			Name:     NodeBootstrapperClusterRoleName,
//			Kind:     "ClusterRole",
//			APIGroup: rbacv1.GroupName,
//		},
//		Subjects: nil,
//	}
//
//	kubeletBootstrapClusterRoleBinding.Subjects = append(kubeletBootstrapClusterRoleBinding.Subjects, rbacv1.Subject{
//		APIGroup: rbacv1.GroupName,
//		Kind:     rbacv1.GroupKind,
//		Name:     kubeadmconstants.NodeBootstrapTokenAuthGroup,
//	})
//
//	if err := apiclient.CreateOrUpdateRole(context.Background(), clientSet, &role); err != nil {
//		return err
//	}
//
//	if err := apiclient.CreateOrUpdateRoleBinding(context.Background(), clientSet, &roleBinding); err != nil {
//		return err
//	}
//
//	if err := apiclient.CreateOrUpdateClusterRole(context.Background(), clientSet, &getNodeClusterRole); err != nil {
//		return err
//	}
//
//	if err := apiclient.CreateOrUpdateClusterRoleBinding(context.Background(), clientSet, &getNodeClusterRoleBinding); err != nil {
//		return err
//	}
//
//	if err := apiclient.CreateOrUpdateClusterRole(context.Background(), clientSet, &kubeletBootstrapClusterRole); err != nil {
//		return err
//	}
//
//	return apiclient.CreateOrUpdateClusterRoleBinding(context.Background(), clientSet, &kubeletBootstrapClusterRoleBinding)
//}
//
//// ensureKubeadmConfigConfigMap ensure a ConfigMap with the generic kubeadm configuration.
//func ensureKubeadmConfigConfigMap(clientSet kubernetes.Interface, configPath string) (*kubeadmapi.ClusterConfiguration, error) {
//	hostname, err := kubeadmutil.GetHostname("")
//	cluster, err := kubeclient.GetClusterInfo(configPath)
//	if err != nil {
//		return nil, err
//	}
//
//	clusterConfigurationToUpload := &kubeadmapi.ClusterConfiguration{
//		APIServer: kubeadmapi.APIServer{
//			TimeoutForControlPlane: &metav1.Duration{
//				Duration: kubeadmconstants.DefaultControlPlaneTimeout,
//			},
//		},
//		DNS: kubeadmapi.DNS{
//			Type: kubeadmapi.CoreDNS,
//		},
//		CertificatesDir: v1beta2.DefaultCertificatesDir,
//		ClusterName:     v1beta2.DefaultClusterName,
//		Etcd: kubeadmapi.Etcd{
//			Local: &kubeadmapi.LocalEtcd{
//				DataDir: v1beta2.DefaultEtcdDataDir,
//			},
//		},
//		ImageRepository:   constant.ImageRepository,
//		KubernetesVersion: kubeadmconstants.CurrentKubernetesVersion.String(),
//		Networking: kubeadmapi.Networking{
//			ServiceSubnet: v1beta2.DefaultServicesSubnet,
//			DNSDomain:     v1beta2.DefaultServiceDNSDomain,
//		},
//	}
//	clusterConfigurationToUpload.ComponentConfigs = kubeadmapi.ComponentConfigMap{}
//
//	// Marshal the ClusterConfiguration into YAML
//	clusterConfigurationYaml, err := configutil.MarshalKubeadmConfigObject(clusterConfigurationToUpload)
//	if err != nil {
//		return nil, err
//	}
//
//	// Prepare the ClusterStatus for upload
//	clusterUrl, err := url.Parse(cluster.Server)
//	if err != nil {
//		return nil, errors.Wrap(err, "error parsing cluster server")
//	}
//	apiEndpoint, err := kubeadmapi.APIEndpointFromString(clusterUrl.Host)
//	if err != nil {
//		return nil, err
//	}
//	clusterStatus := &kubeadmapi.ClusterStatus{
//		APIEndpoints: map[string]kubeadmapi.APIEndpoint{
//			hostname: apiEndpoint,
//		},
//	}
//	// Marshal the ClusterStatus into YAML
//	clusterStatusYaml, err := configutil.MarshalKubeadmConfigObject(clusterStatus)
//	if err != nil {
//		return nil, err
//	}
//
//	configMap := &v1.ConfigMap{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      kubeadmconstants.KubeadmConfigConfigMap,
//			Namespace: constant.NamespaceKubeSystem,
//		},
//		Data: map[string]string{
//			kubeadmconstants.ClusterConfigurationConfigMapKey: string(clusterConfigurationYaml),
//			kubeadmconstants.ClusterStatusConfigMapKey:        string(clusterStatusYaml),
//		},
//	}
//
//	apiclient.CreateOrRetainConfigMap(clientSet, configMap, kubeadmconstants.KubeadmConfigConfigMap)
//	return clusterConfigurationToUpload, nil
//}
//
//// ensureKubeletConfigMap ensure a ConfigMap with the generic kubelet configuration.
//func ensureKubeletConfigMap(clientSet kubernetes.Interface, clusterConfiguration *kubeadmapi.ClusterConfiguration) error {
//	clusterCfg := &kubeadmapiv1beta2.ClusterConfiguration{
//		KubernetesVersion: kubeadmconstants.CurrentKubernetesVersion.String(),
//	}
//	internalcfg, err := configutil.DefaultedInitConfiguration(&kubeadmapiv1beta2.InitConfiguration{}, clusterCfg)
//	if err != nil {
//		errors.Wrapf(err, "unexpected failure by DefaultedInitConfiguration: %v", clusterCfg)
//	}
//	klog.V(1).Infoln("[upload-config] Uploading the kubelet component config to a ConfigMap")
//	if err := createConfigMap(&internalcfg.ClusterConfiguration, clientSet); err != nil {
//		return errors.Wrap(err, "error creating kubelet configuration ConfigMap")
//	}
//	return nil
//}
//
//// createConfigMap creates a ConfigMap with the generic kubelet configuration.
//func createConfigMap(cfg *kubeadmapi.ClusterConfiguration, client kubernetes.Interface) error {
//
//	k8sVersion, err := version.ParseSemantic(cfg.KubernetesVersion)
//	if err != nil {
//		return err
//	}
//
//	configMapName := kubeadmconstants.GetKubeletConfigMapName(k8sVersion)
//	klog.V(1).Infof("[kubelet] Creating a ConfigMap %q in namespace %s with the configuration for the kubelets in the cluster\n", configMapName, constant.NamespaceKubeSystem)
//
//	kubeletCfg, ok := cfg.ComponentConfigs[componentconfigs.KubeletGroup]
//	if !ok {
//		return errors.New("no kubelet component config found in the active component config set")
//	}
//
//	kubeletBytes, err := kubeletCfg.Marshal()
//	if err != nil {
//		return err
//	}
//
//	configMap := &v1.ConfigMap{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      configMapName,
//			Namespace: constant.NamespaceKubeSystem,
//		},
//		Data: map[string]string{
//			kubeadmconstants.KubeletBaseConfigurationConfigMapKey: string(kubeletBytes),
//		},
//	}
//
//	if !kubeletCfg.IsUserSupplied() {
//		componentconfigs.SignConfigMap(configMap)
//	}
//
//	if err := apiclient.CreateOrRetainConfigMap(client, configMap, configMapName); err != nil {
//		return err
//	}
//
//	if err := createConfigMapRBACRules(client, k8sVersion); err != nil {
//		return errors.Wrap(err, "error creating kubelet configuration configmap RBAC rules")
//	}
//	return nil
//}
//
//// createConfigMapRBACRules creates the RBAC rules for exposing the base kubelet ConfigMap in the kube-system namespace to unauthenticated users
//func createConfigMapRBACRules(client kubernetes.Interface, k8sVersion *version.Version) error {
//	if err := apiclient.CreateOrUpdateRole(client, &rbacv1.Role{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      configMapRBACName(k8sVersion),
//			Namespace: metav1.NamespaceSystem,
//		},
//		Rules: []rbacv1.PolicyRule{
//			{
//				Verbs:         []string{"get"},
//				APIGroups:     []string{""},
//				Resources:     []string{"configmaps"},
//				ResourceNames: []string{kubeadmconstants.GetKubeletConfigMapName(k8sVersion)},
//			},
//		},
//	}); err != nil {
//		return err
//	}
//
//	return apiclient.CreateOrUpdateRoleBinding(client, &rbacv1.RoleBinding{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      configMapRBACName(k8sVersion),
//			Namespace: metav1.NamespaceSystem,
//		},
//		RoleRef: rbacv1.RoleRef{
//			APIGroup: rbacv1.GroupName,
//			Kind:     "Role",
//			Name:     configMapRBACName(k8sVersion),
//		},
//		Subjects: []rbacv1.Subject{
//			{
//				Kind: rbacv1.GroupKind,
//				Name: kubeadmconstants.NodesGroup,
//			},
//			{
//				Kind: rbacv1.GroupKind,
//				Name: kubeadmconstants.NodeBootstrapTokenAuthGroup,
//			},
//		},
//	})
//}
//
//// configMapRBACName returns the name for the Role/RoleBinding for the kubelet config configmap for the right branch of k8s
//func configMapRBACName(k8sVersion *version.Version) string {
//	return fmt.Sprintf("%s%d.%d", kubeadmconstants.KubeletBaseConfigMapRolePrefix, k8sVersion.Major(), k8sVersion.Minor())
//}
//
//
//
