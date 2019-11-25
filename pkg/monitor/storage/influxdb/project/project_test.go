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

package project

import (
	influxdbclient "github.com/influxdata/influxdb1-client/v2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
	businessv1 "tkestack.io/tke/api/business/v1"
	fakeclient "tkestack.io/tke/api/client/clientset/versioned/fake"
)

var batchPoints influxdbclient.BatchPoints

type fakeInfluxDBClient struct {
}

func (fakeInfluxDB *fakeInfluxDBClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	return time.Since(time.Now()), "ping", nil
}
func (fakeInfluxDB *fakeInfluxDBClient) Write(bp influxdbclient.BatchPoints) error {
	batchPoints = bp
	return nil
}
func (fakeInfluxDB *fakeInfluxDBClient) Query(q influxdbclient.Query) (*influxdbclient.Response, error) {
	return nil, nil
}
func (fakeInfluxDB *fakeInfluxDBClient) QueryAsChunk(q influxdbclient.Query) (*influxdbclient.ChunkedResponse, error) {
	return nil, nil
}
func (fakeInfluxDB *fakeInfluxDBClient) Close() error {
	return nil
}
func newInfluxDBClient() influxdbclient.Client {
	return &fakeInfluxDBClient{}
}

func initBusinessClient(p *InfluxDB) {
	businessClient := fakeclient.NewSimpleClientset()

	project1 := &businessv1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: "project1",
		},
		Spec: businessv1.ProjectSpec{
			Clusters: businessv1.ClusterHard{
				"cluster1": businessv1.HardQuantity{
					Hard: businessv1.ResourceList{
						"cpu":    resource.MustParse("1"),
						"memory": resource.MustParse("1Gi"),
					},
				},
			},
		},
		Status: businessv1.ProjectStatus{
			Clusters: businessv1.ClusterUsed{
				"cluster1": businessv1.UsedQuantity{
					Used: businessv1.ResourceList{
						"cpu":    resource.MustParse("300m"),
						"memory": resource.MustParse("450Mi"),
					},
				},
			},
		},
	}
	_, _ = businessClient.BusinessV1().Projects().Create(project1)

	namespace1 := &businessv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "project1-namespace1",
			Namespace: project1.Name,
		},
		Spec: businessv1.NamespaceSpec{
			ClusterName: "cluster1",
			Hard: businessv1.ResourceList{
				"cpu":    resource.MustParse("300m"),
				"memory": resource.MustParse("450Mi"),
			},
		},
		Status: businessv1.NamespaceStatus{
			Used: businessv1.ResourceList{
				"cpu":    resource.MustParse("200m"),
				"memory": resource.MustParse("300Mi"),
			},
		},
	}
	_, _ = businessClient.BusinessV1().Namespaces(project1.Name).Create(namespace1)

	p.businessClient = businessClient.BusinessV1()
}

func initInfluxDBClient(p *InfluxDB) {
	influxdbClient := newInfluxDBClient()
	p.clients = []influxdbclient.Client{}
	p.clients = append(p.clients, influxdbClient)
}

func TestCollectProjectMetrics(t *testing.T) {
	pro := &InfluxDB{}
	initBusinessClient(pro)
	initInfluxDBClient(pro)
	pro.Collect()

	points := batchPoints.Points()
	for _, p := range points {
		tags := p.Tags()
		rName := p.Name()
		switch rName {
		case "project_capacity_cpu":
			if len(tags) != 1 {
				t.Fatalf("project_capacity_cpu tags should just one tag: %v", tags)
			}
			if tags["project_name"] == "project1" {
				fields, err := p.Fields()
				if err != nil {
					t.Fatalf("get field value err: %v", err)
				}
				if fields["value"] != 1.0 {
					t.Fatalf("wrong project_capacity_cpu(project1) should be 1.0, but %f", fields["value"])
				}
			}
		case "project_capacity_memory":
			if len(tags) != 1 {
				t.Fatalf("project_capacity_memory tags should just one tag: %v", tags)
			}
			if tags["project_name"] == "project1" {
				fields, err := p.Fields()
				if err != nil {
					t.Fatalf("get field value err: %v", err)
				}
				if fields["value"] != 1073741824.0 {
					t.Fatalf("wrong project_capacity_memory(project1) should be 1073741824.0, but %f", fields["value"])
				}
			}
		case "project_capacity_cluster_memory":
			if len(tags) != 1 {
				t.Fatalf("project_capacity_cluster_memory tags should just one tag: %v", tags)
			}
			if tags["project_name"] == "project1" {
				fields, err := p.Fields()
				if err != nil {
					t.Fatalf("get field value err: %v", err)
				}
				if fields["value"] != 471859200.0 {
					t.Fatalf("wrong project_capacity_cluster_memory(project1) should be 471859200.0, but %f", fields["value"])
				}
			}
		case "project_allocated_memory":
			if len(tags) != 1 {
				t.Fatalf("project_allocated_memory tags should just one tag: %v", tags)
			}
			if tags["project_name"] == "project1" {
				fields, err := p.Fields()
				if err != nil {
					t.Fatalf("get field value err: %v", err)
				}
				if fields["value"] != 314572800.0 {
					t.Fatalf("wrong project_allocated_memory(project1) should be 314572800.0, but %f", fields["value"])
				}
			}
		case "project_cluster_capacity_memory":
			if len(tags) != 2 {
				t.Fatalf("project_cluster_capacity_memory tags should just two tag: %v", tags)
			}
			fields, err := p.Fields()
			if err != nil {
				t.Fatalf("get field value err: %v", err)
			}
			if fields["value"] != 471859200.0 {
				t.Fatalf("wrong project_cluster_capacity_memory(project1) should be 471859200.0, but %f", fields["value"])
			}
		case "project_cluster_allocated_memory":
			if len(tags) != 2 {
				t.Fatalf("project_cluster_allocated_memory tags should just two tag: %v", tags)
			}
			fields, err := p.Fields()
			if err != nil {
				t.Fatalf("get field value err: %v", err)
			}
			if fields["value"] != 314572800.0 {
				t.Fatalf("wrong project_cluster_allocated_memory(project1) should be 314572800.0, but %f", fields["value"])
			}
		default:
			t.Logf("Ignore resource metrics: %s", rName)
		}
	}
}
