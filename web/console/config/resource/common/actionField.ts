import { ActionField } from '../../../src/modules/common/models';

/** resource list资源页面的所提供的操作按钮 */
export const commonActionField: ActionField = {
  create: {
    isAvailable: true
  },
  search: {
    isAvailable: true,
    attributes: []
  },
  manualRenew: {
    isAvailable: true,
    attributes: []
  },
  autoRenew: {
    isAvailable: false
  },
  download: {
    isAvailable: true
  }
};
