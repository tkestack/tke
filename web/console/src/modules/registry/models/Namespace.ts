import { Identifiable } from '@tencent/ff-redux';
import { Validation } from '../../common/models';
import { ChartInfoFilter } from './Chart';

export interface Namespace extends Identifiable {
  metadata?: {
    name?: string;
  };
}

export interface ProjectNamespace extends Identifiable {
  metadata?: {
    name?: string;
    namespace?: string;
  };
  spec?: {
    clusterName?: string;
    namespace?: string;
  };
}

export interface NamespaceFilter {
  cluster?: string;
  chartInfoFilter?: ChartInfoFilter;
}

export interface ProjectNamespaceFilter {
  projectId?: string;
  chartInfoFilter?: ChartInfoFilter;
}
