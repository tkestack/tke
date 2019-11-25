import { Identifiable } from '@tencent/qcloud-lib';
import { any } from 'prop-types';
import { LogStashSpec } from './LogStashEdit';

export interface Log extends Identifiable {
  apiVersion?: string;
  kind?: string;
  metadata?: {
    creationTimestamp: string;
    name: string;
    namespace: string;
    resourceVersion?: string;
  };
  spec: {
    input: {
      type?: string;
      [props: string]: any;
    };
    output: {
      type?: string;
      [props: string]: any;
    };
  };
  [props: string]: any;
}

export interface LogFilter {
  /** 地域的id */
  regionId?: number;

  /** 日志收集器的ID*/
  collectorId?: string;

  /** 日志收集器所属的集群ID*/
  clusterId?: string;

  /** 根据状态进行复选 */
  status?: string;

  /** 是否清除*/
  isClear?: boolean;

  /**命名空间 */
  namespace?: string;
}

export interface LogOperator {
  /**
   * 地域
   */
  regionId?: number;

  /**
   * 集群id
   */
  clusterId?: string;

  /**
   * 当前的编辑类型 create | update
   */
  mode?: string;

  /**
   * 日志收集器id
   */
  collectorName?: string;
}
