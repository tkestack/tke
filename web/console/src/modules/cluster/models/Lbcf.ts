import { extend, Identifiable, RecordSet } from '@tencent/ff-redux';

export interface LbcfResource extends Identifiable {
  /** metadata */
  metadata?: {
    name: string;
    namespace: string;
  };

  /** spec */
  spec?: {
    lbDriver: string;
    lbSpec: {
      lbID?: string;
      lbVpcID?: string;
    };
    backGroups: BackendGroup[];
  };

  /** status */
  status?: any;

  /** other */
  [props: string]: any;
}

export interface BackendGroup {
  name: string;

  pods?: PodBackend;

  service?: ServiceBackend;

  static?: string[];

  status: {
    backends: number;
    registeredBackends: number;
  };
  backendRecords: BackendRecord[];
}

export interface PodBackend {
  labels: {
    [props: string]: string;
  };
  port: {
    portNumber: number;
    protocol: string;
  };
  byName: string[];
}

export interface ServiceBackend {
  name: string;
  port: {
    portNumber: number;
    protocol: string;
  };
  nodeSelector: {
    [props: string]: string;
  };
}

export interface BackendRecord {
  name: string;
  backendAddr: string;
  conditions: Condition[];
}

interface Condition {
  lastTransitionTime: string;
  message: string;
  status: string;
  type: string;
}
