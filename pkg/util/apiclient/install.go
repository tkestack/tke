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
	"context"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	admissionv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	kuberuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	kubeaggregator "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	"tkestack.io/tke/pkg/util/template"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	aaclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

type object struct {
	Kind string `yaml:"kind"`
}

var (
	handlers   map[string]func(context.Context, kubernetes.Interface, []byte) error
	kaHandlers map[string]func(context.Context, kubeaggregator.Interface, []byte) error
	aaHandlers map[string]func(context.Context, aaclientset.Interface, []byte) error
)

func init() {
	handlers = make(map[string]func(context.Context, kubernetes.Interface, []byte) error)
	kaHandlers = make(map[string]func(context.Context, kubeaggregator.Interface, []byte) error)
	aaHandlers = make(map[string]func(context.Context, aaclientset.Interface, []byte) error)

	// apiregistration
	kaHandlers["APIService"] = func(ctx context.Context, client kubeaggregator.Interface, data []byte) error {
		obj := new(apiregistrationv1.APIService)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateAPIService(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}

	//add CustomResourceDefinition
	aaHandlers["CustomResourceDefinition"] = func(ctx context.Context, client aaclientset.Interface, data []byte) error {
		obj := new(apiextensionsv1.CustomResourceDefinition)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateCustomResourceDefinition(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}

	// core
	handlers["ConfigMap"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(corev1.ConfigMap)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateConfigMap(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Endpoints"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(corev1.Endpoints)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateEndpoints(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Namespace"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(corev1.Namespace)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateNamespace(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Secret"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(corev1.Secret)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateSecret(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Service"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(corev1.Service)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateService(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["ServiceAccount"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(corev1.ServiceAccount)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateServiceAccount(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	// batch
	handlers["Job"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(batchv1.Job)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateJob(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	// batchv1beta1
	handlers["CronJob"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(batchv1beta1.CronJob)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateCronJob(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	// apps
	handlers["DaemonSet"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(appsv1.DaemonSet)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateDaemonSet(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Pod"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(corev1.Pod)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdatePod(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["Deployment"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(appsv1.Deployment)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateDeployment(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["StatefulSet"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(appsv1.StatefulSet)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateStatefulSet(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}

	// extentions
	handlers["Ingress"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(extensionsv1beta1.Ingress)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateIngress(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}

	// rbac
	handlers["Role"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(rbacv1.Role)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateRole(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["RoleBinding"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(rbacv1.RoleBinding)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateRoleBinding(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["ClusterRole"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(rbacv1.ClusterRole)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateClusterRole(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["ClusterRoleBinding"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(rbacv1.ClusterRoleBinding)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateClusterRoleBinding(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	// admissionregistration
	handlers["ValidatingWebhookConfiguration"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(admissionv1beta1.ValidatingWebhookConfiguration)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateValidatingWebhookConfiguration(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
	handlers["MutatingWebhookConfiguration"] = func(ctx context.Context, client kubernetes.Interface, data []byte) error {
		obj := new(admissionv1beta1.MutatingWebhookConfiguration)
		if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), data, obj); err != nil {
			return errors.Wrapf(err, "unable to decode %s", reflect.TypeOf(obj).String())
		}
		err := CreateOrUpdateMutatingWebhookConfiguration(ctx, client, obj)
		if err != nil {
			return err
		}
		return nil
	}
}

// CreateResourceWithDir create k8s resource with dir
func CreateResourceWithDir(ctx context.Context, client kubernetes.Interface, pattern string, option interface{}) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return errors.New("no matches found")
	}
	for _, filename := range matches {
		err := CreateResourceWithFile(ctx, client, filename, option)
		if err != nil {
			return errors.Wrapf(err, "filename: %s", filename)
		}
	}

	return nil
}

func CreateResourceWiteContent(ctx context.Context, client kubernetes.Interface, content string,
	option interface{}) error {
	var (
		data []byte
		err  error
	)
	if option != nil {
		data, err = template.ParseString(content, option)
	} else {
		data = []byte(content)
	}

	if err != nil {
		return err
	}

	items := strings.Split(string(data), "\n---")
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
		err = f(ctx, client, objBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateResourceWithFile create k8s resource with file
func CreateResourceWithFile(ctx context.Context, client kubernetes.Interface, filename string, option interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return CreateResourceWiteContent(ctx, client, string(data), option)
}

// CreateKAResourceWithFile create k8s and kube-aggregator resource with file
func CreateKAResourceWithFile(ctx context.Context, client kubernetes.Interface, kaClient kubeaggregator.Interface, filename string, option interface{}) error {
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

	items := strings.Split(string(data), "\n---")
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
		if obj.Kind == "APIService" {
			f, ok := kaHandlers[obj.Kind]
			if !ok {
				return errors.Errorf("unsupport kind %q", obj.Kind)
			}
			err = f(ctx, kaClient, objBytes)
			if err != nil {
				return err
			}
		} else {
			f, ok := handlers[obj.Kind]
			if !ok {
				return errors.Errorf("unsupport kind %q", obj.Kind)
			}
			err = f(ctx, client, objBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}



// CreateAsResourceWithFile create k8s and apiextensions-apiserver resource with file
func CreateAsResourceWithFile(ctx context.Context, client kubernetes.Interface, aaClient aaclientset.Interface, filename string, option interface{}) error {
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

	items := strings.Split(string(data), "\n---")
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
		if obj.Kind == "CustomResourceDefinition" {
			f, ok := aaHandlers[obj.Kind]
			if !ok {
				return errors.Errorf("unsupport kind %q", obj.Kind)
			}
			err = f(ctx, aaClient, objBytes)
			if err != nil {
				return err
			}
		} else {
			f, ok := handlers[obj.Kind]
			if !ok {
				return errors.Errorf("unsupport kind %q", obj.Kind)
			}
			err = f(ctx, client, objBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
