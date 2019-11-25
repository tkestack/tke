import { Identifiable } from '@tencent/qcloud-lib';
import { Validation } from 'src/modules/common';
export interface AllocationRatioEdition extends Identifiable {
  isUseCpu?: boolean;
  isUseMemory?: boolean;
  cpuRatio?: string;
  v_cpuRatio?: Validation;
  memoryRatio?: string;
  v_memoryRatio: Validation;
}
