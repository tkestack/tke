import { Identifiable } from '@tencent/ff-redux';
import { Validation } from '../../common/models';

export interface Policy extends Identifiable {
  metadata?: {
    /** id */
    name: string;
  };
  spec: {
    /** 名称 */
    displayName: string;
    /** 类别 */
    category: string;
    /** 描述 */
    description: string;
    /** 类型 */
    type?: string;
    /** 语法 */
    statement: {
      /** 资源列表 */
      resources: string[];
      /** 结果 */
      effect: string;
      /** 动作 */
      actions: string[];
    };
  };
}

export interface PolicyPlain extends Identifiable {
  /** id */
  name?: string;
  /** 名称 */
  displayName?: string;
  /** 类别 */
  category?: string;
  /** 描述 */
  description: string;
}

export interface PolicyFilter {
  /** 目标资源，如role,localidentity,localgroup */
  resource: string;
  /** 资源id */
  resourceID: string;
  /** 关联/解关联操作后的回调函数 */
  callback?: () => void;
}

/** 用于单个查询 */
export interface PolicyInfoFilter {
  name: string;
}

export interface PolicyEditor extends Identifiable {
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

export interface PolicyAssociation extends Identifiable{
  /** 后端绑定接口不支持同时绑定和解绑，因此，这里设计灵活点，存储原始数据和即将增删的数据 */
  /** 最新数据 */
  policies?: PolicyPlain[];
  /** 原来数据 */
  originPolicies?: PolicyPlain[];
  /** 新增数据 */
  addPolicies?: PolicyPlain[];
  /** 删除数据 */
  removePolicies?: PolicyPlain[];
}
