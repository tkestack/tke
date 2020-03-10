import { Identifiable } from '@tencent/ff-redux';

/** 原有User在localidentities和user概念之间混用了，anyway，用于关联角色等 */
export interface UserPlain extends Identifiable {
  /** 名称 */
  name?: string;
  /** 展示名 */
  displayName?: string;
}

export interface CommonUserFilter {
  /** 目标资源，如localgroup/role/policy */
  resource: string;
  /** 资源id */
  resourceID: string;
  /** 关联/解关联操作后的回调函数 */
  callback?: () => void;
}

export interface CommonUserAssociation extends Identifiable{
  /** 后端绑定接口不支持同时绑定和解绑，因此，这里设计灵活点，存储原始数据和即将增删的数据 */
  /** 最新数据 */
  users?: UserPlain[];
  /** 原来数据 */
  originUsers?: UserPlain[];
  /** 新增数据 */
  addUsers?: UserPlain[];
  /** 删除数据 */
  removeUsers?: UserPlain[];
}
