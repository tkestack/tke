import { Identifiable } from '@tencent/ff-redux';

import { Validation } from './';

export interface EnvItem extends Identifiable {
  /**变量名 */
  envName?: string;
  v_envName?: Validation;

  /**变量值 */
  envValue?: string;
  v_envValue?: Validation;
}
