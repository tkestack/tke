import { Validation } from '../../common';

export interface PeEdit {
  /** es的地址 */
  esAddress?: string;
  v_esAddress?: Validation;

  /** 索引名称 */
  indexName?: string;
  v_indexName?: Validation;

  /** ES 认证用户名 */
  esUsername?: string;

  /** ES 认证密码 */
  esPassword?: string;
}
