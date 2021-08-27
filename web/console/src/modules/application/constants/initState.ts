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
import { uuid } from '@tencent/ff-redux';

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

export const initAppEditorState = {
  id: uuid(),
  metadata: {
    namespace: '',
    name: '',
    creationTimestamp: '',
    generation: 0
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
  },

  status: {
    observedGeneration: 0
  },

  v_editing: false
};

export const initResourceList = {
  id: uuid(),
  resources: []
};

export const initHistoryList = {
  id: uuid(),
  histories: []
};
