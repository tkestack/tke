import { Identifiable } from '@tencent/qcloud-lib';

export interface Region extends Identifiable {
  /**
   * 地域Id
   */
  value: number | string;

  /**
   * 区域名称
   * */
  name?: string;
}
