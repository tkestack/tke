/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-generated-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE. DO NOT EDIT.
var map_ConfigMap = map[string]string{
	"":           "ConfigMap holds configuration data for tke to consume.",
	"data":       "Data contains the configuration data. Each key must consist of alphanumeric characters, '-', '_' or '.'. Values with non-UTF-8 byte sequences must use the BinaryData field. The keys stored in Data must not overlap with the keys in the BinaryData field, this is enforced during validation process.",
	"binaryData": "BinaryData contains the binary data. Each key must consist of alphanumeric characters, '-', '_' or '.'. BinaryData can contain byte sequences that are not in the UTF-8 range. The keys stored in BinaryData must not overlap with the ones in the Data field, this is enforced during validation process.",
}

func (ConfigMap) SwaggerDoc() map[string]string {
	return map_ConfigMap
}

var map_ConfigMapList = map[string]string{
	"":      "ConfigMapList is a resource containing a list of ConfigMap objects.",
	"items": "Items is the list of ConfigMaps.",
}

func (ConfigMapList) SwaggerDoc() map[string]string {
	return map_ConfigMapList
}

var map_DataBase = map[string]string{
	"": "Database describes the attributes of a MeshManager.",
}

func (DataBase) SwaggerDoc() map[string]string {
	return map_DataBase
}

var map_MeshManager = map[string]string{
	"":     "MeshManager is a manager to manager mesh clusters.",
	"spec": "Spec defines the desired identities of MeshManager.",
}

func (MeshManager) SwaggerDoc() map[string]string {
	return map_MeshManager
}

var map_MeshManagerList = map[string]string{
	"":      "MeshManagerList is the whole list of all meshmanagers which owned by a tenant.",
	"items": "List of volume decorators.",
}

func (MeshManagerList) SwaggerDoc() map[string]string {
	return map_MeshManagerList
}

var map_MeshManagerSpec = map[string]string{
	"": "MeshManagerSpec describes the attributes of a MeshManager.",
}

func (MeshManagerSpec) SwaggerDoc() map[string]string {
	return map_MeshManagerSpec
}

var map_MeshManagerStatus = map[string]string{
	"":                            "MeshManagerStatus is information about the current status of a MeshManager.",
	"phase":                       "Phase is the current lifecycle phase of the MeshManager of cluster.",
	"reason":                      "Reason is a brief CamelCase string that describes any failure.",
	"retryCount":                  "RetryCount is a int between 0 and 5 that describes the time of retrying initializing.",
	"lastReInitializingTimestamp": "LastReInitializingTimestamp is a timestamp that describes the last time of retrying initializing.",
}

func (MeshManagerStatus) SwaggerDoc() map[string]string {
	return map_MeshManagerStatus
}

var map_StorageBackend = map[string]string{
	"": "StorageBackend describes the attributes of a backend storage StorageType can be \"influxdb\",\"elasticsearch\",\"es\",\"thanos\"",
}

func (StorageBackend) SwaggerDoc() map[string]string {
	return map_StorageBackend
}

// AUTO-GENERATED FUNCTIONS END HERE
