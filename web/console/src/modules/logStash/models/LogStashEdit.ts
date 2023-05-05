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
import { FetcherState, Identifiable, QueryState, RecordSet } from '@tencent/ff-redux';

import { Validation } from '../../common/models';
import { ContainerLogs, MetadataItem } from './';
import { ContainerFilePathItem } from './ContainerFilePathItem';
import { Pod, PodListFilter } from './pod';
import { Resource, ResourceFilter, ResourceTarget, WorkLoadList } from './Resource';

export interface LogStashEdit extends Identifiable {
  /** 采集规则的名称 */
  logStashName?: string;
  v_logStashName?: Validation;

  /** 校验选择集群的ID */
  v_clusterSelection?: Validation;

  /** 当前日志采集的类型 */
  logMode?: string;

  /**
   * pre: 类型为指定容器日志
   * 是否选择所有容器
   */
  isSelectedAllNamespace?: string;

  /** 指定容器的编辑项数组 */
  containerLogs?: ContainerLogs[];

  /**
   * pre: 类型为指定主机路径
   * 收集路径
   */
  nodeLogPath?: string;
  v_nodeLogPath?: Validation;

  nodeLogPathType: 'host' | 'container';

  /**
   * pre: 类型为指定主机路径
   * metadata标签
   */
  metadatas: MetadataItem[];

  /** 当前的消费端的类型 */
  consumerMode?: string;

  /** 访问地址 ip */
  addressIP?: string;
  v_addressIP?: Validation;

  /** 访问地址 port */
  addressPort?: string;
  v_addressPort?: Validation;

  /** topic 的输入*/
  topic?: string;
  v_topic?: Validation;

  /** es的地址 */
  esAddress?: string;
  v_esAddress?: Validation;

  /** es username */
  esUsername?: string;

  /** es password */
  esPassword?: string;

  /** 索引名称 */
  indexName?: string;
  v_indexName?: Validation;

  /** 资源列表 resourceList */
  resourceList?: FetcherState<RecordSet<Resource>>;

  /** 资源列表的查询 */
  resourceQuery?: QueryState<ResourceFilter>;

  /** 资源的使用對象 */
  resourceTarget?: ResourceTarget;
  /** 是否是第一次獲取資源對象 */
  isFirstFetchResource?: boolean;
  /** 容器文件路径 命名空间*/
  containerFileNamespace?: string;
  v_containerFileNamespace?: Validation;

  /** 容器文件路径 工作负载类型*/
  containerFileWorkloadType?: string;
  v_containerFileWorkloadType?: Validation;

  /** 容器文件路径 工作负载*/
  containerFileWorkload?: string;
  v_containerFileWorkload?: Validation;

  /**容器文件路径 容器名+文件路径*/
  containerFilePaths?: ContainerFilePathItem[];
  /**容器工作負載選項 */
  containerFileWorkloadList?: WorkLoadList[];

  podList?: FetcherState<RecordSet<Pod>>;
  podListQuery?: QueryState<PodListFilter>;
}

export interface LogStashEditOperator {
  /** 地域 */
  regionId: number;

  /** 当前的编辑类型 */
  mode: string;

  /** 当前的集群名称 */
  clusterId: string;

  /**当前集权的名称 */
  clusterVersion: string;
}

export interface LogStashEditYaml {
  /**资源类型 */
  kind: string;
  /**API的版本 */
  apiVersion: string;
  /**metadata */
  metadata: LogStashMetadata;
  /** */
  spec: LogStashSpec;
}

export interface LogStashMetadata {
  /**日志收集规则的名字 */
  name?: string;
  /**日志收集规则的命名空间 固定为kube-system */
  namespace?: string;
  /**固定参数 只有设置了pod-log ,才需要将这个label设置为'true'*/
  labels?: {
    'log.tke.cloud.tencent.com/pod-log': string;
  };

  resourceVersion?: string;
}

export interface LogStashSpec {
  /**相关描述 */
  description?: string;
  /**输入端参数 */
  input?: PodLogInput | HostLogInput | ContainerLogInput;
  /**输出端参数 */
  output?: ClsOutput | CkafkaOutput | KafkaOutpot | ElasticsearchOutput;
}

export interface PodLogInput {
  pod_log_input: {
    container_log_files: {
      [props: string]: {
        path: string[];
      };
    };
    metadata: boolean;
    workload: {
      name: string;
      type: string;
    };
  };
  type: 'pod-log';
}

export interface HostLogInput {
  host_log_input: {
    labels: {
      [props: string]: string;
    };
    path: string;
  };
  type: 'host-log';
}

export interface ContainerLogInput {
  container_log_input: { all_namespaces?: boolean; namespaces?: ContainerLogNamespace[] };
  type: 'container-log';
}

export interface ContainerLogNamespace {
  all_containers: boolean;
  namespace: string;
  services?: string[];
  workloads?: {
    name: string;
    type: string;
  }[];
}

export interface ClsOutput {
  cls_output: {
    host?: string;
    logset_id: string;
    port?: number;
    topic_id: string;
  };
  type: 'cls';
}

export interface CkafkaOutput {
  ckafka_output: {
    host: string;
    instance_id: string;
    port: number;
    topic: string;
    topic_id: string;
  };
  type: 'ckafka';
}

export interface KafkaOutpot {
  kafka_output: {
    host: string;
    port: number;
    topic: string;
  };
  type: 'kafka';
}

export interface ElasticsearchOutput {
  elasticsearch_output: {
    hosts: string[];
    index: string;
    user: string;
    password: string;
  };
  type: 'elasticsearch';
}
