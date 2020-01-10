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

package context

import (
	"context"
	dcontext "github.com/docker/distribution/context"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/log/distribution"
)

// BuildDistributionContext create a background context with logger for
// distribution and returns it.
func BuildDistributionContext() context.Context {
	ctx := context.Background()
	return BuildRequestContext(ctx)
}

// BuildRequestContext create a new context with logger by given context.
func BuildRequestContext(ctx context.Context) context.Context {
	logger := distribution.NewLogger(log.ZapLogger())
	return dcontext.WithLogger(ctx, logger)
}
