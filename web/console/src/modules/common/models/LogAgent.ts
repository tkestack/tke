export interface LAMetadata {
  name?: string;

  generateName?: string;

  selfLink?: string;

  uid?: string;

  resourceVersion?: string;

  creationTimestamp?: string;
}

export interface LASpec {
  tenantID?: string;

  clusterName?: string;

  version?: string;
}

export interface LAStatus {
  version?: string;

  phase?: string;

  retryCount?: number;

  lastReInitializingTimestamp?: any;
}

/**
 * LogAgent 的结构定义
 */
export interface LogAgent {
  metadata?: LAMetadata;

  spec?: LASpec;

  status?: LAStatus;
}
