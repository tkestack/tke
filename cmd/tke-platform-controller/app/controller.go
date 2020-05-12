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
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server/mux"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// ControllerStartJitter used when starting controller managers.
	ControllerStartJitter = 1.0
)

// ControllersDisabledByDefault configured all controllers that are turned off
// by default.
var ControllersDisabledByDefault = sets.NewString()

// KnownControllers returns the known controllers.
func KnownControllers() []string {
	ret := sets.StringKeySet(NewControllerInitializers())
	return ret.List()
}

// NewControllerInitializers is a public map of named controller groups (you can start more than one in an init func)
// paired to their InitFunc.  This allows for structured downstream composition and subdivision.
func NewControllerInitializers() map[string]InitFunc {
	controllers := map[string]InitFunc{}

	controllers["cluster"] = startClusterController
	controllers["machine"] = startMachineController
	controllers["persistentevent"] = startPersistentEventController
	controllers["helm"] = startHelmController
	controllers["tappcontroller"] = startTappControllerController
	controllers["cronhpa"] = startCronHPAController
	controllers["csioperator"] = startCSIOperatorController
	controllers["volumedecorators"] = startVolumeDecoratorController
	controllers["logcollectors"] = startLogCollectorController
	controllers["prometheus"] = startPrometheusController
	controllers["ipam"] = startIPAMController
	controllers["lbcf"] = startLBCFControllerController
	return controllers
}

// StartControllers to start the controller.
func StartControllers(ctx ControllerContext, controllers map[string]InitFunc, unsecuredMux *mux.PathRecorderMux) error {
	for controllerName, initFn := range controllers {
		if !ctx.IsControllerEnabled(controllerName) {
			log.Warnf("%q is disabled", controllerName)
			continue
		}

		time.Sleep(wait.Jitter(ctx.ControllerStartInterval, ControllerStartJitter))

		log.Infof("Starting %q", controllerName)
		debugHandler, started, err := initFn(ctx)
		if err != nil {
			log.Errorf("Error starting %q", controllerName)
			return err
		}
		if !started {
			log.Warnf("Skipping %q", controllerName)
			continue
		}
		if debugHandler != nil && unsecuredMux != nil {
			basePath := "/debug/controllers/" + controllerName
			unsecuredMux.UnlistedHandle(basePath, http.StripPrefix(basePath, debugHandler))
			unsecuredMux.UnlistedHandlePrefix(basePath+"/", http.StripPrefix(basePath, debugHandler))
		}
		log.Infof("Started %q", controllerName)
	}

	return nil
}
