import { Identifiable } from '@tencent/ff-redux';

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
