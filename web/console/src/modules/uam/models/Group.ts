import { Identifiable } from '@tencent/ff-redux';
import { Validation } from '@tencent/ff-validator';

export interface Group extends Identifiable {
  metadata?: {
    /** 名称 */
    name: string;
    /** 创建时间 */
    creationTimestamp?: string;
  };
  spec: {
    /** 展示名 */
    displayName: string;
    /** 描述 */
    description: string;
  };
  [props: string]: any;
}

export interface GroupPlain extends Identifiable {
  /** id */
  name?: string;
  /** 名称 */
  displayName?: string;
  /** 描述 */
  description: string;
}

/** 用于列表查询 */
export interface GroupFilter {
  /** 目标资源，如role,policy,localidentity */
  resource: string;
  /** 资源id */
  resourceID: string;
  /** 关联/解关联操作后的回调函数 */
  callback?: () => void;
}

/** 用于单个查询 */
export interface GroupInfoFilter {
  name: string;
}

export interface GroupCreation extends Identifiable {
  spec: {
    /** 展示名 */
    displayName: string;
    /** 描述 */
    description: string;
  };
  status?: {
    /** 用户 */
    users?: {
      id: string;
    }[];
  };
}

export interface GroupEditor extends Identifiable {
  metadata?: {
    /** 名称 */
    name: string;
    /** 创建时间 */
    creationTimestamp?: string;
  };
  spec: {
    /** 展示名 */
    displayName: string;
    /** 描述 */
    description: string;
  };

  /** 是否正在编辑 */
  v_editing?: boolean;
}

export interface GroupAssociation extends Identifiable{
  /** 后端绑定接口不支持同时绑定和解绑，因此，这里设计灵活点，存储原始数据和即将增删的数据 */
  /** 最新数据 */
  groups?: GroupPlain[];
  /** 原来数据 */
  originGroups?: GroupPlain[];
  /** 新增数据 */
  addGroups?: GroupPlain[];
  /** 删除数据 */
  removeGroups?: GroupPlain[];
}
