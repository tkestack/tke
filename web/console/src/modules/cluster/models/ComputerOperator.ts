export interface ComputerOperator {
  /**
   * 集群Id
   */
  clusterId?: string;

  /**
   * 地域
   */
  regionId?: number;

  /**
   * 移出操作方式
   */
  nodeDeleteMode?: string;

  /**
   * node是否进行 Unschedule 的操作
   */
  isUnSchedule?: boolean;
}
