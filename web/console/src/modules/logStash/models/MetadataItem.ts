import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface MetadataItem extends Identifiable {
  /** 变量Key */
  metadataKey?: string;
  v_metadataKey?: Validation;

  /** 变量值 */
  metadataValue?: string;
  v_metadataValue?: Validation;
}
