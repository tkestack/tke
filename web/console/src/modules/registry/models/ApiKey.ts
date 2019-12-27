import { Identifiable } from '@tencent/qcloud-lib';
import { Validation } from '../../common/models';

export interface ApiKey extends Identifiable {
  metadata?: {
    name: string;
  };

  spec: {
    /** key 内容 */
    apiKey?: string;

    /** 描述 */
    description?: string;

    /** 过期时间 */
    expire_at?: string;

    /** 创建时间 */
    issue_at?: string;
  };

  status: {
    /** 启用|禁用 */
    disabled?: boolean;

    /** 是否过期 */
    expired?: boolean;
  };

  /** 软删除标记? */
  deleted?: boolean;
}

export interface ApiKeyFilter {
  /** 描述字段 */
  desc?: string;
}

export interface ApiKeyCreation extends Identifiable {
  /** key 描述 */
  description?: string;
  /** key 过期时间，单位 h */
  expire?: number;
  v_expire: Validation;
  unit?: string;
}
