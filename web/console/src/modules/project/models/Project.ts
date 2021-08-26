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
import { Manager } from './Manager';

export interface Project extends Identifiable {
  /** metadata */
  metadata: ProjectMetadata;

  /** spec */
  spec: ProjectSpec;

  /** status */
  status: ProjectStatus;
}

interface ProjectMetadata {
  /** projectId */
  name?: string;

  /** 创建时间 */
  creationTimestamp?: string;

  /** 其余属性 */
  [props: string]: any;
}

interface ProjectSpec {
  /** project名称 */
  displayName: string;

  /** Project成员 */
  members: string[];

  /** project分配的quota */
  clusters: {
    [props: string]: {
      hard?: StatusResource;
    };
  };

  parentProjectName?: string;
}

interface ProjectStatus {
  /** 可分配 */
  used?: StatusResource;
  /** 未分配 */
  calculatedChildProjects?: string;

  calculatedNamespaces?: string[];

  /** 状态 */
  phase?: string;
}

interface StatusResource {
  [props: string]: string;
}

export interface ProjectEdition extends Identifiable {
  id: string;

  resourceVersion: string;

  members: Manager[];

  displayName: string;
  v_displayName: Validation;

  clusters: {
    name: string;
    v_name: Validation;
    resourceLimits: ProjectResourceLimit[];
  }[];

  parentProject: string;

  status: any;
}

export interface ProjectResourceLimit extends Identifiable {
  type: string;
  v_type: Validation;
  value: string;
  v_value: Validation;
}

export interface ProjectFilter {
  /** 业务id */
  ProjectId?: string;

  /**业务名称 */
  displayName?: string;

  parentProject?: string;
}

export interface ProjectUserMap {
  [props: string]: {
    id: string;
    username: string;
  }[];
}

export interface UserManagedProject extends Identifiable {
  name: string;
}
export interface UserManagedProjectFilter {
  userId: string;
}
