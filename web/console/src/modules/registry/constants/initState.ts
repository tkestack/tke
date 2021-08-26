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

import { uuid } from '@tencent/ff-redux';

export const initChartGroupCreationState = {
  id: uuid(),
  spec: {
    name: '',
    displayName: '',
    visibility: 'Public',
    description: '',
    type: 'SelfBuilt',
    projects: [],
    users: [],
    importedInfo: {
      addr: '',
      username: '',
      password: ''
    }
  }
};

export const initChartGroupEditorState = {
  id: uuid(),
  metadata: {
    name: '',
    creationTimestamp: ''
  },
  spec: {
    name: '',
    displayName: '',
    visibility: '',
    description: '',
    type: '',
    projects: [],
    users: [],
    importedInfo: {
      addr: '',
      username: '',
      password: ''
    }
  },

  v_editing: false
};

export const initUserInfoState = {
  name: '',
  uid: '',
  groups: [''],
  extra: {
    displayname: '',
    tenantid: ''
  }
};

export const initChartEditorState = {
  id: uuid(),
  metadata: {
    namespace: '',
    name: '',
    creationTimestamp: ''
  },
  spec: {
    chartGroupName: '',
    displayName: '',
    name: '',
    tenantID: '',
    visibility: ''
  },
  status: {
    pullCount: 0,
    versions: []
  },

  v_editing: false,
  sortedVersions: [],
  selectedVersion: {}
};

export const initRemovedChartVersionsState = {
  versions: []
};

export const initAppCreationState = {
  id: uuid(),
  metadata: {
    namespace: ''
  },
  spec: {
    chart: {
      chartGroupName: '',
      chartName: '',
      chartVersion: '',
      tenantID: ''
    },
    name: '',
    targetCluster: '',
    tenantID: '',
    type: 'HelmV3',
    values: {
      rawValues: '',
      rawValuesType: 'yaml',
      values: ['']
    }
  }
};

export const initCommonUserAssociationState = {
  /** 最新数据 */
  users: [],
  /** 原始数据 */
  originUsers: [],
  /** 新增数据 */
  addUsers: [],
  /** 删除数据 */
  removeUsers: []
};
