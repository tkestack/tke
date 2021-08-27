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
import { PeEdit } from './PeEdit';

export interface AddonEdit {
  /** 当前选择的组件 */
  addonName?: string;

  /** 事件持久化的编辑项 */
  peEdit?: PeEdit;
}

interface AddonEditBasicJsonYaml {
  /** 资源的类型 */
  kind: string;

  /** api的版本 */
  apiVersion: string;

  /** metadata */
  metadata?: any;

  /** spec */
  spec?: any;
}

/** ====================== Helm、GameApp 创建相关的yaml ====================== */
export interface AddonEditUniversalJsonYaml extends AddonEditBasicJsonYaml {
  metadata: {
    generateName: string;
  };

  spec: {
    clusterName: string;
  };
}
/** ====================== Helm、GameApp 创建相关的yaml ====================== */

/** ====================== persistentEvent创建相关的yaml ===================== */
export interface AddonEditPeJsonYaml extends AddonEditBasicJsonYaml {
  metadata: {
    generateName: string;
  };

  spec: {
    clusterName: string;
    persistentBackEnd: PersistentBackEnd;
  };
}

export interface PersistentBackEnd {
  /** es的配置 */
  es: EsInfo;
}

export interface EsInfo {
  ip: string;
  port: number;
  scheme: string;
  indexName: string;
  user: string;
  password: string;
}
/** ====================== persistentEvent创建相关的yaml ===================== */
