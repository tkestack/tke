import { Identifiable } from '@tencent/qcloud-lib';
import { ResourceFilter } from './ResourceOption';

export interface Event extends Identifiable {
  /** 事件出现的次数 */
  count?: string;

  /** 首次出现的时间 */
  firstTimestamp?: string;

  /** 最后出现的时间 */
  lastTimestamp?: string;

  /** message */
  message?: string;

  /** metadata */
  metadata?: Metadata;

  /** involvedObject */
  involvedObject?: InvolvedObject;

  /** reason */
  reason?: string;

  /** source */
  source?: any;

  /** 事件的级别 */
  type?: string;
}

interface InvolvedObject {
  /** apiVersion */
  apiVersion?: string;

  /** kind */
  kind?: string;

  /** name */
  name?: string;

  /** namespace */
  namespace?: string;

  /** resourceVersion */
  resourceVersion?: string;

  /** uid */
  uid?: string;
}

interface Metadata {
  /** creationTimestamp */
  creationTimestamp?: string;

  /** name */
  name?: string;

  /** namespace */
  namespace?: string;

  /** resourceVersion */
  resourceVersion?: string;

  /** selfLink */
  selfLink?: string;

  /** uid */
  uid?: string;
}

/** 资源详情页当中的 事件filter数据类型EventFilter */
export interface EventFilter extends ResourceFilter {
  /** kind */
  kind?: string;

  /** workload的名称 */
  name?: string;
}
