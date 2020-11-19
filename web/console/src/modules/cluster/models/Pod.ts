import { Identifiable } from '@tencent/ff-redux';

export interface Pod extends Identifiable {
  /** metadata */
  metadata?: Metadata;

  /** spec */
  spec?: Spec;

  /** status */
  status?: Status;
}

interface Metadata {
  annotations?: {
    'kubernetes.io/created-by'?: string;

    [props: string]: string;
  };

  creationTimestamp?: string;

  name?: string;

  namespace?: string;

  [props: string]: any;
}

interface Spec {
  containers?: PodContainer[];

  [props: string]: any;
}

interface Status {
  containerStatuses?: any[];

  conditions?: any[];

  phase?: string;

  qosClass?: string;

  /** pod所在node 的ip */
  hostIP?: string;

  /** pod的ip */
  podIP?: string;

  /** pod启动时间 */
  startTime?: string;
}

export interface PodContainer extends Identifiable {
  env?: Env[];

  image?: string;

  imagePullPolicy?: string;

  name?: string;

  resources?: any;

  terminationMessagePath?: string;

  terminationMessagePolicy?: string;

  [props: string]: any;
}

interface Env {
  name?: string;

  value?: string;
}

/** Node详情页内的pod列表的筛选框的项 */
export interface PodFilterInNode {
  /** 命名空间 */
  namespace?: string;

  /** podName */
  podName?: string;

  /** pod的状态值 */
  phase?: string;

  ip?: string;
}
