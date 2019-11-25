import { Identifiable } from '@tencent/qcloud-lib';

export interface Manager extends Identifiable {
  /** 名称 */
  displayName?: string;

  /** id */
  name?: string;
}

export interface ManagerFilter {}
