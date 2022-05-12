/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
/** 地域的相关操作 */
export const QueryRegion = 'QueryRegion';
export const FetchRegion = 'FetchRegion';
export const SelectRegion = 'SelectRegion';

/** 集群的相关操作 */
export const FetchClusterList = 'FetchClusterList';
export const QueryClusterList = 'QueryClusterList';
export const ClusterVersion = 'ClusterVersion';
export const SelectCluster = 'SelectCluster';

/** 日志采集列表、详情相关操作 */
export const FetchLogList = 'FetchLogList';
export const QueryLogList = 'QueryLogList';
export const SelectLog = 'SelectLog';
export const IsOpenLogStash = 'IsOpenLogStash';
export const IsDaemonsetNormal = 'IsDaemonsetNormal';
export const IsFetchDoneSpecificLog = 'IsFetchDoneSpecificLog';
export const FetchNamespaceList = 'FetchNamespaceList';
export const QueryNamespaceList = 'QueryNamespaceList';
export const NamespaceSelection = 'NamespaceSelection';
export const FetchLogDaemonset = 'FetchLogDaemonset';
export const QueryLogDaemonset = 'QueryLogDaemonset';
/** workloadActions相关操作 */
export const AuthorizeOpenLog = 'AuthorizeOpenLog';
export const ModifyLogStashFlow = 'ModifyLogStashFlow';
export const InlineDeleteLog = 'InlineDeleteLog';
/** 创建日志采集规则的相关操作 */
export const ClearLogStashEdit = 'ClearLogStashEdit';
export const LogStashName = 'LogStashName';
export const V_LogStashName = 'V_LogStashName';
export const V_SelectClusterSelection = 'V_SelectClusterSelection';
export const ChangeLogMode = 'ChangeLogMode';
export const IsSelectedAllNamespace = 'IsSelectedAllNamespace';
export const UpdateContainerLogs = 'UpdateContainerLogs';
export const NodeLogPath = 'NodeLogPath';
export const V_NodeLogPath = 'V_NodeLogPath';
export const NodeLogPathType = 'NodeLogPathType';
export const UpdateMetadata = 'UpdateMetadata';
export const ChangeConsumerMode = 'ChangeConsumerMode';
export const IsCanUseCkafka = 'IsCanUseCkafka';
export const IsCanUseCls = 'IsCanUseCls';
export const IsSelectedCkafka = 'IsSelectedCkafka';
export const FetchCkafkaList = 'FetchCkafkaList';
export const QueryCkafkaList = 'QueryCkafkaList';
export const SelectCkafka = 'SelectCkafka';
export const V_SelectCkafka = 'V_SelectCkafka';
export const FetchCTopicList = 'FetchCTopicList';
export const QueryCTopicList = 'QueryCTopicList';
export const SelectCTopic = 'SelectCTopic';
export const V_SelectCTopic = 'V_SelectCTopic';
export const AddressIP = 'AddressIP';
export const V_AddressIP = 'V_AddressIP';
export const AddressPort = 'AddressPort';
export const V_AddressPort = 'V_AddressPort';
export const Topic = 'Topic';
export const V_Topic = 'V_Topic';
export const FetchClsList = 'FetchClsList';
export const QueryClsList = 'QueryClsList';
export const SelectCls = 'SelectCls';
export const V_SelectCls = 'V_SelectCls';
export const FetchClsTopicList = 'FetchClsTopicList';
export const QueryClsTopicList = 'QueryClsTopicList';
export const SelectClsTopic = 'SelectClsTopic';
export const V_SelectClsTopic = 'V_SelectClsTopic';
export const EsAddress = 'EsAddress';
export const V_EsAddress = 'V_EsAddress';
export const IndexName = 'IndexName';
export const V_IndexName = 'V_IndexName';
export const EsUsername = 'EsUsername';
export const EsPassword = 'EsPassword';
export const FetchResourceList = 'FetchResourceList';
export const QueryResourceList = 'QueryResourceList';
export const UpdateResourceTarget = 'UpdateResourceTarget';
/**mode为容器文件路径 相关操作 */
export const SelectContainerFileNamespace = 'SelectContainerFileNamespace';
export const V_ContainerFileNamespace = 'V_ContainerFileNamespace';

export const SelectContainerFileWorkloadType = 'SelectContainerFileWorkloadType';
export const V_ContainerFileWorkloadType = 'V_ContainerFileWorkloadType';

export const SelectContainerFileWorkload = 'SelectContainerFileWorkload';
export const V_ContainerFileWorkload = 'V_ContainerFileWorkload';

export const FetchPodList = 'FetchPodList';
export const QueryPodList = 'QueryPodList';

export const UpdateContainerFilePaths = 'UpdateContainerFilePaths';
export const UpdateContaierFileWorkloadList = 'UpdateContaierFileWorkloadList ';
export const isFirstFetchResource = 'isFirstFetchResource';

//业务侧
// export const QueryNamespaceList = 'QueryNamespaceList';
export const SelectNamespace = 'SelectNamespace';
export const InitProjectList = 'InitProjectList';
export const ProjectSelection = 'ProjectSelection';
// export const FetchNamespaceList = 'FetchNamespaceList';
export const FetchProjectList = 'FetchProjectList';
export const QueryProject = 'QueryProject';
