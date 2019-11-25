import { Identifiable } from '@tencent/qcloud-lib';
export interface Category extends Identifiable {
  /** value值 */
  name: string;

  /** 展示名称 */
  displayName: string;

  description: string;

  /** 操作 */
  actions: object;
  [props: string]: any;
}
