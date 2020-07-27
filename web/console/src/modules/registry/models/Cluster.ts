import { Identifiable } from '@tencent/ff-redux';
import { Validation } from '../../common/models';
import { ChartInfoFilter } from './Chart';

export interface Cluster extends Identifiable {
  metadata?: {
    name?: string;
  };

  spec?: {
    displayName?: string;
    tenantID?: string;
    type?: string;
    version?: string;
  };
}

export interface ClusterFilter {
  chartInfoFilter?: ChartInfoFilter;
}
