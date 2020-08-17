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
