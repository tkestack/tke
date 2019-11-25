export interface UserFilter {
  /** 用户名(唯一) */
  name?: string;

  /** 展示名 */
  displayName?: string;

  /** 相关参数 */
  search?: string;

  ifAll?: boolean;
}
