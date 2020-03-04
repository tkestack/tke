import { Identifiable } from '@tencent/ff-redux';

export interface Manager extends Identifiable {
  /** 名称 */
  displayName?: string;

  /** id */
  name?: string;
}

export interface ManagerFilter {}
