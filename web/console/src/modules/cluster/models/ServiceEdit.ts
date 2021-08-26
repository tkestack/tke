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

import { FetcherState, Identifiable, QueryState, RecordSet } from '@tencent/ff-redux';

import { Validation } from '../../common/models';
import { PortMap } from './PortMap';
import { Resource, ResourceFilter } from './ResourceOption';

export interface ServiceEdit extends Identifiable {
  /**
   * 服务名称
   */
  serviceName?: string;
  v_serviceName?: Validation;

  /**
   * 描述
   */
  description?: string;
  v_description?: Validation;

  /**
   * namespace
   */
  namespace?: string;
  v_namespace?: Validation;

  /** 访问方式 */
  communicationType?: string;

  /** 端口映射 */
  portsMap?: PortMap[];

  /** 是否开启headless */
  isOpenHeadless?: boolean;

  /** selector */
  selector?: Selector[];

  /** 是否展示引用workload标签的弹窗 */
  isShowWorkloadDialog?: boolean;

  /** 当前选择的workload的类型 */
  workloadType?: string;

  /** workload的查询 */
  workloadQuery?: QueryState<ResourceFilter>;

  /** workload的列表 */
  workloadList?: FetcherState<RecordSet<Resource>>;

  /** workload的选择 */
  workloadSelection?: Resource[];

  /** externalTrafficPolicy */
  externalTrafficPolicy?: string;

  /**会话保持 */
  sessionAffinity?: string;

  sessionAffinityTimeout?: number;

  v_sessionAffinityTimeout?: Validation;
}

export interface CLB extends Identifiable {
  /** loadBalancerId */
  loadBalancerId: string;

  /** loadBalancerName */
  loadBalancerName: string;

  /** 当前的lb类型，0位传统型，1为应用型 */
  forward: string;

  /** 网络类型 2为公网、3为内网 */
  loadBalancerType?: string;

  /** uniqVpcId */
  uniqVpcId?: string;

  [props: string]: any;
}

export interface Selector extends Identifiable {
  /** selector 的 key */
  key: string;
  v_key: Validation;

  /** selector 的 value */
  value: string;
  v_value: Validation;
}

/** 创建服务的时候，提交的jsonSchema */
export interface ServiceEditJSONYaml {
  /** 资源的类型 */
  kind?: string;

  /** api的版本 */
  apiVersion?: string;

  /** metadata */
  metadata?: ServiceMetadata;

  /** spec */
  spec?: ServiceSpec;
}

/** metadata的配置，非全部选项 */
interface ServiceMetadata {
  /** 插件能力 */
  annotations?: {
    [props: string]: string;
  };

  /** 集群名称 */
  clusterName?: string;

  /** labels */
  labels?: {
    [props: string]: string;
  };

  /** service的名称 */
  name?: string;

  /** service的命名空间 */
  namespace?: string;
}

/** spec的配置，非全部选项 */
interface ServiceSpec {
  /** clusterIP 用于设置 headless */
  clusterIP?: string;

  /** statefulset的headless 需要设置serviceName */
  serviceName?: string;

  /** 访问的类型 */
  type: string;

  /** 端口映射 */
  ports?: ServicePorts[];

  /** selector */
  selector?: {
    [props: string]: string;
  };

  /**会话保持 */
  sessionAffinity?: string;

  /** 会话保持的相关配置 */
  sessionAffinityConfig?: any;

  /** 会话保持的流量设置 */
  externalTrafficPolicy?: string;
}

export interface ServicePorts {
  /** 名称 */
  name: string;

  /** nodePort */
  nodePort?: number;

  /** 服务端口 */
  port?: number;

  /** 协议 */
  protocol?: string;

  /** 容器端口 */
  targetPort?: number;
}
