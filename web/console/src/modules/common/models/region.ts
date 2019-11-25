import { Identifiable } from '@tencent/qcloud-lib';

export interface Region extends Identifiable {
  /** 地域的值 */
  value: number;

  /** 地域的名称 */
  name?: string;

  /** 地域是否可用 */
  disabled?: boolean;

  /** 地域所属大区 */
  area?: string;

  /** 是否新版控制台当中的地域 */
  Remark?: string;

  //[props: string]: any;
}

export interface RegionFilter {
  /** 地域id */
  regionId?: number;
}
