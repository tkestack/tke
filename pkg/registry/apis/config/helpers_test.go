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

func TestRegistryConfigurationPathFields(t *testing.T) {
	// ensure the intersection of registryConfigurationPathFieldPaths and RegistryConfigurationNonPathFields is empty
	if i := registryConfigurationPathFieldPaths.Intersection(registryConfigurationNonPathFieldPaths); len(i) > 0 {
		t.Fatalf("expect the intersection of registryConfigurationPathFieldPaths and "+
			"RegistryConfigurationNonPathFields to be empty, got:\n%s",
			strings.Join(i.List(), "\n"))
	}

	// ensure that registryConfigurationPathFields U registryConfigurationNonPathFields == allPrimitiveFieldPaths(RegistryConfiguration)
	expect := sets.NewString().Union(registryConfigurationPathFieldPaths).Union(registryConfigurationNonPathFieldPaths)
	result := allPrimitiveFieldPaths(t, reflect.TypeOf(&RegistryConfiguration{}), nil)
	if !expect.Equal(result) {
		// expected fields missing from result
		missing := expect.Difference(result)
		// unexpected fields in result but not specified in expect
		unexpected := result.Difference(expect)
		if len(missing) > 0 {
			t.Errorf("the following fields were expected, but missing from the result. "+
				"If the field has been removed, please remove it from the registryConfigurationPathFieldPaths set "+
				"and the RegistryConfigurationPathRefs function, "+
				"or remove it from the registryConfigurationNonPathFieldPaths set, as appropriate:\n%s",
				strings.Join(missing.List(), "\n"))
		}
		if len(unexpected) > 0 {
			t.Errorf("the following fields were in the result, but unexpected. "+
				"If the field is new, please add it to the registryConfigurationPathFieldPaths set "+
				"and the RegistryConfigurationPathRefs function, "+
				"or add it to the registryConfigurationNonPathFieldPaths set, as appropriate:\n%s",
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
	// RegistryConfiguration fields that contain file paths. If you update this, also update RegistryConfigurationPathRefs!
	registryConfigurationPathFieldPaths = sets.NewString(
		"Security.TokenPrivateKeyFile",
		"Security.TokenPublicKeyFile",
	)

	// RegistryConfiguration fields that do not contain file paths.
	registryConfigurationNonPathFieldPaths = sets.NewString(
		"TypeMeta.APIVersion",
		"TypeMeta.Kind",
		"Security.TokenExpiredHours",
		"Security.HTTPSecret",
		"Security.AdminUsername",
		"Security.AdminPassword",
		"DomainSuffix",
		"DefaultTenant",
		"Storage.FileSystem.MaxThreads",
		"Storage.FileSystem.RootDirectory",
		"Storage.S3.AccessKey",
		"Storage.S3.Bucket",
		"Storage.S3.ChunkSize",
		"Storage.S3.Encrypt",
		"Storage.S3.KeyID",
		"Storage.S3.MultipartCopyChunkSize",
		"Storage.S3.MultipartCopyMaxConcurrency",
		"Storage.S3.MultipartCopyThresholdSize",
		"Storage.S3.ObjectACL",
		"Storage.S3.Region",
		"Storage.S3.RegionEndpoint",
		"Storage.S3.RootDirectory",
		"Storage.S3.SecretKey",
		"Storage.S3.Secure",
		"Storage.S3.SkipVerify",
		"Storage.S3.StorageClass",
		"Storage.S3.UserAgent",
		"Storage.S3.V4Auth",
		"Redis.Addr",
		"Redis.DB",
		"Redis.DialTimeoutMillisecond",
		"Redis.Password",
		"Redis.PoolIdleTimeoutSeconds",
		"Redis.PoolMaxActive",
		"Redis.PoolMaxIdle",
		"Redis.ReadTimeoutMillisecond",
		"Redis.WriteTimeoutMillisecond",
	)
)
