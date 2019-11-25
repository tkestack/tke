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

package authenticator

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"tkestack.io/tke/pkg/auth/types"

	"github.com/coreos/etcd/clientv3"
	"github.com/dgrijalva/jwt-go"
	"gotest.tools/assert"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/auth/registry/apikey"
	"tkestack.io/tke/pkg/util/etcd"
	"tkestack.io/tke/pkg/util/log"
)

var (
	defaultDialTimeout = 2 * time.Second

	client *clientv3.Client

	store apikey.Storage

	userName = "test"

	tenantID = "default"

	desc = "fortest"
)

func setup(t *testing.T) {
	testEtcdEnv := "ETCD_ENDPOINTS"
	endpointsStr := os.Getenv(testEtcdEnv)
	if endpointsStr == "" {
		t.Skipf("test environment variable %q not set, skipping", testEtcdEnv)
		return
	}
	endpoints := strings.Split(endpointsStr, ",")
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: defaultDialTimeout,
	}

	var err error
	client, err = clientv3.New(cfg)
	if err != nil {
		t.Skipf("init etcd client %v failed, skipping", endpoints)
	}

	store = apikey.NewAPIKeyStorage(client)
}

func cleanup(t *testing.T) {
	if _, err := client.Delete(context.Background(), "/signkeys/", clientv3.WithPrefix()); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Delete(context.Background(), "/apikeys/", clientv3.WithPrefix()); err != nil {
		t.Fatal(err)
	}
}

