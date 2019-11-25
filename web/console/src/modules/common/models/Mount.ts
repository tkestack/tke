import { Validation } from './';
import { Identifiable } from '@tencent/qcloud-lib';

export interface MountItem extends Identifiable {
  /**数据卷 */
  volume?: string;
  v_volume?: Validation;

  /**目标路径 */
  mountPath?: string;
  v_mountPath?: Validation;

  /**权限 */
  mode?: string;
  v_mode: Validation;
}
