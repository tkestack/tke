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
