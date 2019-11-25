import { ResourceInfo } from '../../common';

export interface DetailResourceOption {
  /**detail 其他资源resource和selection */
  detailResourceName?: string;

  detailResourceInfo?: ResourceInfo;

  detailResourceSelection?: string;

  detailResourceList?: any[];

  detailDeleteResourceSelection?: string;
}
