import { Identifiable, uuid } from '@tencent/qcloud-lib';
import { Validation, initValidator } from '../../common/models';

export interface ConfigMapEdit extends Identifiable {
  /** cm的名称 */
  name?: string;
  v_name?: Validation;

  /** namespace */
  namespace?: string;

  /** 配置内容 */
  variables?: Array<Variable>;
}

export interface Variable extends Identifiable {
  /** key */
  key?: string;
  v_key?: Validation;

  /** value */
  value?: string;
}

export const initVariable = {
  id: uuid(),
  key: '',
  v_key: initValidator,
  value: ''
};
