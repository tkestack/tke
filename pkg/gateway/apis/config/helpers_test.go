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

package config

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"reflect"
	"strings"
	"testing"
)

func TestGatewayConfigurationPathFields(t *testing.T) {
	// ensure the intersection of gatewayConfigurationPathFieldPaths and GatewayConfigurationNonPathFields is empty
	if i := gatewayConfigurationPathFieldPaths.Intersection(gatewayConfigurationNonPathFieldPaths); len(i) > 0 {
		t.Fatalf("expect the intersection of gatewayConfigurationPathFieldPaths and "+
			"GatewayConfigurationNonPathFields to be empty, got:\n%s",
			strings.Join(i.List(), "\n"))
	}

	// ensure that gatewayConfigurationPathFields U gatewayConfigurationNonPathFields == allPrimitiveFieldPaths(GatewayConfiguration)
	expect := sets.NewString().Union(gatewayConfigurationPathFieldPaths).Union(gatewayConfigurationNonPathFieldPaths)
	result := allPrimitiveFieldPaths(t, reflect.TypeOf(&GatewayConfiguration{}), nil)
	if !expect.Equal(result) {
		// expected fields missing from result
		missing := expect.Difference(result)
		// unexpected fields in result but not specified in expect
		unexpected := result.Difference(expect)
		if len(missing) > 0 {
			t.Errorf("the following fields were expected, but missing from the result. "+
				"If the field has been removed, please remove it from the gatewayConfigurationPathFieldPaths set "+
				"and the GatewayConfigurationPathRefs function, "+
				"or remove it from the gatewayConfigurationNonPathFieldPaths set, as appropriate:\n%s",
				strings.Join(missing.List(), "\n"))
		}
		if len(unexpected) > 0 {
			t.Errorf("the following fields were in the result, but unexpected. "+
				"If the field is new, please add it to the gatewayConfigurationPathFieldPaths set "+
				"and the GatewayConfigurationPathRefs function, "+
				"or add it to the gatewayConfigurationNonPathFieldPaths set, as appropriate:\n%s",
				strings.Join(unexpected.List(), "\n"))
		}
	}
}

func allPrimitiveFieldPaths(t *testing.T, tp reflect.Type, path *field.Path) sets.String {
	paths := sets.NewString()
	switch tp.Kind() {
	case reflect.Ptr:
		paths.Insert(allPrimitiveFieldPaths(t, tp.Elem(), path).List()...)
	case reflect.Struct:
		for i := 0; i < tp.NumField(); i++ {
			f := tp.Field(i)
			paths.Insert(allPrimitiveFieldPaths(t, f.Type, path.Child(f.Name)).List()...)
		}
	case reflect.Map, reflect.Slice:
		paths.Insert(allPrimitiveFieldPaths(t, tp.Elem(), path.Key("*")).List()...)
	case reflect.Interface:
		t.Fatalf("unexpected interface{} field %s", path.String())
	default:
		// if we hit a primitive type, we're at a leaf
		paths.Insert(path.String())
	}
	return paths
}

// dummy helper types
type foo struct {
	foo int
}
type bar struct {
	str    string
	strptr *string

	ints      []int
	stringMap map[string]string

	foo    foo
	fooptr *foo

	bars   []foo
	barMap map[string]foo
}

func TestAllPrimitiveFieldPaths(t *testing.T) {
	expect := sets.NewString(
		"str",
		"strptr",
		"ints[*]",
		"stringMap[*]",
		"foo.foo",
		"fooptr.foo",
		"bars[*].foo",
		"barMap[*].foo",
	)
	result := allPrimitiveFieldPaths(t, reflect.TypeOf(&bar{}), nil)
	if !expect.Equal(result) {
		// expected fields missing from result
		missing := expect.Difference(result)

		// unexpected fields in result but not specified in expect
		unexpected := result.Difference(expect)

		if len(missing) > 0 {
			t.Errorf("the following fields were exepcted, but missing from the result:\n%s", strings.Join(missing.List(), "\n"))
		}
		if len(unexpected) > 0 {
			t.Errorf("the following fields were in the result, but unexpected:\n%s", strings.Join(unexpected.List(), "\n"))
		}
	}
}

var (
	// GatewayConfiguration fields that contain file paths. If you update this, also update GatewayConfigurationPathRefs!
	gatewayConfigurationPathFieldPaths = sets.NewString(
		"Components.Auth.FrontProxy.CAFile",
		"Components.Auth.FrontProxy.ClientCertFile",
		"Components.Auth.FrontProxy.ClientKeyFile",
		"Components.Business.FrontProxy.CAFile",
		"Components.Business.FrontProxy.ClientCertFile",
		"Components.Business.FrontProxy.ClientKeyFile",
		"Components.Monitor.FrontProxy.CAFile",
		"Components.Monitor.FrontProxy.ClientCertFile",
		"Components.Monitor.FrontProxy.ClientKeyFile",
		"Components.Notify.FrontProxy.CAFile",
		"Components.Notify.FrontProxy.ClientCertFile",
		"Components.Notify.FrontProxy.ClientKeyFile",
		"Components.Platform.FrontProxy.CAFile",
		"Components.Platform.FrontProxy.ClientCertFile",
		"Components.Platform.FrontProxy.ClientKeyFile",
		"Components.Registry.FrontProxy.CAFile",
		"Components.Registry.FrontProxy.ClientCertFile",
		"Components.Registry.FrontProxy.ClientKeyFile",
		"Components.Auth.Passthrough.CAFile",
		"Components.Business.Passthrough.CAFile",
		"Components.Monitor.Passthrough.CAFile",
		"Components.Notify.Passthrough.CAFile",
		"Components.Platform.Passthrough.CAFile",
		"Components.Registry.Passthrough.CAFile",
	)

	// GatewayConfiguration fields that do not contain file paths.
	gatewayConfigurationNonPathFieldPaths = sets.NewString(
		"TypeMeta.APIVersion",
		"TypeMeta.Kind",
		"DisableOIDCProxy",
		"Components.Auth.Address",
		"Components.Auth.FrontProxy.UsernameHeader",
		"Components.Auth.FrontProxy.GroupsHeader",
		"Components.Auth.FrontProxy.ExtraPrefixHeader",
		"Components.Business.Address",
		"Components.Business.FrontProxy.UsernameHeader",
		"Components.Business.FrontProxy.GroupsHeader",
		"Components.Business.FrontProxy.ExtraPrefixHeader",
		"Components.Monitor.Address",
		"Components.Monitor.FrontProxy.UsernameHeader",
		"Components.Monitor.FrontProxy.GroupsHeader",
		"Components.Monitor.FrontProxy.ExtraPrefixHeader",
		"Components.Notify.Address",
		"Components.Notify.FrontProxy.UsernameHeader",
		"Components.Notify.FrontProxy.GroupsHeader",
		"Components.Notify.FrontProxy.ExtraPrefixHeader",
		"Components.Platform.Address",
		"Components.Platform.FrontProxy.UsernameHeader",
		"Components.Platform.FrontProxy.GroupsHeader",
		"Components.Platform.FrontProxy.ExtraPrefixHeader",
		"Components.Registry.Address",
		"Components.Registry.FrontProxy.UsernameHeader",
		"Components.Registry.FrontProxy.GroupsHeader",
		"Components.Registry.FrontProxy.ExtraPrefixHeader",
		"Registry.DefaultTenant",
		"Registry.DomainSuffix",
		"Auth.DefaultTenant",
	)
)
