import { Validation } from 'src/modules/common';

import { Identifiable } from '@tencent/ff-redux';

export interface AllocationRatioEdition extends Identifiable {
  isUseCpu?: boolean;
  isUseMemory?: boolean;
  cpuRatio?: string;
  v_cpuRatio?: Validation;
  memoryRatio?: string;
  v_memoryRatio: Validation;
}
