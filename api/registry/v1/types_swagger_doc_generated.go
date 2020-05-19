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
var map_Chart = map[string]string{
	"":     "Chart is a chart in chart group of chartmuseum registry.",
	"spec": "Spec defines the desired identities of chart in this set.",
}

func (Chart) SwaggerDoc() map[string]string {
	return map_Chart
}

var map_ChartGroup = map[string]string{
	"":     "ChartGroup is a chart container in chartmuseum registry.",
	"spec": "Spec defines the desired identities of chart group in this set.",
}

func (ChartGroup) SwaggerDoc() map[string]string {
	return map_ChartGroup
}

var map_ChartGroupList = map[string]string{
	"":      "ChartGroupList is the whole list of all chart groups which owned by a tenant.",
	"items": "List of chart groups",
}

func (ChartGroupList) SwaggerDoc() map[string]string {
	return map_ChartGroupList
}

var map_ChartGroupSpec = map[string]string{
	"": "ChartGroupSpec is a description of a chart group.",
}

func (ChartGroupSpec) SwaggerDoc() map[string]string {
	return map_ChartGroupSpec
}

var map_ChartGroupStatus = map[string]string{
	"":                   "ChartGroupStatus represents information about the status of a chart group.",
	"lastTransitionTime": "The last time the condition transitioned from one status to another.",
	"reason":             "The reason for the condition's last transition.",
	"message":            "A human readable message indicating details about the transition.",
}

func (ChartGroupStatus) SwaggerDoc() map[string]string {
	return map_ChartGroupStatus
}

var map_ChartInfo = map[string]string{
	"":     "ChartInfo describes detail of a chart version.",
	"spec": "Spec defines the desired identities of a chart.",
}

func (ChartInfo) SwaggerDoc() map[string]string {
	return map_ChartInfo
}

var map_ChartInfoSpec = map[string]string{
	"": "ChartInfoSpec is a description of a ChartInfo.",
}

func (ChartInfoSpec) SwaggerDoc() map[string]string {
	return map_ChartInfoSpec
}

var map_ChartList = map[string]string{
	"":      "ChartList is the whole list of all charts which owned by a chart group.",
	"items": "List of charts",
}

func (ChartList) SwaggerDoc() map[string]string {
	return map_ChartList
}

var map_ChartProxyOptions = map[string]string{
	"": "ChartProxyOptions is the query options to a ChartInfo proxy call.",
}

func (ChartProxyOptions) SwaggerDoc() map[string]string {
	return map_ChartProxyOptions
}

var map_ChartStatus = map[string]string{
	"lastTransitionTime": "The last time the condition transitioned from one status to another.",
	"reason":             "The reason for the condition's last transition.",
	"message":            "A human readable message indicating details about the transition.",
}

func (ChartStatus) SwaggerDoc() map[string]string {
	return map_ChartStatus
}

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

var map_Namespace = map[string]string{
	"":     "Namespace is an image container in registry.",
	"spec": "Spec defines the desired identities of namespace in this set.",
}

func (Namespace) SwaggerDoc() map[string]string {
	return map_Namespace
}

var map_NamespaceList = map[string]string{
	"":      "NamespaceList is the whole list of all namespaces which owned by a tenant.",
	"items": "List of namespaces",
}

func (NamespaceList) SwaggerDoc() map[string]string {
	return map_NamespaceList
}

var map_NamespaceSpec = map[string]string{
	"": "NamespaceSpec is a description of a namespace.",
}

func (NamespaceSpec) SwaggerDoc() map[string]string {
	return map_NamespaceSpec
}

var map_NamespaceStatus = map[string]string{
	"": "NamespaceStatus represents information about the status of a namespace.",
}

func (NamespaceStatus) SwaggerDoc() map[string]string {
	return map_NamespaceStatus
}

var map_Repository = map[string]string{
	"":     "Repository is a repo in namespace of registry.",
	"spec": "Spec defines the desired identities of repository in this set.",
}

func (Repository) SwaggerDoc() map[string]string {
	return map_Repository
}

var map_RepositoryList = map[string]string{
	"":      "RepositoryList is the whole list of all repositories which owned by a namespace.",
	"items": "List of repositories",
}

func (RepositoryList) SwaggerDoc() map[string]string {
	return map_RepositoryList
}

// AUTO-GENERATED FUNCTIONS END HERE
