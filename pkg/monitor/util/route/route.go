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

package route

import (
	"github.com/pkg/errors"
	alertconfig "github.com/prometheus/alertmanager/config"
	"gopkg.in/yaml.v2"
	"io"
	"tkestack.io/tke/pkg/monitor/util"
	"tkestack.io/tke/pkg/util/log"
)

type route struct {
	alertConfig alertconfig.Config
}

var _ util.GenericRouteOperator = &route{}

// NewRouteOperator returns a implements of GenericRouteOperator to operate alertmanager config
func NewRouteOperator() util.GenericRouteOperator {
	return &route{}
}

func (r *route) InsertRoute(route *alertconfig.Route) (*alertconfig.Route, error) {
	if err := util.ValidateLabels(route.Match); err != nil {
		return nil, err
	}

	insertUniqValue := route.Match[util.DefaultAlertKey]
	for _, entry := range r.alertConfig.Route.Routes {
		if insertUniqValue == entry.Match[util.DefaultAlertKey] {
			return nil, errors.Errorf("duplicate value for %s", util.DefaultAlertKey)
		}
	}

	r.alertConfig.Route.Routes = append(r.alertConfig.Route.Routes, route)

	return r.alertConfig.Route, nil
}

func (r *route) DeleteRoute(alertValue string) (*alertconfig.Route, error) {
	for i, entry := range r.alertConfig.Route.Routes {
		if entry.Match[util.DefaultAlertKey] == alertValue {
			rebuildRoutes := make([]*alertconfig.Route, 0)
			rebuildRoutes = append(rebuildRoutes, r.alertConfig.Route.Routes[:i]...)
			if i < len(r.alertConfig.Route.Routes)-1 {
				rebuildRoutes = append(rebuildRoutes, r.alertConfig.Route.Routes[i+1:]...)
			}
			r.alertConfig.Route.Routes = rebuildRoutes
			return entry, nil
		}
	}

	return nil, errors.Errorf("%s=%s label not found", util.DefaultAlertKey, alertValue)
}

func (r *route) UpdateRoute(alertValue string, route *alertconfig.Route) (*alertconfig.Route, error) {
	if err := util.ValidateLabels(route.Match); err != nil {
		return nil, err
	}

	for i, entry := range r.alertConfig.Route.Routes {
		if entry.Match[util.DefaultAlertKey] == alertValue {
			r.alertConfig.Route.Routes[i] = route
			return route, nil
		}
	}

	return nil, errors.Errorf("%s=%s label not found", util.DefaultAlertKey, alertValue)
}

func (r *route) GetRoute(alertValue string) (*alertconfig.Route, error) {
	for _, entry := range r.alertConfig.Route.Routes {
		if entry.Match[util.DefaultAlertKey] == alertValue {
			return entry, nil
		}
	}

	return nil, errors.Errorf("%s=%s label not found", util.DefaultAlertKey, alertValue)
}

func (r *route) ListRoute() ([]*alertconfig.Route, error) {
	return r.alertConfig.Route.Routes, nil
}

func (r *route) Load(reader io.Reader) error {
	err := yaml.NewDecoder(reader).Decode(&r.alertConfig)
	if err != nil {
		log.Error("Failed to decode alert config file", log.Err(err))
		return err
	}
	return nil
}

// Save implements PersistentOperator
func (r *route) Save(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(&r.alertConfig)
}
