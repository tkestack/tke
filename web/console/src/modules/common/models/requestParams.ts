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

/** 发起参数请求的数据格式 */
export interface RequestParams {
  /** 获取数据的方法，用于独立部署版 */
  method?: string;

  /** 发起请求的url */
  url?: string;

  /** userDefinedHeader，用户自定义头部 */
  userDefinedHeader?: UserDefinedHeader;

  /** data ，post的数据的时候，传的body */
  data?: any;

  /** 公有云 TCE 云API的参数请求，请求接口所需的参数 */
  apiParams?: ApiParams;

  /** 请求的基本域名url */
  baseURL?: string;
}

/** 公有云 TCE 云API的参数请求 */
interface ApiParams {
  /** 地域ID */
  regionId: number;

  /** 接口名称，默认为 ForwardRequest */
  interfaceName?: string;

  /** 请求的模块 */
  module: string;

  /** 其余的参数，主要是云API的参数列表 */
  restParams?: any;

  /** 云api的 options配置参数 */
  opts?: {
    // 是否展示tipErr，顶部的提示框
    tipErr?: boolean;

    global?: boolean;
  };
}

/** 用户自定义头部 */
export interface UserDefinedHeader {
  /** 请求接收的格式 */
  Accept?: string;

  /** 请求发送的格式 */
  'Content-Type'?: string;

  /** 集群的名称 */
  'X-TKE-ClusterName'?: string;

  'X-TKE-ProjectName'?: string;
}
