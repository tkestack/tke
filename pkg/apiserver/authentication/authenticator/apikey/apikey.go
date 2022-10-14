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

package apikey

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	v1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	"net/http"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/transport"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Options struct {
	TokenReviewCAFile string
	TokenReviewURL    string
	AdminUsername     string
	AdminPassword     string
}

// NewAPIKeyAuthenticator creates a request auth authenticator and returns it.
func NewAPIKeyAuthenticator(opts *Options) (authenticator.Request, error) {
	tr, err := transport.NewOneWayTLSTransport(opts.TokenReviewCAFile, true)
	if err != nil {
		return nil, err
	}

	return &Authenticator{
		adminUsername:        opts.AdminUsername,
		adminPassword:        opts.AdminPassword,
		tokenReviewURL:       opts.TokenReviewURL,
		tokenReviewTransport: tr,
	}, nil
}

// Authenticator according to the basic auth header information in the http
// request, treats the password as the APIKey of the OIDC server, requests the
// token review of the OIDC server, and returns user authentication information.
type Authenticator struct {
	adminUsername        string
	adminPassword        string
	tokenReviewURL       string
	tokenReviewTransport *http.Transport
}

// AuthenticateRequest implements authenticator.Request.
func (a *Authenticator) AuthenticateRequest(req *http.Request) (*authenticator.Response, bool, error) {
	username, password, ok := req.BasicAuth()
	if !ok {
		return nil, false, fmt.Errorf("cannot get user info from request")
	}
	if a.adminPassword != "" &&
		a.adminUsername != "" &&
		username == a.adminUsername &&
		password == a.adminPassword {
		log.Debug("Authenticated a system administrator")
		u := &user.DefaultInfo{
			Name:   a.adminUsername,
			UID:    "",
			Groups: nil,
			Extra:  make(map[string][]string, 0),
		}
		return &authenticator.Response{
			User: u,
		}, true, nil
	}
	if a.tokenReviewURL == "" {
		log.Warn("Token review url not specify, failed to review token")
		return nil, false, fmt.Errorf("token review url not specify")
	}
	log.Debug("Start review token", log.String("username", username), log.String("tokenReviewURL", a.tokenReviewURL))
	tokenReviewRequest := &v1.TokenReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TokenReview",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "tke",
		},
		Spec: v1.TokenReviewSpec{
			Token: password,
		},
	}
	bs, err := json.Marshal(tokenReviewRequest)
	if err != nil {
		log.Error("Failed to marshal token review request", log.Any("tokenReview", tokenReviewRequest), log.Err(err))
		return nil, false, err
	}

	req, err = http.NewRequest(http.MethodPost, a.tokenReviewURL, bytes.NewBuffer(bs))
	if err != nil {
		log.Error("Failed to create token review request", log.Err(err))
		return nil, false, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Transport: a.tokenReviewTransport,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error("Failed to request token review", log.Err(err))
	}

	defer func() {
		_ = res.Body.Close()
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to read response body", log.Err(err))
		return nil, false, err
	}
	if res.StatusCode != http.StatusOK {
		log.Error("Authentication server returns http status code error when token review",
			log.Int("statusCode", res.StatusCode),
			log.ByteString("responseBody", body))
		return nil, false, err
	}
	var tokenReviewResponse v1.TokenReview
	if err := json.Unmarshal(body, &tokenReviewResponse); err != nil {
		log.Error("Failed to unmarshal token review response", log.ByteString("body", body), log.Err(err))
		return nil, false, err
	}
	if tokenReviewResponse.Status.Error != "" {
		log.Error("Authentication server returns error when token review",
			log.String("error", tokenReviewResponse.Status.Error))
		return nil, false, fmt.Errorf(tokenReviewResponse.Status.Error)
	}
	if !tokenReviewResponse.Status.Authenticated {
		log.Error("Authentication server returns not authenticated when token review",
			log.Any("tokenReviewStatus", tokenReviewResponse.Status))
		return nil, false, fmt.Errorf("unknown error when token review")
	}

	u := &user.DefaultInfo{
		Name:   tokenReviewResponse.Status.User.Username,
		UID:    tokenReviewResponse.Status.User.UID,
		Groups: tokenReviewResponse.Status.User.Groups,
		Extra:  make(map[string][]string, 0),
	}
	if len(tokenReviewResponse.Status.User.Extra) > 0 {
		for k, v := range tokenReviewResponse.Status.User.Extra {
			u.Extra[k] = v
		}
	}
	return &authenticator.Response{
		User: u,
	}, true, nil
}
