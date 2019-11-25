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

// Portions Copyright 2014 The Kubernetes Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"strings"
	"testing"

	"k8s.io/client-go/util/flowcontrol"
)

func TestRegisterMetricAndTrackRateLimiterUsage(t *testing.T) {
	testCases := []struct {
		ownerName   string
		rateLimiter flowcontrol.RateLimiter
		err         string
	}{
		{
			ownerName:   "owner_name",
			rateLimiter: flowcontrol.NewTokenBucketRateLimiter(1, 1),
			err:         "",
		},
		{
			ownerName:   "owner_name",
			rateLimiter: flowcontrol.NewTokenBucketRateLimiter(1, 1),
			err:         "already registered",
		},
		{
			ownerName:   "invalid-owner-name",
			rateLimiter: flowcontrol.NewTokenBucketRateLimiter(1, 1),
			err:         "error registering rate limiter usage metric",
		},
	}

	for i, tc := range testCases {
		e := RegisterMetricAndTrackRateLimiterUsage(tc.ownerName, tc.rateLimiter)
		if e != nil {
			if tc.err == "" {
				t.Errorf("[%d] unexpected error: %v", i, e)
			} else if !strings.Contains(e.Error(), tc.err) {
				t.Errorf("[%d] expected an error containing %q: %v", i, tc.err, e)
			}
		}
	}
}
