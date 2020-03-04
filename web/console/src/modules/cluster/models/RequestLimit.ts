import { Identifiable } from '@tencent/ff-redux';

export interface RequestLimit extends Identifiable {
  /**集群ID */
  clusterId?: string;

  /**Cpu Request总和*/
  totalCpuRequest?: number;

  /**内存 Request总和*/
  totalMemRequest?: number;

  /**总内存 */
  totalCpu?: number;

  /**总cpu */
  totalMem?: number;

  /**总gpu */
  totalGpu?: number;

  /**错误信息 */
  result?: any;
}

export interface RequestLimitFilter {
  regionId?: number;

  clusterIds?: string[];
}