func TestNewAPIKeyAuthRSAAndVerify(t *testing.T) {
	setup(t)
	defer cleanup(t)
	apiKeyAuth, err := NewAPIKeyAuthenticator("RSA", store, nil)
	if err != nil {
		t.Fatal(err)
	}

	apiKey, err := apiKeyAuth.CreateToken(tenantID, userName, desc, 2*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	resp, valid, err := apiKeyAuth.AuthenticateToken(context.Background(), apiKey.APIkey)
	if !valid || err != nil {
		t.Fatalf("auth api key faild: %v, %v", valid, err)
	}

	assert.Assert(t, resp.User.GetName() == userName)
	assert.Assert(t, resp.User.GetExtra()[oidc.TenantIDKey][0] == tenantID)
	assert.Assert(t, resp.User.GetExtra()["description"][0] == desc)

	//Test api key expired
	at(time.Now().Add(3*time.Hour), func() {
		resp, valid, err = apiKeyAuth.AuthenticateToken(context.Background(), apiKey.APIkey)
		if err == nil || err.Error() != "token is either expired or not active yet" {
			t.Fatal("token is expected expire")
		}
	})

	token := "12345"
	_, valid, err = apiKeyAuth.AuthenticateToken(context.Background(), token)
	if valid || err == nil {
		t.Fatal("expected verify failed, but success")
	}
}

func TestNewAPIKeyAuthHMACAndVerify(t *testing.T) {
	setup(t)
	defer cleanup(t)
	apiKeyAuth, err := NewAPIKeyAuthenticator("HMAC", store, nil)
	if err != nil {
		t.Fatal(err)
	}

	apiKey, err := apiKeyAuth.CreateToken(tenantID, userName, desc, 2*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	resp, valid, err := apiKeyAuth.AuthenticateToken(context.Background(), apiKey.APIkey)
	if !valid || err != nil {
		t.Fatalf("auth api key faild: %v, %v", valid, err)
	}

	assert.Assert(t, resp.User.GetName() == userName)
	assert.Assert(t, resp.User.GetExtra()[oidc.TenantIDKey][0] == tenantID)

	//Test api key expired
	at(time.Now().Add(3*time.Hour), func() {
		resp, valid, err = apiKeyAuth.AuthenticateToken(context.Background(), apiKey.APIkey)
		if err == nil || err.Error() != "token is either expired or not active yet" {
			t.Fatal("token is expected expire")
		}
	})

	token := "12345"
	_, valid, err = apiKeyAuth.AuthenticateToken(context.Background(), token)
	if valid || err == nil {
		t.Fatal("expected verify failed, but success")
	}
}

func TestNewAPIKeyAuthWithoutExpire(t *testing.T) {
	setup(t)
	defer cleanup(t)
	apiKeyAuth, err := NewAPIKeyAuthenticator("HMAC", store, nil)
	if err != nil {
		t.Fatal(err)
	}

	apiKey, err := apiKeyAuth.CreateToken(tenantID, userName, desc, 0)
	if err != nil {
		t.Fatal(err)
	}

	resp, valid, err := apiKeyAuth.AuthenticateToken(context.Background(), apiKey.APIkey)
	if !valid || err != nil {
		t.Fatalf("auth api key faild: %v, %v", valid, err)
	}

	assert.Assert(t, resp.User.GetName() == userName)
	assert.Assert(t, resp.User.GetExtra()[oidc.TenantIDKey][0] == tenantID)
	// default is 7 days
	assert.Assert(t, apiKey.ExpireAt.After(time.Now().Add(6*24*time.Hour)))
}

func TestInvalidApiKey(t *testing.T) {
	setup(t)
	defer cleanup(t)
	apiKeyAuth, err := NewAPIKeyAuthenticator("HMAC", store, nil)
	if err != nil {
		t.Fatal(err)
	}

	apiKey, err := apiKeyAuth.CreateToken(tenantID, userName, desc, 2*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	value := true
	apiKey.Deleted = &value
	err = apiKeyAuth.UpdateToken(apiKey, tenantID, userName)
	if err != nil {
		t.Fatal(err)
	}

	_, valid, err := apiKeyAuth.AuthenticateToken(context.Background(), apiKey.APIkey)
	if valid || err == nil {
		t.Fatal("expected verify failed, but success")
	}

	apiKey, err = apiKeyAuth.CreateToken(tenantID, userName, desc, 3*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	apiKey.Disabled = &value
	err = apiKeyAuth.UpdateToken(apiKey, tenantID, userName)
	if err != nil {
		t.Fatal(err)
	}

	_, valid, err = apiKeyAuth.AuthenticateToken(context.Background(), apiKey.APIkey)
	if valid || err == nil {
		t.Fatal("expected verify failed, but success")
	}

	apiKeyAuth, err = NewAPIKeyAuthenticator("RSA", store, nil)
	if err != nil {
		t.Fatal(err)
	}

	apiKey, err = apiKeyAuth.CreateToken(tenantID, userName, desc, 2*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	apiKey.Disabled = &value
	err = apiKeyAuth.UpdateToken(apiKey, tenantID, userName)
	if err != nil {
		t.Fatal(err)
	}

	_, valid, err = apiKeyAuth.AuthenticateToken(context.Background(), apiKey.APIkey)
	if valid || err == nil {
		t.Fatal("expected verify failed, but success")
	}

	apiKey, err = apiKeyAuth.CreateToken(tenantID, userName, desc, 3*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	apiKey.Deleted = &value
	err = apiKeyAuth.UpdateToken(apiKey, tenantID, userName)
	if err != nil {
		t.Fatal(err)
	}

	_, valid, err = apiKeyAuth.AuthenticateToken(context.Background(), apiKey.APIkey)
	if valid || err == nil {
		t.Fatal("expected verify failed, but success")
	}
}

func TestInvalidApiKeyRotate(t *testing.T) {
	setup(t)
	defer cleanup(t)

	apiKeyRotateInterval = 1 * time.Second
	apiKeyAuth, err := NewAPIKeyAuthenticator("HMAC", store, nil)
	if err != nil {
		t.Fatal(err)
	}

	apiKey, err := apiKeyAuth.CreateToken(tenantID, userName, desc, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	value := true
	apiKey.Deleted = &value
	err = apiKeyAuth.UpdateToken(apiKey, tenantID, userName)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(4 * time.Second)
	res, err := apiKeyAuth.store.GetAPIKey(tenantID, userName, apiKey.APIkey)
	if err == nil || err != etcd.ErrNotFound {
		t.Fatal("expect get no apikey", log.Any("res", res), log.Err(err))
	}
}

func TestUniqApiKey(t *testing.T) {
	setup(t)
	defer cleanup(t)
	for i := 0; i < 10; i++ {
		var (
			apiKey1 *types.APIKeyData
			apiKey2 *types.APIKeyData
		)

		apiKeyAuth, err := NewAPIKeyAuthenticator("HMAC", store, nil)
		if err != nil {
			t.Fatal(err)
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			apiKey1, err = apiKeyAuth.CreateToken(tenantID, userName, desc, 2*time.Hour)
			if err != nil {
				t.Error(err)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			apiKey2, err = apiKeyAuth.CreateToken(tenantID, userName, desc, 2*time.Hour)
			if err != nil {
				t.Error(err)
			}
		}()
		wg.Wait()

		assert.Assert(t, apiKey1 != nil && apiKey2 != nil && apiKey1.APIkey != apiKey2.APIkey)
	}
}

func TestListAPIKeys(t *testing.T) {
	setup(t)
	defer cleanup(t)
	apiKeyAuth, err := NewAPIKeyAuthenticator("HMAC", store, nil)
	if err != nil {
		t.Fatal(err)
	}

	apiKey1, err := apiKeyAuth.CreateToken(tenantID, userName, desc, 2*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	apiKey2, err := apiKeyAuth.CreateToken(tenantID, userName, desc, 2*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	keyList, err := apiKeyAuth.ListAPIKeys(tenantID, userName)
	if err != nil {
		t.Fatal(err)
	}

	if len(keyList.Items) != 2 {
		t.Fatalf("Expected get 2 api key values, but got %+v", len(keyList.Items))
	}

	keySets := sets.NewString()

	for _, key := range keyList.Items {
		keySets.Insert(key.APIkey)
	}

	assert.Assert(t, keySets.Has(apiKey1.APIkey))
	assert.Assert(t, keySets.Has(apiKey2.APIkey))
}

// Override jwt time value for tests.  Restore default value after.
func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}
