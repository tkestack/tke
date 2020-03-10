import { Identifiable } from '@tencent/ff-redux';

export interface K8sVersion extends Identifiable {
  /**版本名称 */
  name?: string;

  /**版本号 */
  version?: string;

  /**状态 */
  status?: string;

  /**备注 */
  remark?: string;
}

export interface K8sVersionFilter {
  /**状态 */
  status?: string;

  /**地域 */
  regionId?: string | number;
}
