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

package storage

import (
	"context"
	"fmt"
	corev1api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericrest "k8s.io/apiserver/pkg/registry/generic/rest"
	"k8s.io/apiserver/pkg/registry/rest"
	"net/url"
	"strconv"
	"time"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
	restutil "tkestack.io/tke/pkg/platform/util/rest"
)

const isNegativeErrorMsg = apimachineryvalidation.IsNegativeErrorMsg

// LogREST implements the log endpoint for a Pod
type LogREST struct {
	platformClient platforminternalclient.PlatformInterface
}

var _ rest.StorageMetadata = &LogREST{}
var _ rest.GetterWithOptions = &LogREST{}

// NewGetOptions returns versioned resource that represents proxy parameters
func (r *LogREST) NewGetOptions() (runtime.Object, bool, string) {
	return &corev1api.PodLogOptions{}, false, ""
}

// OverrideMetricsVerb override the GET verb to CONNECT for pod log resource
func (r *LogREST) OverrideMetricsVerb(oldVerb string) (newVerb string) {
	newVerb = oldVerb

	if oldVerb == "GET" {
		newVerb = "CONNECT"
	}

	return
}

// New creates a new Pod log options object
func (r *LogREST) New() runtime.Object {
	return &corev1api.Pod{}
}

// ProducesMIMETypes implements StorageMetadata
func (r *LogREST) ProducesMIMETypes(verb string) []string {
	return []string{"text/plain"}
}

// ProducesObject implements StorageMetadata, return string as the generating object
func (r *LogREST) ProducesObject(verb string) interface{} {
	return ""
}

// Get retrieves a runtime.Object that will stream the contents of the pod log
func (r *LogREST) Get(ctx context.Context, name string, opts runtime.Object) (runtime.Object, error) {
	logOpts, ok := opts.(*corev1api.PodLogOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", opts)
	}
	if errs := validatePodLogOptions(logOpts); len(errs) > 0 {
		return nil, errors.NewInvalid(corev1api.SchemeGroupVersion.WithKind("PodLogOptions").GroupKind(), name, errs)
	}

	location, transport, token, err := util.APIServerLocation(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	if logOpts.Container != "" {
		params.Add("container", logOpts.Container)
	}
	if logOpts.Follow {
		params.Add("follow", "true")
	}
	if logOpts.Previous {
		params.Add("previous", "true")
	}
	if logOpts.Timestamps {
		params.Add("timestamps", "true")
	}
	if logOpts.SinceSeconds != nil {
		params.Add("sinceSeconds", strconv.FormatInt(*logOpts.SinceSeconds, 10))
	}
	if logOpts.SinceTime != nil {
		params.Add("sinceTime", logOpts.SinceTime.Format(time.RFC3339))
	}
	if logOpts.TailLines != nil {
		params.Add("tailLines", strconv.FormatInt(*logOpts.TailLines, 10))
	}
	if logOpts.LimitBytes != nil {
		params.Add("limitBytes", strconv.FormatInt(*logOpts.LimitBytes, 10))
	}

	location.RawQuery = params.Encode()

	return &restutil.LocationStreamer{
		Token:           token,
		Location:        location,
		Transport:       transport,
		ContentType:     "text/plain",
		Flush:           logOpts.Follow,
		ResponseChecker: genericrest.NewGenericHttpResponseChecker(corev1api.Resource("pods/log"), name),
		RedirectChecker: genericrest.PreventRedirects,
	}, nil
}

func validatePodLogOptions(opts *corev1api.PodLogOptions) field.ErrorList {
	allErrs := field.ErrorList{}
	if opts.TailLines != nil && *opts.TailLines < 0 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("tailLines"), *opts.TailLines, isNegativeErrorMsg))
	}
	if opts.LimitBytes != nil && *opts.LimitBytes < 1 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("limitBytes"), *opts.LimitBytes, "must be greater than 0"))
	}
	switch {
	case opts.SinceSeconds != nil && opts.SinceTime != nil:
		allErrs = append(allErrs, field.Forbidden(field.NewPath(""), "at most one of `sinceTime` or `sinceSeconds` may be specified"))
	case opts.SinceSeconds != nil:
		if *opts.SinceSeconds < 1 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("sinceSeconds"), *opts.SinceSeconds, "must be greater than 0"))
		}
	}
	return allErrs
}
