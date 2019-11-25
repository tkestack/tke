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

package apiclient

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	kuberuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	"tkestack.io/tke/pkg/util/template"
)

type object struct {
	Kind string `yaml:"kind"`
}

var handlers map[string]func(kubernetes.Interface, []byte) error

func init() {
	handlers = make(map[string]func(kubernetes.Interface, []byte) error)

	// core
	handlers["ConfigMap"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(corev1.ConfigMap)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateConfigMap(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Endpoints"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(corev1.Endpoints)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateEndpoints(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Namespace"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(corev1.Namespace)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateNamespace(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Secret"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(corev1.Secret)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateSecret(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Service"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(corev1.Service)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateService(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["ServiceAccount"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(corev1.ServiceAccount)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateServiceAccount(client, obj)
		if err != nil {
			return err
		}
		return nil
	}

	// apps
	handlers["DaemonSet"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(appsv1.DaemonSet)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateDaemonSet(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Deployment"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(appsv1.Deployment)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateDeployment(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["StatefulSet"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(appsv1.StatefulSet)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateStatefulSet(client, obj)
		if err != nil {
			return err
		}
		return nil
	}

	// extentions
	handlers["Ingress"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(extensionsv1beta1.Ingress)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateIngress(client, obj)
		if err != nil {
			return err
		}
		return nil
	}

	// rbac
	handlers["Role"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(rbacv1.Role)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateRole(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["RoleBinding"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(rbacv1.RoleBinding)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateRoleBinding(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["ClusterRole"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(rbacv1.ClusterRole)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateClusterRole(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["ClusterRoleBinding"] = func(client kubernetes.Interface, data []byte) error {
		obj := new(rbacv1.ClusterRoleBinding)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateClusterRoleBinding(client, obj)
		if err != nil {
			return err
		}
		return nil
	}
}

// CreateResourceWithDir create k8s resource with dir
func CreateResourceWithDir(client kubernetes.Interface, pattern string, option interface{}) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return errors.New("no matches found")
	}
	for _, filename := range matches {
		err := CreateResourceWithFile(client, filename, option)
		if err != nil {
			return errors.Wrapf(err, "filename: %s", filename)
		}
	}

	return nil
}

// CreateResourceWithFile create k8s resource with file
func CreateResourceWithFile(client kubernetes.Interface, filename string, option interface{}) error {
	var (
		data []byte
		err  error
	)
	if option != nil {
		data, err = template.ParseFile(filename, option)
	} else {
		data, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	reg := regexp.MustCompile(`(?m)^-{3,}$`)
	items := reg.Split(string(data), -1)
	for _, item := range items {
		objBytes := []byte(item)
		obj := new(object)
		err := yaml.Unmarshal(objBytes, obj)
		if err != nil {
			return err
		}
		if obj.Kind == "" {
			continue
		}
		f, ok := handlers[obj.Kind]
		if !ok {
			return errors.Errorf("unsupport kind %q", obj.Kind)
		}
		err = f(client, objBytes)
		if err != nil {
			return err
		}
	}

	return nil
}
