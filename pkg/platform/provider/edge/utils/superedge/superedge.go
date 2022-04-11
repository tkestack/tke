///*
//Copyright 2020 The SuperEdge Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//*/
package superedge

//
//import (
//	"attlee-wang/tke/pkg/controller"
//	"context"
//	"encoding/base64"
//	"fmt"
//	v1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/labels"
//	"k8s.io/apimachinery/pkg/util/version"
//	"k8s.io/client-go/tools/cache"
//	"k8s.io/klog"
//	"net/url"
//	"path/filepath"
//
//	//"github.com/superedge/superedge/pkg/edgeadm/constant/manifests"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/klog/v2"
//	"k8s.io/kubernetes/pkg/apis/apps"
//	"reflect"
//	kubeclient "tkestack.io/tke/pkg/util/apiclient"
//)
//
//func DeployEdgeAPPS(client *kubernetes.Clientset, manifestsDir, caCertFile, caKeyFile, masterPublicAddr string, certSANs []string, configPath string) error {
//	if err := EnsureEdgeSystemNamespace(client); err != nil {
//		return err
//	}
//	if err := DeployEdgePreflight(client, manifestsDir, masterPublicAddr, configPath); err != nil {
//		return err
//	}
//	// Deploy tunnel
//	if err := DeployTunnelAddon(client, manifestsDir, caCertFile, caKeyFile, masterPublicAddr, certSANs); err != nil {
//		return err
//	}
//	klog.Infof("Deploy %s success!", manifests.APP_TUNNEL_EDGE)
//
//	// Deploy edge-health
//	if err := DeployEdgeHealth(client, manifestsDir); err != nil {
//		klog.Errorf("Deploy edge health, error: %s", err)
//		return err
//	}
//	klog.Infof("Deploy edge-health success!")
//
//	// Deploy service-group
//	if err := DeployServiceGroup(client, manifestsDir); err != nil {
//		klog.Errorf("Deploy serivce group, error: %s", err)
//		return err
//	}
//	klog.Infof("Deploy service-group success!")
//
//	// Deploy edge-coredns
//	if err := DeployEdgeCorednsAddon(configPath, manifestsDir); err != nil {
//		klog.Errorf("Deploy edge-coredns error: %v", err)
//		return err
//	}
//
//	// Update Kube-* Config
//	if err := UpdateKubeConfig(client); err != nil {
//		klog.Errorf("Deploy serivce group, error: %s", err)
//		return err
//	}
//	klog.Infof("Update Kubernetes cluster config support marginal autonomy success")
//
//	//Prepare config join Node
//	if err := JoinNodePrepare(client, manifestsDir, caCertFile, caKeyFile); err != nil {
//		klog.Errorf("Prepare config join Node error: %s", err)
//		return err
//	}
//	klog.Infof("Prepare join Node configMap success")
//
//	return nil
//}
//
//func EnsureEdgeSystemNamespace(client kubernetes.Interface) error {
//	if err := kubeclient.CreateOrUpdateNamespace(context.Background(), client, &v1.Namespace{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: NamespaceEdgeSystem,
//		},
//	}); err != nil {
//		return err
//	}
//	return nil
//}
//
