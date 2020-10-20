import { Identifiable } from '@tencent/ff-redux';
import { Validation } from '../../common/models';

export interface Chart extends Identifiable {
  metadata?: {
    creationTimestamp?: string;
    name?: string;
    namespace?: string;
  };

  spec?: {
    chartGroupName?: string;
    displayName?: string;
    name?: string;
    tenantID?: string;
    visibility?: string;
  };

  status?: {
    pullCount?: number;
    versions?: ChartVersion[];
  };

  // custom: store last version data
  lastVersion?: ChartVersion;
  sortedVersions?: ChartVersion[];
  projectID?: string;
}

export interface ChartVersion {
  chartSize?: number;
  description?: string;
  timeCreated?: string;
  version?: string;
  icon?: string;
  appVersion?: string;
}

export interface ChartFilter {
  namespace?: string;
  repoType?: string;
  projectID?: string;
}

export interface ChartInfo {
  metadata?: {
    name?: string;
    namespace?: string;
  };

  spec: {
    files: {
      [props: string]: string;
    };
    values: {
      [props: string]: string;
    };
  };
}

export interface ChartInfoFilter {
  cluster: string;
  namespace: string;
  metadata: {
    namespace: string;
    name: string;
  };
  chartVersion: string;
  projectID: string;
}
