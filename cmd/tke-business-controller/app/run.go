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

package app

import (
	"context"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apiserver/pkg/server/healthz"
	"tkestack.io/tke/api/business"
	"tkestack.io/tke/cmd/tke-business-controller/app/config"
	"tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util/leaderelection"
	"tkestack.io/tke/pkg/util/leaderelection/resourcelock"
	"tkestack.io/tke/pkg/util/log"
)

// Run runs the specified business controller manager. This should never exit.
func Run(cfg *config.Config, stopCh <-chan struct{}) error {
	log.Info("Starting Tencent Kubernetes Engine business controller manager")

	// Setup any healthz checks we will want to use.
	var checks []healthz.HealthChecker
	var electionChecker *leaderelection.HealthzAdaptor
	if cfg.Component.LeaderElection.LeaderElect {
		electionChecker = leaderelection.NewLeaderHealthzAdaptor(time.Second * 20)
		checks = append(checks, electionChecker)
	}

	// Start the controller manager HTTP server
	// serverMux is the handler for these controller *after* authn/authz filters have been applied
	serverMux := controller.NewBaseHandler(&cfg.Component.Debugging, checks...)
	handler := controller.BuildHandlerChain(serverMux, &cfg.Authorization, &cfg.Authentication, business.Codecs)
	if _, err := cfg.SecureServing.Serve(handler, 0, stopCh); err != nil {
		return err
	}

	run := func(ctx context.Context) {
		rootClientBuilder := controller.SimpleControllerClientBuilder{
			ClientConfig: cfg.BusinessAPIServerClientConfig,
		}

		controllerContext, err := CreateControllerContext(cfg, rootClientBuilder, ctx.Done())
		if err != nil {
			log.Fatalf("error building controller context: %v", err)
		}

		if err := StartControllers(controllerContext, NewControllerInitializers(), serverMux); err != nil {
			log.Fatalf("error starting controllers: %v", err)
		}

		controllerContext.InformerFactory.Start(controllerContext.Stop)
		close(controllerContext.InformersStarted)

		select {}
	}

	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		<-stopCh
		cancel()
	}()

	if !cfg.Component.LeaderElection.LeaderElect {
		run(ctx)
		panic("unreachable")
	}

	id, err := os.Hostname()
	if err != nil {
		return err
	}

	// add a uniquifier so that two processes on the same host don't accidentally both become active
	id = id + "_" + string(uuid.NewUUID())
	rl := resourcelock.NewBusiness("tke-business-controller",
		cfg.LeaderElectionClient.BusinessV1(),
		resourcelock.Config{
			Identity: id,
		})

	leaderelection.RunOrDie(ctx, leaderelection.ElectionConfig{
		Lock:          rl,
		LeaseDuration: cfg.Component.LeaderElection.LeaseDuration.Duration,
		RenewDeadline: cfg.Component.LeaderElection.RenewDeadline.Duration,
		RetryPeriod:   cfg.Component.LeaderElection.RetryPeriod.Duration,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				log.Fatalf("leaderelection lost")
			},
		},
		WatchDog: electionChecker,
		Name:     "tke-business-controller",
	})
	panic("unreachable")
}
