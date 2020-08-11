import { Identifiable } from '@tencent/ff-redux';

export interface ChartGroup extends Identifiable {
  metadata?: {
    name: string;
    creationTimestamp?: string;
  };
  spec: {
    name: string;
    tenantID?: string;
    displayName?: string;
    visibility: string;
    description?: string;
    type: string;
    projects?: string[];
  };
  status?: {
    chartCount?: number;
    phase: string;
    [props: string]: any;
  };
}

export interface ChartGroupFilter {
  repoType?: string;
}

export interface ChartGroupDetailFilter {
  name: string;
  projectID: string;
}

export interface ChartGroupCreation extends Identifiable {
  spec: {
    name: string;
    displayName?: string;
    visibility: string;
    description?: string;
    type: string;
    projects?: string[];
  };
}

export interface ChartGroupEditor extends Identifiable {
  metadata?: {
    name: string;
    creationTimestamp?: string;
  };
  spec: {
    name: string;
    displayName?: string;
    visibility: string;
    description?: string;
    type: string;
    projects?: string[];
  };

  /** 是否正在编辑 */
  v_editing?: boolean;
}
