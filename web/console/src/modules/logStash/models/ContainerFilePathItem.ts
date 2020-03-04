import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface ContainerFilePathItem extends Identifiable {
  containerName?: string;
  containerFilePath?: string;
  v_containerFilePath?: Validation;
  v_containerName?: Validation;
}
