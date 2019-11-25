import { Identifiable, uuid } from '@tencent/qcloud-lib';
import { Validation, initValidator } from './';

export interface Label extends Identifiable {
  /**Lable名称 */
  key?: string;
  v_key?: Validation;

  /**Label值 */
  value?: string;
  v_value?: Validation;
}

export const initLabel = {
  id: uuid(),

  /**Lable名称 */
  key: '',
  v_key: initValidator,

  /**Label值 */
  value: '',
  v_value: initValidator
};
