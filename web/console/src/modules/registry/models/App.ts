import { Identifiable } from '@tencent/ff-redux';
import { Validation } from '../../common/models';

export interface App extends Identifiable {
  metadata?: {
    creationTimestamp?: string;
    name?: string;
    namespace?: string;
    generation?: number;
  };

  spec?: {
    chart?: {
      chartGroupName?: string;
      chartName?: string;
      chartVersion?: string;
      tenantID?: string;
    };
    name?: string;
    targetCluster?: string;
    tenantID?: string;
    type?: string;
    values?: {
      rawValues?: string;
      rawValuesType?: string;
      values?: string[];
    };
  };
  status?: {
    lastTransitionTime?: string;
    phase?: string;
    releaseLastUpdated?: string;
    releaseStatus?: string;
    revision?: number;
    observedGeneration?: number;
  };
}

export interface AppCreation extends Identifiable {
  metadata: {
    namespace: string;
  };

  spec?: {
    chart?: {
      chartGroupName?: string;
      chartName?: string;
      chartVersion?: string;
      tenantID?: string;
    };
    tenantID?: string;
    name?: string;
    targetCluster?: string;
    type?: string;
    values?: {
      rawValues?: string;
      rawValuesType?: string;
      values?: string[];
    };
  };
}

export interface AppFilter {
  cluster: string;
  namespace: string;
  projectId?: string;
}
