import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface Image extends Identifiable {
  apiVersion?: string;
  kind?: string;
  metadata?: {
    annotations?: any;
    clusterName?: string;
    creationTimestamp?: string;
    generateName?: string;
    generation?: number;
    name?: 'string';
    namespace?: 'string';
    resourceVersion?: 'string';
    selfLink?: 'string';
    uid?: 'string';
  };
  spec?: {
    displayName?: 'string';
    name?: 'string';
    namespaceName?: 'string';
    tenantID?: 'string';
    visibility?: 'Public' | 'Private';
  };
  status?: {
    locked?: boolean;
    pullCount?: number;
    tags?: Tag[];
  };
}

export interface Tag extends Identifiable {
  digest?: string;
  name?: string;
  timeCreated?: string;
}

export interface ImageFilter {
  namespace?: string;
  namespaceName?: string;
}

export interface ImageCreation extends Identifiable {
  // kind: Repository
  // metadata?: {
  namespace?: string;
  // };
  // spec?: {
  displayName?: string;
  name?: string;
  v_name?: Validation;
  namespaceName?: string;
  visibility?: 'Public' | 'Private';
  // };
}
