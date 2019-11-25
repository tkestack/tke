import { Identifiable, RecordSet, extend } from '@tencent/qcloud-lib';

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
  labels: {
    [props: string]: string;
  };
  port: {
    portNumber: number;
    protocol: string;
  };
  status: {
    backends: number;
    registeredBackends: number;
  };
  backendRecords: BackendRecord[];
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
