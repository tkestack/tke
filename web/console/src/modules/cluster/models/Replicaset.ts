import { Identifiable } from '@tencent/ff-redux';

export interface Replicaset extends Identifiable {
  /** metadata */
  metadata?: Metadata;

  /** spec */
  spec?: Spec;

  /** status */
  status?: Status;
}

interface Metadata {
  /** annotations */
  annotations?: {
    'deployment.kubernetes.io/desired-replicas'?: string;

    'deployment.kubernetes.io/max-replicas'?: string;

    'deployment.kubernetes.io/revision'?: string;

    [props: string]: string;
  };

  /** creationTimestamp */
  creationTimestamp?: string;

  /** generation */
  generation?: string;

  /** labels */
  labels?: {
    [props: string]: string;
  };

  /**name */
  name?: string;

  namespace?: string;

  [props: string]: any;
}

interface Spec {
  /** replicas */
  replicas?: string;

  /** selector */
  selector?: {
    matchLabels: {
      [props: string]: string;
    };
  };

  /** template */
  template?: {
    metadata?: {
      [props: string]: any;
    };

    spec?: {
      [props: string]: any;
    };
  };
}

interface Status {
  /** availableReplicas */
  availableReplicas?: string;

  /** conditions */
  conditions?: any;

  /** fullyLabeledReplicas */
  fullyLabeledReplicas?: string;

  /** observedGeneration */
  observedGeneration?: string;

  /** readyReplicas */
  readyReplicas?: string;

  /** replicas */
  replicas?: string;
}
