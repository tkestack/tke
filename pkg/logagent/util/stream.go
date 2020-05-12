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
	"context"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericrest "k8s.io/apiserver/pkg/registry/generic/rest"
	"k8s.io/apiserver/pkg/registry/rest"
	"net/http"
	"net/url"
)


// LocationStreamer is a resource that streams the contents of a particular
// location URL.
type LocationStreamer struct {
	Ip              string
	Request         ReaderCloserGetter
	Location        *url.URL
	Transport       http.RoundTripper
	ContentType     string
	Flush           bool
	ResponseChecker genericrest.HttpResponseChecker
	RedirectChecker func(req *http.Request, via []*http.Request) error
}

// a LocationStreamer must implement a rest.ResourceStreamer
var _ rest.ResourceStreamer = &LocationStreamer{}

func (obj *LocationStreamer) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}
func (obj *LocationStreamer) DeepCopyObject() runtime.Object {
	panic("rest.LocationStreamer does not implement DeepCopyObject")
}

// InputStream returns a stream with the contents of the URL location. If no location is provided,
// a null stream is returned.
func (s *LocationStreamer) InputStream(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, contentType string, err error) {
	stream, err = s.Request.GetReaderCloser()
	flush = false
	contentType = s.ContentType
	return
}
