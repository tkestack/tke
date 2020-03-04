import { Identifiable } from '@tencent/ff-redux';

export interface Zone extends Identifiable {
  /**
   * 可用区ID
   */
  id: string;

  /**
   * 可用区名称
   * */
  name?: string;

  /**
   * 可用区状态
   * */
  status?: number;

  /**
   * cbs
   */
  cbs?: number;

  /**
   * 是否默认选中
   */
  isdefault?: boolean;

  /**
   * 是否可用
   */
  disable?: boolean;

  /**
   * 提示
   */
  tip?: string;

  /**
   * 是否下线
   */
  offline?: boolean;

  /**
   * 白名单
   */
  whiteList?: string;
}

export interface ZoneInfo {
  default?: number;

  id?: number;

  name?: string;

  payMode?: string[];

  zoneId?: string;
}

export interface ZoneFilter {
  regionId?: number | string;

  devPayMode?: string;
}

export interface ZoneQuotaFilter {
  cvmPayMode?: number;
}
