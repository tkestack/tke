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

package options

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
	"github.com/coreos/etcd/pkg/transport"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/storage/value"
	"sync/atomic"
	"time"
	"tkestack.io/tke/pkg/util/log"
)

const (
	flagETCDServers               = "etcd-servers"
	flagETCDPrefix                = "etcd-prefix"
	flagETCDKeyFile               = "etcd-keyfile"
	flagETCDCertFile              = "etcd-certfile"
	flagETCDCAFile                = "etcd-cafile"
	flagETCDCompactionInterval    = "etcd-compaction-interval"
	flagETCDCountMetricPollPeriod = "etcd-count-metric-poll-period"
)

const (
	configETCDServers               = "etcd.servers"
	configETCDPrefix                = "etcd.prefix"
	configETCDKeyFile               = "etcd.keyfile"
	configETCDCertFile              = "etcd.certfile"
	configETCDCAFile                = "etcd.cafile"
	configETCDCompactionInterval    = "etcd.compaction_interval"
	configETCDCountMetricPollPeriod = "etcd.count_metric_poll_period"
)

// The short keepalive timeout and interval have been chosen to aggressively
// detect a failed etcd server without introducing much overhead.
const keepaliveTime = 30 * time.Second
const keepaliveTimeout = 10 * time.Second

// dialTimeout is the timeout for failing to establish a connection.
// It is set to 20 seconds as times shorter than that will cause TLS connections to fail
// on heavily loaded arm64 CPUs (issue #64649)
const dialTimeout = 20 * time.Second

// ETCDClientOptions contains the options that storage backend by etcd.
type ETCDClientOptions struct {
	// Prefix is the prefix to all keys passed to storage.Interface methods.
	Prefix string
	// ServerList is the list of storage servers to connect with.
	ServerList []string
	// TLS credentials
	KeyFile  string
	CertFile string
	CAFile   string
	// Paging indicates whether the server implementation should allow paging (if it is
	// supported). This is generally configured by feature gating, or by a specific
	// resource type not wishing to allow paging, and is not intended for end users to
	// set.
	Paging bool
	Codec  runtime.Codec
	// EncodeVersioner is the same groupVersioner used to build the
	// storage encoder. Given a list of kinds the input object might belong
	// to, the EncodeVersioner outputs the gvk the object will be
	// converted to before persisted in etcd.
	EncodeVersioner runtime.GroupVersioner
	// Transformer allows the value to be transformed prior to persisting into etcd.
	Transformer value.Transformer
	// CompactionInterval is an interval of requesting compaction from apiserver.
	// If the value is 0, no compaction will be issued.
	CompactionInterval time.Duration
	// CountMetricPollPeriod specifies how often should count metric be updated
	CountMetricPollPeriod time.Duration
}

