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

export interface ChartGroup extends Identifiable {
  metadata?: {
    name: string;
    creationTimestamp?: string;
  };
  spec: {
    name: string;
    tenantID?: string;
    displayName?: string;
    visibility: string;
    description?: string;
    type: string;
    projects?: string[];
    users?: string[];
    importedInfo?: {
      addr: string;
      username?: string;
      password?: string;
    };
  };
  status?: {
    chartCount?: number;
    phase: string;
    [props: string]: any;
  };
}

export interface ChartGroupFilter {
  repoType?: string;
}

export interface ChartGroupDetailFilter {
  name: string;
  projectID: string;
}

export interface ChartGroupCreation extends Identifiable {
  spec: {
    name: string;
    displayName?: string;
    visibility: string;
    description?: string;
    type: string;
    projects?: string[];
    users?: string[];
    importedInfo?: {
      addr: string;
      username?: string;
      password?: string;
    };
  };
}

export interface ChartGroupEditor extends Identifiable {
  metadata?: {
    name: string;
    creationTimestamp?: string;
  };
  spec: {
    name: string;
    displayName?: string;
    visibility: string;
    description?: string;
    type: string;
    projects?: string[];
    users?: string[];
    importedInfo?: {
      addr: string;
      username?: string;
      password?: string;
    };
  };

  /** 是否正在编辑 */
  v_editing?: boolean;
}
