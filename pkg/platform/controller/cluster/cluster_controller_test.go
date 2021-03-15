/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package cluster

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/prashantv/gostub"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/client-go/rest"
	restfake "k8s.io/client-go/rest/fake"
	"tkestack.io/tke/api/client/clientset/versioned/fake"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/controller/cluster/deletion"
	"tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/scheme"
	"tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/types"
	platformtypev1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/log"
)

type tkeProvider struct {
}

func (c *tkeProvider) Name() string {
	return "Tke"
}

func (c *tkeProvider) RegisterHandler(mux *mux.PathRecorderMux) {
}

func (c *tkeProvider) Validate(cluster *types.Cluster) field.ErrorList {
	return field.ErrorList{}
}

func (c *tkeProvider) PreCreate(cluster *types.Cluster) error {
	return nil
}

func (c *tkeProvider) AfterCreate(cluster *types.Cluster) error {
	return nil
}

// Setup called by controller to give an chance for plugin do some init work.
func (c *tkeProvider) Setup() error {
	return nil
}

// Teardown called by controller for plugin do some clean job.
func (c *tkeProvider) Teardown() error {
	return nil
}

func (c *tkeProvider) OnCreate(ctx context.Context, cluster *platformtypev1.Cluster) error {
	log.Info("======= OnCreate =========")
	cluster.Status.Phase = platformv1.ClusterRunning
	return nil
}

func (c *tkeProvider) OnUpdate(ctx context.Context, cluster *platformtypev1.Cluster) error {
	log.Info("======= OnUpdate =========")
	return nil
}

func (c *tkeProvider) OnDelete(ctx context.Context, cluster *platformtypev1.Cluster) error {
	log.Info("======= OnDelete =========")
	return nil
}

// OnRunning call on first running.
func (c *tkeProvider) OnRunning(ctx context.Context, cluster *platformtypev1.Cluster) error {
	log.Info("=======  OnRunning =========")
	return nil
}

