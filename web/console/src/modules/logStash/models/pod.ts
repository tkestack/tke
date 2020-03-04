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
export interface PodListFilter {
  /** 命名空间 */
  namespace?: string;

  /** 集群id */
  clusterId?: string;

  /** 地域id */
  regionId?: number;

  /** name */
  specificName?: string;

  isCanFetchPodList?: boolean;
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
