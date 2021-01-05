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

package application

import (
	"context"
	"reflect"
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/application"
	"tkestack.io/tke/api/client/clientset/internalversion/fake"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
)

func TestValidateApplication_Normal(t *testing.T) {
	cluster1 := "cls-test1"
	clientset := fake.NewSimpleClientset()

	type args struct {
		ctx               context.Context
		app               *application.App
		applicationClient applicationinternalclient.ApplicationInterface
	}
	tests := []struct {
		name string
		args args
		want field.ErrorList
	}{
		{
			name: "normal app",
			args: args{
				ctx: context.TODO(),
				app: &application.App{
					TypeMeta: v1.TypeMeta{
						Kind:       "App",
						APIVersion: "application.tkestack.io/v1",
					},
					ObjectMeta: v1.ObjectMeta{
						Name:      "app3",
						Namespace: "default",
					},
					Spec: application.AppSpec{
						Type:          "HelmV3",
						TenantID:      "10001",
						Name:          "p2p",
						TargetCluster: cluster1,
						Chart: application.Chart{
							TenantID:       "10001",
							ChartGroupName: "local",
							ChartName:      "p2p",
							ChartVersion:   "1.0.0",
							RepoURL:        "http://chartmuseum:8080",
							ImportedRepo:   true,
						},
						Values: application.AppValues{
							RawValuesType: "yaml",
							RawValues:     "",
							Values:        []string{},
						},
						DryRun: false,
					},
				},
				applicationClient: clientset.Application(),
			},
			want: field.ErrorList{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateApplication(tt.args.ctx, tt.args.app, tt.args.applicationClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateApplication_DuplicateAppInOneCluster(t *testing.T) {
	cluster1 := "cls-test1"
	// see https://github.com/kubernetes/code-generator/blob/release-1.18/cmd/client-gen/generators/fake/generator_fake_for_type.go
	// k8s code generator, not support field selector.
	app := &application.App{
		TypeMeta: v1.TypeMeta{
			Kind:       "App",
			APIVersion: "application.tkestack.io/v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: application.AppSpec{
			Type:          "HelmV3",
			TenantID:      "10001",
			Name:          "p2p",
			TargetCluster: cluster1,
			Chart: application.Chart{
				TenantID:       "10001",
				ChartGroupName: "local",
				ChartName:      "p2p",
				ChartVersion:   "1.0.0",
				RepoURL:        "http://chartmuseum:8080",
				ImportedRepo:   true,
			},
			Values: application.AppValues{
				RawValuesType: "yaml",
				RawValues:     "",
				Values:        []string{},
			},
			DryRun: false,
		},
	}
	clientset := fake.NewSimpleClientset(app)

	type args struct {
		ctx               context.Context
		app               *application.App
		applicationClient applicationinternalclient.ApplicationInterface
	}
	tests := []struct {
		name string
		args args
		want field.ErrorList
	}{
		{
			name: "normal app",
			args: args{
				ctx: context.TODO(),
				app: &application.App{
					TypeMeta: v1.TypeMeta{
						Kind:       "App",
						APIVersion: "application.tkestack.io/v1",
					},
					ObjectMeta: v1.ObjectMeta{
						Name:      "app2",
						Namespace: "default",
					},
					Spec: application.AppSpec{
						Type:          "HelmV3",
						TenantID:      "10001",
						Name:          "p2p",
						TargetCluster: cluster1,
						Chart: application.Chart{
							TenantID:       "10001",
							ChartGroupName: "local",
							ChartName:      "p2p",
							ChartVersion:   "1.0.0",
							RepoURL:        "http://chartmuseum:8080",
							ImportedRepo:   true,
						},
						Values: application.AppValues{
							RawValuesType: "yaml",
							RawValues:     "",
							Values:        []string{},
						},
						DryRun: false,
					},
				},
				applicationClient: clientset.Application(),
			},
			want: field.ErrorList{
				&field.Error{
					Type:     "FieldValueDuplicate",
					Field:    "spec.name",
					BadValue: "p2p",
					Detail:   "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateApplication(tt.args.ctx, tt.args.app, tt.args.applicationClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}
