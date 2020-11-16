import { Identifiable } from '@tencent/ff-redux';
import { Validation } from '../../common/models';

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
    clusterDisplayName?: string;
  };
}

export interface NamespaceFilter {
  cluster?: string;
}

export interface ProjectNamespaceFilter {
  projectId?: string;
}
