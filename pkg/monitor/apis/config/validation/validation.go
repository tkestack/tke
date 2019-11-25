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

package validation

import (
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
)

// ValidateMonitorConfiguration validates `mc` and returns an error if it is invalid
func ValidateMonitorConfiguration(mc *monitorconfig.MonitorConfiguration) error {
	var allErrors []error

	fld := field.NewPath("storage")
	storageCount := 0
	if mc.Storage.InfluxDB != nil {
		storageCount++

		influxDBFld := fld.Child("influxDB")
		if len(mc.Storage.InfluxDB.Servers) == 0 {
			allErrors = append(allErrors, field.Required(influxDBFld.Child("servers"), ""))
		} else {
			for index, v := range mc.Storage.InfluxDB.Servers {
				if v.Address == "" {
					allErrors = append(allErrors, field.Required(influxDBFld.Child("servers").Index(index).Child("address"), "must be specify"))
				}
			}
		}
	}

	if mc.Storage.ElasticSearch != nil {
		storageCount++

		esFld := fld.Child("elasticSearch")
		if len(mc.Storage.ElasticSearch.Servers) == 0 {
			allErrors = append(allErrors, field.Required(esFld.Child("servers"), ""))
		} else {
			for index, v := range mc.Storage.ElasticSearch.Servers {
				if v.Address == "" {
					allErrors = append(allErrors, field.Required(esFld.Child("servers").Index(index).Child("address"), "must be specify"))
				}
			}
		}
	}

	if storageCount == 0 {
		allErrors = append(allErrors, field.Required(fld, "at least 1 storage is required"))
	} else if storageCount > 1 {
		allErrors = append(allErrors, field.Required(fld, "storage can only specify at most one"))
	}

	return utilerrors.NewAggregate(allErrors)
}
