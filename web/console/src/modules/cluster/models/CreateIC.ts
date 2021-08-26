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

import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';
import { ContainerRuntimeEnum } from '../constants/Config';

export interface LabelsKeyValue {
  key?: string;
  value?: string;
}
export interface ICComponter {
  ipList?: string;
  ssh?: string;
  cidr?: string;
  role?: string;
  labels?: LabelsKeyValue[];
  authType?: string;
  username?: string;
  password?: string;
  privateKey?: string;
  passPhrase?: string;
  isEditing?: boolean;

  //添加节点时候复用了
  isGpu?: boolean;
}
export interface CreateIC extends Identifiable {
  /**集群名称 */
  name?: string;
  v_name?: Validation;

  k8sVersion?: string;

  networkDevice?: string;
  v_networkDevice?: Validation;

  cidr?: string;

  maxClusterServiceNum?: string;

  maxNodePodNum?: string;

  k8sVersionList?: any[];

  computerList?: ICComponter[];
  computerEdit?: ICComponter;

  vipAddress?: string;
  v_vipAddress?: Validation;

  vipPort?: string;
  v_vipPort?: Validation;

  vipType?: string;

  gpu?: boolean;

  gpuType?: string;

  merticsServer?: boolean;

  cilium?: string;

  networkMode?: string;

  asNumber?: string;
  v_asNumber: Validation;

  switchIp?: string;
  v_switchIp: Validation;

  containerRuntime: ContainerRuntimeEnum;
}