// NewETCDClientOptions creates a Options object with default parameters.
func NewETCDClientOptions(defaultETCDPathPrefix string) *ETCDClientOptions {
	return &ETCDClientOptions{
		Prefix:             defaultETCDPathPrefix,
		CompactionInterval: 5 * time.Minute,
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *ETCDClientOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSlice(flagETCDServers, o.ServerList,
		"List of etcd servers to connect with (scheme://ip:port), comma separated.")
	_ = viper.BindPFlag(configETCDServers, fs.Lookup(flagETCDServers))

	fs.String(flagETCDPrefix, o.Prefix,
		"The prefix to prepend to all resource paths in etcd.")
	_ = viper.BindPFlag(configETCDPrefix, fs.Lookup(flagETCDPrefix))

	fs.String(flagETCDKeyFile, o.KeyFile,
		"SSL key file used to secure etcd communication.")
	_ = viper.BindPFlag(configETCDKeyFile, fs.Lookup(flagETCDKeyFile))

	fs.String(flagETCDCertFile, o.CertFile,
		"SSL certification file used to secure etcd communication.")
	_ = viper.BindPFlag(configETCDCertFile, fs.Lookup(flagETCDCertFile))

	fs.String(flagETCDCAFile, o.CAFile,
		"SSL Certificate Authority file used to secure etcd communication.")
	_ = viper.BindPFlag(configETCDCAFile, fs.Lookup(flagETCDCAFile))

	fs.Duration(flagETCDCompactionInterval, o.CompactionInterval,
		"The interval of compaction requests. If 0, the compaction request from apiserver is disabled.")
	_ = viper.BindPFlag(configETCDCompactionInterval, fs.Lookup(flagETCDCompactionInterval))

	fs.Duration(flagETCDCountMetricPollPeriod, o.CountMetricPollPeriod, ""+
		"Frequency of polling etcd for number of resources per type. 0 disables the metric collection.")
	_ = viper.BindPFlag(configETCDCountMetricPollPeriod, fs.Lookup(flagETCDCountMetricPollPeriod))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *ETCDClientOptions) ApplyFlags() []error {
	var errs []error

	o.ServerList = viper.GetStringSlice(configETCDServers)
	o.CAFile = viper.GetString(configETCDCAFile)
	o.CertFile = viper.GetString(configETCDCertFile)
	o.KeyFile = viper.GetString(configETCDKeyFile)
	o.Prefix = viper.GetString(configETCDPrefix)
	o.CompactionInterval = viper.GetDuration(configETCDCompactionInterval)
	o.CountMetricPollPeriod = viper.GetDuration(configETCDCountMetricPollPeriod)

	if len(o.ServerList) == 0 {
		errs = append(errs, fmt.Errorf("--%s must be specified", flagETCDServers))
	}

	return errs
}

// NewClient creates the etcd v3 client object and returns it.
func (o *ETCDClientOptions) NewClient() (*clientv3.Client, error) {
	tlsInfo := transport.TLSInfo{
		CertFile: o.CertFile,
		KeyFile:  o.KeyFile,
		CAFile:   o.CAFile,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}
	// NOTE: Client relies on nil tlsConfig
	// for non-secure connections, update the implicit variable
	if len(o.CertFile) == 0 && len(o.KeyFile) == 0 && len(o.CAFile) == 0 {
		tlsConfig = nil
	}
	cfg := clientv3.Config{
		DialTimeout:          dialTimeout,
		DialKeepAliveTime:    keepaliveTime,
		DialKeepAliveTimeout: keepaliveTimeout,
		DialOptions: []grpc.DialOption{
			grpc.WithUnaryInterceptor(grpcprom.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(grpcprom.StreamClientInterceptor),
		},
		Endpoints: o.ServerList,
		TLS:       tlsConfig,
	}

	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	if o.Prefix != "" {
		c.KV = namespace.NewKV(c.KV, o.Prefix)
	}
	return c, nil
}

// NewHealthCheck creates the health check callback by given backend config.
func (o *ETCDClientOptions) NewHealthCheck() (func() error, error) {
	// constructing the etcd v3 client blocks and times out if etcd is not available.
	// retry in a loop in the background until we successfully create the client, storing the client or error encountered

	clientValue := &atomic.Value{}

	clientErrMsg := &atomic.Value{}
	clientErrMsg.Store("etcd client connection not yet established")

	go func() {
		if err := wait.PollUntil(time.Second, func() (bool, error) {
			client, err := o.NewClient()
			if err != nil {
				clientErrMsg.Store(err.Error())
				return false, nil
			}
			clientValue.Store(client)
			clientErrMsg.Store("")
			return true, nil
		}, wait.NeverStop); err != nil {
			log.Error("Failed to wait poll until", log.Err(err))
		}
	}()

	return func() error {
		if errMsg := clientErrMsg.Load().(string); len(errMsg) > 0 {
			return fmt.Errorf(errMsg)
		}
		client := clientValue.Load().(*clientv3.Client)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if _, err := client.Cluster.MemberList(ctx); err != nil {
			return fmt.Errorf("error listing etcd members: %v", err)
		}
		return nil
	}, nil
}
