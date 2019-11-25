import { Identifiable } from '@tencent/qcloud-lib';

export interface TableFilterOption extends Identifiable {
  /**显示名称 */
  label?: string;

  /**是否默认 */
  default?: boolean;
}
