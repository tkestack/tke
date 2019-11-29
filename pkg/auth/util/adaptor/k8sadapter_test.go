/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package adapter

import (
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
	"testing"
	"time"

	fakeauth "tkestack.io/tke/api/client/clientset/versioned/fake"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
)

var defaultModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub)  && keyMatchCustom(r.obj, p.obj) && keyMatchCustom(r.act, p.act)
`

func testGetPolicy(t *testing.T, e *casbin.SyncedEnforcer, res [][]string) {
	t.Helper()
	myRes := e.GetPolicy()

	if !util.Array2DEquals(res, myRes) {
		t.Error("Test failed, Policy: ", myRes, ", supposed to be ", res)
		return
	}

	t.Log("Test pass")
}

func TestInitPolicy(t *testing.T) {

	authClient := fakeauth.NewSimpleClientset()
	sharedInformers := versionedinformers.NewSharedInformerFactory(authClient, 10*time.Second)
	a := NewAdapter(authClient, sharedInformers.Auth().V1().Rules(), "key")
	m := casbin.NewModel(defaultModel)
	e, err := casbin.NewSyncedEnforcerSafe(m, a)

	e.AddPolicySafe("alice", "data1", "read")

	err = a.SavePolicy(e.GetModel())
	if err != nil {
		panic(err)
	}

	// Clear the current policy.
	e.ClearPolicy()
	testGetPolicy(t, e, [][]string{})

	// Load the policy from ETCD.
	err = a.LoadPolicy(e.GetModel())
	if err != nil {
		panic(err)
	}
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

//func testSaveLoad(t *testing.T, pathKey string, etcdEndpoints []string) {
//	// Initialize some policy in ETCD.
//	initPolicy(t, pathKey, etcdEndpoints)
//	// Note: you don't need to look at the above code
//	// if you already have a working ETCD with policy inside.
//
//	// Now the ETCD has policy, so we can provide a normal use case.
//	// Create an adapter and an enforcer.
//	// NewEnforcer() will load the policy automatically.
//	a := NewAdapter(etcdEndpoints, pathKey)
//	e := casbin.NewEnforcer("examples/rbac_model.conf", a)
//	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
//}
//
//func testAutoSave(t *testing.T, pathKey string, etcdEndpoints []string) {
//	// Initialize some policy in ETCD.
//	initPolicy(t, pathKey, etcdEndpoints)
//	// Note: you don't need to look at the above code
//	// if you already have a working ETCD with policy inside.
//
//	// Now the ETCD has policy, so we can provide a normal use case.
//	// Create an adapter and an enforcer.
//	// NewEnforcer() will load the policy automatically.
//	a := NewAdapter(etcdEndpoints, pathKey)
//	e := casbin.NewEnforcer("examples/rbac_model.conf", a)
//
//	// AutoSave is enabled by default.
//	// Now we disable it.
//	e.EnableAutoSave(false)
//
//	// Because AutoSave is disabled, the policy change only affects the policy in Casbin enforcer,
//	// it doesn't affect the policy in the storage.
//	e.AddPolicy("alice", "data1", "write")
//	// Reload the policy from the storage to see the effect.
//	e.LoadPolicy()
//	// This is still the original policy.
//	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
//
//	// Now we enable the AutoSave.
//	e.EnableAutoSave(true)
//
//	// Because AutoSave is enabled, the policy change not only affects the policy in Casbin enforcer,
//	// but also affects the policy in the storage.
//	e.AddPolicy("alice", "data1", "write")
//	// Reload the policy from the storage to see the effect.
//	e.LoadPolicy()
//	// The policy has a new rule: {"alice", "data1", "write"}.
//	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"alice", "data1", "write"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
//
//	// Remove the added rule.
//	e.RemovePolicy("alice", "data1", "write")
//	e.LoadPolicy()
//	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
//
//	// Remove "data2_admin" related policy rules via a filter.
//	// Two rules: {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"} will be deleted.
//	e.RemoveFilteredPolicy(0, "data2_admin")
//	e.LoadPolicy()
//	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}})
//
//}

func TestAdapters(t *testing.T) {

//	testSaveLoad(t, "casbin_policy_test", []string{"http://127.0.0.1:2379"})
//	testAutoSave(t, "casbin_policy_test", []string{"http://127.0.0.1:2379"})
}
