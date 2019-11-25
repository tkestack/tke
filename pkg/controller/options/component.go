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
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	componentconfig "k8s.io/component-base/config"
)

const (
	flagComponentMinResyncPeriod         = "min-resync-period"
	flagComponentControllerStartInterval = "controller-start-interval"
	flagComponentControllers             = "controllers"
	flagComponentLeaderElect             = "leader-elect"
	flagComponentLeaseDuration           = "leader-elect-lease-duration"
	flagComponentRenewDeadline           = "leader-elect-renew-deadline"
	flagComponentRetryPeriod             = "leader-elect-retry-period"
)

const (
	configComponentMinResyncPeriod         = "component.min_resync_period"
	configComponentControllerStartInterval = "component.controller_start_interval"
	configComponentControllers             = "component.controllers"
	configComponentLeaderElect             = "component.leader_elect"
	configComponentLeaseDuration           = "component.leader_elect_lease_duration"
	configComponentRenewDeadline           = "component.leader_elect_renew_deadline"
	configComponentRetryPeriod             = "component.leader_elect_retry_period"
)

// ComponentOptions contains configuration items related to controller manager
// component attributes.
type ComponentOptions struct {
	MinResyncPeriod time.Duration

	ControllerStartInterval time.Duration
	LeaderElection          *LeaderElectionConfiguration
	Controllers             []string
	ContainerRegistryDomain string

	allControllers               []string
	disabledByDefaultControllers []string
}

// ComponentConfiguration holds configuration for a generic controller-manager
type ComponentConfiguration struct {
	// minResyncPeriod is the resync period in reflectors; will be random between
	// minResyncPeriod and 2*minResyncPeriod.
	MinResyncPeriod time.Duration
	// ClientConnection specifies the kubeconfig file and client connection
	// settings for the proxy server to use when communicating with the apiserver.
	ClientConnection componentconfig.ClientConnectionConfiguration
	// How long to wait between starting controller managers
	ControllerStartInterval time.Duration
	// leaderElection defines the configuration of leader election client.
	LeaderElection componentconfig.LeaderElectionConfiguration
	// Controllers is the list of controllers to enable or disable
	// '*' means "all enabled by default controllers"
	// 'foo' means "enable 'foo'"
	// '-foo' means "disable 'foo'"
	// first item for a particular name wins
	Controllers []string
	// DebuggingConfiguration holds configuration for Debugging related features.
	Debugging componentconfig.DebuggingConfiguration
}

// LeaderElectionConfiguration defines the configuration of leader election
// clients for components that can run with leader election enabled.
type LeaderElectionConfiguration struct {
	// leaderElect enables a leader election client to gain leadership
	// before executing the main loop. Enable this when running replicated
	// components for high availability.
	LeaderElect bool
	// leaseDuration is the duration that non-leader candidates will wait
	// after observing a leadership renewal until attempting to acquire
	// leadership of a led but unrenewed leader slot. This is effectively the
	// maximum duration that a leader can be stopped before it is replaced
	// by another candidate. This is only applicable if leader election is
	// enabled.
	LeaseDuration time.Duration
	// renewDeadline is the interval between attempts by the acting master to
	// renew a leadership slot before it stops leading. This must be less
	// than or equal to the lease duration. This is only applicable if leader
	// election is enabled.
	RenewDeadline time.Duration
	// retryPeriod is the duration the clients should wait between attempting
	// acquisition and renewal of a leadership. This is only applicable if
	// leader election is enabled.
	RetryPeriod time.Duration
}

