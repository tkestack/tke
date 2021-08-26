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

export interface VolumeItem extends Identifiable {
  /** 卷类型 */
  volumeType: string;
  v_volumeType: Validation;

  /** 卷名称 */
  name: string;
  v_name: Validation;

  /** 源路径 */
  hostPathType: string;
  hostPath: string;
  v_hostPath: Validation;

  /** nfs路径 */
  nfsPath: string;
  v_nfsPath: Validation;

  /** 配置项 */
  configKey: ConfigItems[];
  configName: string;

  /** secret的相关配置 */
  secretKey: ConfigItems[];
  secretName: string;

  /** pvc的选择 */
  pvcSelection: string;
  v_pvcSelection: Validation;

  /** 新创建的pvc */
  newPvcName: string;
  pvcEditInfo: PvcEditInfo;

  /** 当前数据卷是否被挂载 */
  isMounted: boolean;
}

export interface PvcEditInfo {
  /** accessMode */
  accessMode: string;

  /** storageClassName */
  storageClassName: string;

  /** storage */
  storage: string;
}

export interface ConfigItems extends Identifiable {
  /** config的Key */
  configKey?: string;
  v_configKey?: Validation;

  /** 配置的子路径 */
  path?: string;
  v_path?: Validation;

  /** 当前mode */
  mode?: string;
  v_mode?: Validation;
}
