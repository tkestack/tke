import { Identifiable } from '@tencent/ff-redux';
import { Resource } from './Resource';

export interface ClusterFilter {
  /** 具体名称 */
  specificName?: string;

  /** 地域id */
  regionId?: number;
}

export interface Cluster extends Identifiable {
  /** metadata */
  metadata: ClusterMetadata;

  /** spec */
  spec: ClusterSpec;

  /** status */
  status: ClusterStatus;
}

interface ClusterMetadata {
  /** 集群id */
  name?: string;

  /** 创建时间 */
  creationTimestamp?: string;

  /** 其余属性 */
  [props: string]: any;
}

interface ClusterSpec {
  /** 集群名称 */
  displayName?: string;

  /** 集群的features */
  features?: {
    ipvs: boolean;
    public: boolean;
  };

  /** 集群类型 */
  type?: string;

  /** 集群的版本 */
  version?: string;

  /** 是否安装了prometheus */
  hasPrometheus?: boolean;

  /** promethus详情 */
  promethus?: Resource;

  /** logagent详情 */
  logAgent?: Resource;

  properties?: any;

  [props: string]: any;
}

interface ClusterStatus {
  /** 集群的地址相关信息 */
  addresses?: ClusterAddress[];

  /** 集群当前的状态 */
  conditions?: ClusterCondition[];

  /** 集群的相关凭证 */
  credential?: {
    caCert: string;
    token: string;
  };

  /** 当前的状态 */
  phase?: string;

  /** 资源的相关配置 request limit */
  resource?: {
    /** 可分配 */
    allocatable: StatusResource;

    /** 已分配 */
    allocated: StatusResource;

    /** 集群的配置 */
    capacity: StatusResource;
  };

  /** 集群的版本 */
  version?: string;
}

interface StatusResource {
  /** cpu的相关配置 */
  cpu: string;

  /** memory的相关配置 */
  memory: string;
}

interface ClusterAddress {
  /** 集群的域名 */
  host: string;

  /** 端口名 */
  port: number;

  path?: string;

  /** 集群的类型 */
  type: string;
}

export interface ClusterCondition {
  /** 上次健康检查时间 */
  lastProbeTime?: string;
  /**节点可能是心跳时间 */
  lastHeartbeatTime?: string;

  /** 上次变更时间 */
  lastTransitionTime?: string;

  /** 原因 */
  reason?: string;

  /** 错误信息 */
  message?: string;

  /** 状态是否正常 */
  status?: string;

  /** 条件类型 */
  type?: string;
}

export interface RegionCluster {
  /**地域 */
  region: string | number;

  /**集群列表 */
  clusters: Cluster[];
}

export interface ClusterOperator {
  /**所属地域 */
  regionId?: number | string;
}
