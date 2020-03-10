import { initValidator } from '../../common/models';
import { uuid } from '@tencent/ff-redux';

export const initRoleCreationState = {
  id: uuid(),
  spec: {
    /** 展示名 */
    displayName: '',
    /** 描述 */
    description: '',
    /** 策略 */
    policies: []
  },
  status: {
    /** 用户 */
    users: [],
    /** 用户组 */
    groups: []
  },
};

export const initRoleEditorState = {
  id: uuid(),
  metadata: {
    /** 名称 */
    name: '',
    /** 创建时间 */
    creationTimestamp: ''
  },
  spec: {
    /** 展示名 */
    displayName: '',
    /** 描述 */
    description: '',
  },

  /** 是否正在编辑 */
  v_editing: false,
};

export const initRoleAssociationState = {
  /** 最新数据 */
  roles: [],
  /** 原来数据 */
  originRoles: [],
  /** 新增数据 */
  addRoles: [],
  /** 删除数据 */
  removeRoles: []
};

export const initRoleFilterState = {
  /** 目标资源，如localgroup,policy,localidentity */
  resource: '',
  /** 资源id */
  resourceID: '',
  /** 关联/解关联操作后的回调函数 */
  callback: undefined
};

export const initGroupCreationState = {
  id: uuid(),
  spec: {
    /** 展示名 */
    displayName: '',
    /** 描述 */
    description: '',
  },
  status: {
    /** 用户 */
    users: []
  },
};

export const initGroupEditorState = {
  id: uuid(),
  metadata: {
    /** 名称 */
    name: '',
    /** 创建时间 */
    creationTimestamp: ''
  },
  spec: {
    /** 展示名 */
    displayName: '',
    /** 描述 */
    description: '',
  },

  /** 是否正在编辑 */
  v_editing: false,
};

export const initGroupAssociationState = {
  /** 最新数据 */
  groups: [],
  /** 原来数据 */
  originGroups: [],
  /** 新增数据 */
  addGroups: [],
  /** 删除数据 */
  removeGroups: []
};

export const initGroupFilterState = {
  /** 目标资源，如role,policy,localidentity */
  resource: '',
  /** 资源id */
  resourceID: '',
  /** 关联/解关联操作后的回调函数 */
  callback: undefined
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

export const initCommonUserFilterState = {
  /** 目标资源，如localgroup/role/policy */
  resource: '',
  /** 资源id */
  resourceID: '',
  /** 关联/解关联操作后的回调函数 */
  callback: undefined
};

export const initPolicyEditorState = {
  id: uuid(),
  metadata: {
    /** 名称 */
    name: '',
    /** 创建时间 */
    creationTimestamp: ''
  },
  spec: {
    /** 展示名 */
    displayName: '',
    /** 描述 */
    description: '',
  },

  /** 是否正在编辑 */
  v_editing: false,
};

export const initPolicyAssociationState = {
  /** 最新数据 */
  policies: [],
  /** 原来数据 */
  originPolicies: [],
  /** 新增数据 */
  addPolicies: [],
  /** 删除数据 */
  removePolicies: []
};

export const initPolicyFilterState = {
  /** 目标资源，如role,localidentity,localgroup */
  resource: '',
  /** 资源id */
  resourceID: '',
  /** 关联/解关联操作后的回调函数 */
  callback: undefined
};
