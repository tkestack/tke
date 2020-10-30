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
}
