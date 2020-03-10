import { cluster } from 'config/resource/k8sConfig';

import { Identifiable } from '@tencent/ff-redux';

export interface LogDaemonset extends Identifiable {
  /** kind */
  kind?: string;

  /**apiVersion */
  apiVersion?: string;

  /** metadata */
  metadata?: Metadata;

  /** spec */
  spec?: Spec;

  /** status */
  status?: Status;
}

interface Metadata {
  creationTimestamp?: string;

  name?: string;

  [props: string]: any;
}

interface Spec {
  clusterName?: string;

  [props: string]: any;
}

interface Status {
  phase?: string;

  reason?: string;

  retryCount?: number;

  [props: string]: any;
}

export interface LogDaemonSetFliter {
  clusterId?: string;

  specificName?: string;
}

export interface LogDaemonSetStatus {
  phase?: string;
  reason?: string;
}
