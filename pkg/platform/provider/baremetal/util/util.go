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
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/rand"
	tkev1 "tkestack.io/tke/api/platform/v1"
)

func GetMasterEndpoint(addresses []tkev1.ClusterAddress) (string, error) {
	var advertise, internal []*tkev1.ClusterAddress
	for _, one := range addresses {
		if one.Type == tkev1.AddressAdvertise {
			advertise = append(advertise, &one)
		}
		if one.Type == tkev1.AddressReal {
			internal = append(internal, &one)
		}
	}

	var address *tkev1.ClusterAddress
	if advertise != nil {
		address = advertise[rand.Intn(len(advertise))]
	} else {
		if internal != nil {
			address = internal[rand.Intn(len(internal))]
		}
	}
	if address == nil {
		return "", errors.New("no advertise or internal address for the cluster")
	}

	return fmt.Sprintf("https://%s:%d", address.Host, address.Port), nil
}