// NewComponentOptions creates a ComponentOptions object with default parameters.
func NewComponentOptions(allControllers []string, disabledByDefaultControllers []string) *ComponentOptions {
	return &ComponentOptions{
		MinResyncPeriod:         12 * time.Hour,
		ControllerStartInterval: 0 * time.Second,
		LeaderElection: &LeaderElectionConfiguration{
			LeaderElect:   true,
			LeaseDuration: 15 * time.Second,
			RenewDeadline: 10 * time.Second,
			RetryPeriod:   2 * time.Second,
		},
		Controllers:                  []string{"*"},
		ContainerRegistryDomain:      "docker.io",
		allControllers:               allControllers,
		disabledByDefaultControllers: disabledByDefaultControllers,
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *ComponentOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Duration(flagComponentMinResyncPeriod, o.MinResyncPeriod,
		"The resync period in reflectors will be random between MinResyncPeriod and 2*MinResyncPeriod.")
	_ = viper.BindPFlag(configComponentMinResyncPeriod, fs.Lookup(flagComponentMinResyncPeriod))
	fs.Duration(flagComponentControllerStartInterval, o.ControllerStartInterval,
		"Interval between starting controller managers.")
	_ = viper.BindPFlag(configComponentControllerStartInterval, fs.Lookup(flagComponentControllerStartInterval))
	fs.StringSlice(flagComponentControllers, o.Controllers, fmt.Sprintf(""+
		"A list of controllers to enable. '*' enables all on-by-default controllers, 'foo' enables the controller "+
		"named 'foo', '-foo' disables the controller named 'foo'.\nAll controllers: %s\nDisabled-by-default controllers: %s",
		strings.Join(o.allControllers, ", "), strings.Join(o.disabledByDefaultControllers, ", ")))
	_ = viper.BindPFlag(configComponentControllers, fs.Lookup(flagComponentControllers))
	fs.Bool(flagComponentLeaderElect, o.LeaderElection.LeaderElect,
		"Start a leader election client and gain leadership before "+
			"executing the main loop. Enable this when running replicated "+
			"components for high availability.")
	_ = viper.BindPFlag(configComponentLeaderElect, fs.Lookup(flagComponentLeaderElect))
	fs.Duration(flagComponentLeaseDuration, o.LeaderElection.LeaseDuration,
		"The duration that non-leader candidates will wait after observing a leadership "+
			"renewal until attempting to acquire leadership of a led but unrenewed leader "+
			"slot. This is effectively the maximum duration that a leader can be stopped "+
			"before it is replaced by another candidate. This is only applicable if leader "+
			"election is enabled.")
	_ = viper.BindPFlag(configComponentLeaseDuration, fs.Lookup(flagComponentLeaseDuration))
	fs.Duration(flagComponentRenewDeadline, o.LeaderElection.RenewDeadline,
		"The interval between attempts by the acting master to renew a leadership slot "+
			"before it stops leading. This must be less than or equal to the lease duration. "+
			"This is only applicable if leader election is enabled.")
	_ = viper.BindPFlag(configComponentRenewDeadline, fs.Lookup(flagComponentRenewDeadline))
	fs.Duration(flagComponentRetryPeriod, o.LeaderElection.RetryPeriod,
		"The duration the clients should wait between attempting acquisition and renewal "+
			"of a leadership. This is only applicable if leader election is enabled.")
	_ = viper.BindPFlag(configComponentRetryPeriod, fs.Lookup(flagComponentRetryPeriod))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *ComponentOptions) ApplyFlags() []error {
	var errs []error

	o.LeaderElection.LeaderElect = viper.GetBool(configComponentLeaderElect)
	o.LeaderElection.LeaseDuration = viper.GetDuration(configComponentLeaseDuration)
	o.LeaderElection.RenewDeadline = viper.GetDuration(configComponentRenewDeadline)
	o.LeaderElection.RetryPeriod = viper.GetDuration(configComponentRetryPeriod)
	o.ControllerStartInterval = viper.GetDuration(configComponentControllerStartInterval)
	o.Controllers = viper.GetStringSlice(configComponentControllers)
	o.MinResyncPeriod = viper.GetDuration(configComponentMinResyncPeriod)

	errs = append(errs, o.Validate()...)

	return errs
}

// Validate checks validation of GenericOptions.
func (o *ComponentOptions) Validate() []error {
	var errs []error

	allControllersSet := sets.NewString(o.allControllers...)
	for _, controller := range o.Controllers {
		if controller == "*" {
			continue
		}
		controller = strings.TrimPrefix(controller, "-")
		if !allControllersSet.Has(controller) {
			errs = append(errs, fmt.Errorf("%q is not in the list of known controllers", controller))
		}
	}

	return errs
}

// ApplyTo parsing parameters from the command line or configuration file
// to the options instance.
func (o *ComponentOptions) ApplyTo(cfg *ComponentConfiguration) error {
	cfg.MinResyncPeriod = o.MinResyncPeriod
	cfg.Controllers = o.Controllers
	cfg.ControllerStartInterval = o.ControllerStartInterval
	cfg.LeaderElection.LeaderElect = o.LeaderElection.LeaderElect
	cfg.LeaderElection.RetryPeriod = metav1.Duration{Duration: o.LeaderElection.RetryPeriod}
	cfg.LeaderElection.RenewDeadline = metav1.Duration{Duration: o.LeaderElection.RenewDeadline}
	cfg.LeaderElection.LeaseDuration = metav1.Duration{Duration: o.LeaderElection.LeaseDuration}
	return nil
}