func TestController_Run(t *testing.T) {
	// remove unused error handler
	utilruntime.ErrorHandlers = []func(error){}

	// register test provider
	cluster.Register("Tke", &tkeProvider{})

	type args struct {
		workers int
		stopCh  chan struct{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete cluster when initializing",
			args: args{
				workers: 5,
				stopCh:  make(chan struct{}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// stub deletion get RESTClient, in this case return success
			stubs := gostub.Stub(&deletion.StubGetPlatformRESTClient, func(client v1clientset.PlatformV1Interface) rest.Interface {
				handler := func(req *http.Request) (resp *http.Response, err error) {
					data := `{"metadata":{"name":"cls-jlzv6qqv","generateName":"cls-","selfLink":"/apis/platform.tkestack.io/v1/clusters/cls-jlzv6qqv","uid":"ce199dba-dd0a-446e-97eb-0e27c80e5709","resourceVersion":"6537333","creationTimestamp":"2021-02-09T07:43:01Z"},"spec":{"finalizers":[],"tenantID":"100004603305","displayName":"test","type":"Tke","version":"1.16.3","networkDevice":"eth0","clusterCIDR":"172.18.0.0/16","features":{"upgrade":{"strategy":{}}},"properties":{},"clusterCredentialRef":{"name":"cc-t4dskl8l"}},"status":{"version":"","phase":"Initializing","addresses":[{"type":"Internal","host":"10.29.45.5","port":10413,"path":""},{"type":"Advertise","host":"169.254.128.57","port":6443,"path":""},{"type":"Real","host":"169.254.128.57","port":6443,"path":""}],"resource":{}}}`
					resp = &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(data))),
					}
					return
				}

				return &restfake.RESTClient{
					NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
					Client:               restfake.CreateHTTPClient(handler),
				}
			})
			defer stubs.Reset()

			// mock apiserver
			f := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/xml")
				fmt.Fprintln(w, "{}")
			}
			server := httptest.NewTLSServer(http.HandlerFunc(f))
			defer server.Close()
			u, _ := url.Parse(server.URL)
			uPort, _ := strconv.ParseInt(u.Port(), 10, 32)

			// start controller
			clientset := fake.NewSimpleClientset()
			sharedInformers := versionedinformers.NewSharedInformerFactory(clientset, 200*time.Second)
			c := NewController(
				clientset.PlatformV1(),
				sharedInformers.Platform().V1().Clusters(),
				60*time.Second,
				platformv1.ClusterFinalize,
			)
			sharedInformers.Start(tt.args.stopCh)

			go func() {
				credential := &platformv1.ClusterCredential{
					ObjectMeta: v1.ObjectMeta{
						Name: "cc-t4dskl8l",
					},
					TenantID:    "100004603305",
					ClusterName: "cls-jlzv6qqv",
				}
				cluster := &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{
						Name:              "cls-jlzv6qqv",
						GenerateName:      "cls-",
						SelfLink:          "/apis/platform.tkestack.io/v1/clusters/cls-jlzv6qqv",
						UID:               "ce199dba-dd0a-446e-97eb-0e27c80e5709",
						ResourceVersion:   "6537333",
						CreationTimestamp: v1.NewTime(time.Now()),
					},
					Spec: platformv1.ClusterSpec{
						Finalizers:    []platformv1.FinalizerName{platformv1.ClusterFinalize},
						TenantID:      "100004603305",
						DisplayName:   "test",
						Type:          "Tke",
						Version:       "1.16.3",
						NetworkDevice: "eth0",
						ClusterCIDR:   "172.18.0.0/16",
						ClusterCredentialRef: &corev1.LocalObjectReference{
							Name: "cc-t4dskl8l",
						},
					},
					Status: platformv1.ClusterStatus{
						Phase: platformv1.ClusterInitializing,
						Addresses: []platformv1.ClusterAddress{
							{
								Type: "Internal",
								Host: u.Hostname(),
								Port: int32(uPort),
								Path: "",
							},
							{
								Type: "Advertise",
								Host: u.Hostname(),
								Port: int32(uPort),
								Path: "",
							},
							{
								Type: "Real",
								Host: u.Hostname(),
								Port: int32(uPort),
								Path: "",
							},
						},
					},
				}

				time.Sleep(100 * time.Millisecond)
				log.Info("mock create cluster......")
				clientset.PlatformV1().ClusterCredentials().Create(context.TODO(), credential, v1.CreateOptions{})
				clientset.PlatformV1().Clusters().Create(context.TODO(), cluster, v1.CreateOptions{})

				time.Sleep(500 * time.Millisecond)
				log.Info("mock delete cluster......")
				now := v1.NewTime(time.Now())
				cluster.ObjectMeta.DeletionTimestamp = &now
				cluster.Status.Phase = platformv1.ClusterTerminating
				clientset.PlatformV1().Clusters().Update(context.TODO(), cluster, v1.UpdateOptions{})

				time.Sleep(100 * time.Millisecond)
				close(tt.args.stopCh)
			}()

			if err := c.Run(tt.args.workers, tt.args.stopCh); (err != nil) != tt.wantErr {
				t.Errorf("Controller.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestController_Run_Supervene(t *testing.T) {
	// remove unused error handler
	utilruntime.ErrorHandlers = []func(error){}

	// register test provider
	cluster.Register("Tke", &tkeProvider{})

	type args struct {
		workers int
		stopCh  chan struct{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete cluster when initializing",
			args: args{
				workers: 5,
				stopCh:  make(chan struct{}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// stub deletion get RESTClient, in this case return success
			stubs := gostub.Stub(&deletion.StubGetPlatformRESTClient, func(client v1clientset.PlatformV1Interface) rest.Interface {
				handler := func(req *http.Request) (resp *http.Response, err error) {
					data := `{"metadata":{"name":"cls-jlzv6qqv","generateName":"cls-","selfLink":"/apis/platform.tkestack.io/v1/clusters/cls-jlzv6qqv","uid":"ce199dba-dd0a-446e-97eb-0e27c80e5709","resourceVersion":"6537333","creationTimestamp":"2021-02-09T07:43:01Z"},"spec":{"finalizers":[],"tenantID":"100004603305","displayName":"test","type":"Tke","version":"1.16.3","networkDevice":"eth0","clusterCIDR":"172.18.0.0/16","features":{"upgrade":{"strategy":{}}},"properties":{},"clusterCredentialRef":{"name":"cc-t4dskl8l"}},"status":{"version":"","phase":"Initializing","addresses":[{"type":"Internal","host":"10.29.45.5","port":10413,"path":""},{"type":"Advertise","host":"169.254.128.57","port":6443,"path":""},{"type":"Real","host":"169.254.128.57","port":6443,"path":""}],"resource":{}}}`
					resp = &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(data))),
					}
					return
				}

				return &restfake.RESTClient{
					NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
					Client:               restfake.CreateHTTPClient(handler),
				}
			})
			defer stubs.Reset()

			f := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/xml")
				fmt.Fprintln(w, "{}")
			}
			server := httptest.NewTLSServer(http.HandlerFunc(f))
			defer server.Close()
			u, _ := url.Parse(server.URL)
			uPort, _ := strconv.ParseInt(u.Port(), 10, 32)

			clientset := fake.NewSimpleClientset()
			sharedInformers := versionedinformers.NewSharedInformerFactory(clientset, 200*time.Second)

			c := NewController(
				clientset.PlatformV1(),
				sharedInformers.Platform().V1().Clusters(),
				60*time.Second,
				platformv1.ClusterFinalize,
			)

			sharedInformers.Start(tt.args.stopCh)

			go func() {
				credential := &platformv1.ClusterCredential{
					ObjectMeta: v1.ObjectMeta{
						Name: "cc-t4dskl8l",
					},
					TenantID:    "100004603305",
					ClusterName: "cls-jlzv6qqv",
				}
				cluster := &platformv1.Cluster{
					ObjectMeta: v1.ObjectMeta{
						Name:              "cls-jlzv6qqv",
						GenerateName:      "cls-",
						SelfLink:          "/apis/platform.tkestack.io/v1/clusters/cls-jlzv6qqv",
						UID:               "ce199dba-dd0a-446e-97eb-0e27c80e5709",
						ResourceVersion:   "6537333",
						CreationTimestamp: v1.NewTime(time.Now()),
					},
					Spec: platformv1.ClusterSpec{
						Finalizers:    []platformv1.FinalizerName{platformv1.ClusterFinalize},
						TenantID:      "100004603305",
						DisplayName:   "test",
						Type:          "Tke",
						Version:       "1.16.3",
						NetworkDevice: "eth0",
						ClusterCIDR:   "172.18.0.0/16",
						ClusterCredentialRef: &corev1.LocalObjectReference{
							Name: "cc-t4dskl8l",
						},
					},
					Status: platformv1.ClusterStatus{
						Phase: platformv1.ClusterInitializing,
						Addresses: []platformv1.ClusterAddress{
							{
								Type: "Internal",
								Host: u.Hostname(),
								Port: int32(uPort),
								Path: "",
							},
							{
								Type: "Advertise",
								Host: "169.254.128.57",
								Port: 6443,
								Path: "",
							},
							{
								Type: "Real",
								Host: "169.254.128.57",
								Port: 6443,
								Path: "",
							},
						},
					},
				}

				time.Sleep(100 * time.Millisecond)
				log.Info("mock create cluster......")
				clientset.PlatformV1().ClusterCredentials().Create(context.TODO(), credential, v1.CreateOptions{})
				clientset.PlatformV1().Clusters().Create(context.TODO(), cluster, v1.CreateOptions{})
				time.Sleep(500 * time.Millisecond)

				for i := 0; i < 5; i++ {
					cluster, _ = clientset.PlatformV1().Clusters().Get(context.TODO(), "cls-jlzv6qqv", v1.GetOptions{})
					cluster.Spec.Version = fmt.Sprintf("%d", i)
					clientset.PlatformV1().Clusters().Update(context.TODO(), cluster, v1.UpdateOptions{})
					time.Sleep(100 * time.Millisecond)
				}

				time.Sleep(100 * time.Millisecond)
				close(tt.args.stopCh)
			}()

			if err := c.Run(tt.args.workers, tt.args.stopCh); (err != nil) != tt.wantErr {
				t.Errorf("Controller.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
