import { Identifiable } from '@tencent/qcloud-lib';
import { Validation } from '../../common/models';

export interface Namespace extends Identifiable {
  /** 命名空间名称 */
  name: string;
}
