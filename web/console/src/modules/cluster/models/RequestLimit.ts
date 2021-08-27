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
import { Identifiable } from '@tencent/ff-redux';

export interface RequestLimit extends Identifiable {
  /**集群ID */
  clusterId?: string;

  /**Cpu Request总和*/
  totalCpuRequest?: number;

  /**内存 Request总和*/
  totalMemRequest?: number;

  /**总内存 */
  totalCpu?: number;

  /**总cpu */
  totalMem?: number;

  /**总gpu */
  totalGpu?: number;

  /**错误信息 */
  result?: any;
}

export interface RequestLimitFilter {
  regionId?: number;

  clusterIds?: string[];
}
