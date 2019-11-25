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

package util

import (
	"fmt"
	monitoringclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	influxclient "github.com/influxdata/influxdb1-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"sync"
	"time"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	esclient "tkestack.io/tke/pkg/monitor/storage/es/client"
	"tkestack.io/tke/pkg/platform/util"
)

const (
	// ProjectDatabaseName defines database name for project metrics
	ProjectDatabaseName = "projects"
)

// ClusterNameToClient mapping cluster to kubernetes client
// clusterName => kubernetes.Interface
var ClusterNameToClient sync.Map

// ClusterNameToMonitor mapping cluster to monitoring client
// clusterName => monitoringclient.Clientset
var ClusterNameToMonitor sync.Map

// GetClusterClient get kubernetes client via cluster name
func GetClusterClient(clusterName string, platformClient platformversionedclient.PlatformV1Interface) (kubernetes.Interface, error) {
	// First check from cache
	if item, ok := ClusterNameToClient.Load(clusterName); ok {
		// Check if is available
		kubeClient := item.(kubernetes.Interface)
		_, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).List(metav1.ListOptions{})
		if err == nil {
			return kubeClient, nil
		}
		ClusterNameToClient.Delete(clusterName)
	}

	kubeClient, err := util.BuildExternalClientSetWithName(platformClient, clusterName)
	if err != nil {
		return nil, err
	}

	ClusterNameToClient.Store(clusterName, kubeClient)

	return kubeClient, nil
}

// GetMonitoringClient get monitoring client via cluster name
func GetMonitoringClient(clusterName string, platformClient platformversionedclient.PlatformV1Interface) (monitoringclient.Interface, error) {
	// First check from cache
	if item, ok := ClusterNameToMonitor.Load(clusterName); ok {
		// Check if is available
		monitoringClient := item.(monitoringclient.Interface)
		_, err := monitoringClient.MonitoringV1().Prometheuses(metav1.NamespaceSystem).List(metav1.ListOptions{})
		if err == nil {
			return monitoringClient, nil
		}
		ClusterNameToClient.Delete(clusterName)
	}

	monitoringClient, err := util.BuildExternalMonitoringClientSetWithName(platformClient, clusterName)
	if err != nil {
		return nil, err
	}

	ClusterNameToClient.Store(clusterName, monitoringClient)

	return monitoringClient, nil
}

// RenameInfluxDB replace invalid character for influxDB database
func RenameInfluxDB(name string) string {
	db := strings.ReplaceAll(name, "-", "_")
	db = strings.ReplaceAll(db, ".", "_")
	return db
}

// ParseInfluxdb parse influxdb address and pwd
func ParseInfluxdb(influxdbStr string) (influxclient.Client, error) {
	url := ""
	user := ""
	pwd := ""
	influxdbStr = strings.TrimSpace(influxdbStr)

	if influxdbStr == "" {
		return nil, fmt.Errorf("no influxdb address found")
	}

	attrs := strings.Split(influxdbStr, "&")
	for _, attr := range attrs {
		if strings.HasPrefix(attr, "http") {
			url = attr
		}
		if strings.HasPrefix(attr, "u=") {
			user = attr[2:]
		}
		if strings.HasPrefix(attr, "p=") {
			pwd = attr[2:]
		}
	}

	config := influxclient.HTTPConfig{
		Addr:     url,
		Username: user,
		Password: pwd,
		Timeout:  10 * time.Second,
	}

	return influxclient.NewHTTPClient(config)
}

// ParseES parse es address and user pwd
func ParseES(esStr string) (esclient.Client, error) {
	url := ""
	user := ""
	pwd := ""
	esStr = strings.TrimSpace(esStr)

	if esStr == "" {
		return esclient.Client{}, fmt.Errorf("no influxdb address found")
	}

	attrs := strings.Split(esStr, "&")
	for _, attr := range attrs {
		if strings.HasPrefix(attr, "http") {
			url = attr
		}
		if strings.HasPrefix(attr, "u=") {
			user = attr[2:]
		}
		if strings.HasPrefix(attr, "p=") {
			pwd = attr[2:]
		}
	}

	config := esclient.Client{
		URL:      url,
		Username: user,
		Password: pwd,
	}

	return config, nil
}
