import { Identifiable } from '@tencent/qcloud-lib';
import { Cluster } from '../../common/models';

export interface Ckafka extends Identifiable {
  /** ckafka的 instanceId */
  instanceId?: string;

  /** ckafka instanceName */
  instanceName?: string;

  /**实例状态 0: 创建中，1：运行中， 2：删除中 */
  status?: number;

  bandwith?: number;

  diskSize?: number;

  /** vpc网络id */
  vpcId?: string;

  /** 子网id */
  subnitId?: string;

  /** zoneId */
  zoneId?: number;

  /** topic 的数量 */
  topicNum?: number;

  /** vipList */
  vipList?: any[];
}

export interface CkafkaFilter {
  /** ckafka的状态 */
  status?: number;

  /** 集群的信息 */
  cluster?: Cluster;

  /** 当前的地域ID */
  regionId?: number;

  /** 是否能够拉取ckafka的列表 */
  isCanFetchCkafka?: boolean;
}

export interface CTopic extends Identifiable {
  /** topicId */
  topicId?: string;

  /** topicName */
  topicName?: string;
}

export interface CTopicFilter {
  /** instanceId */
  instanceId?: string;

  /** 当前的地域ID */
  regionId?: number;

  /** 是否能够拉取CTopic的列表 */
  isCanFetchCTopic?: boolean;
}
