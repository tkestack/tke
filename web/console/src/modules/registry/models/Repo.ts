import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface Repo extends Identifiable {
  apiVersion?: string;

  kind?: boolean;

  /** 元数据 */
  metadata?: {
    annotations?: any;
    clusterName?: string;
    creationTimestamp?: string;
    deletionGracePeriodSeconds?: number;
    deletionTimestamp?: string;
    finalizers?: string[];
    generateName?: string;
    generation?: string;
    labels?: any;
    managedFields?: any[];
    /** 键值 name */
    name?: string;
    namespace?: string;
    ownerReferences?: any;
    resourceVersion?: string;
    selfLink?: string;
    uid?: string;
  };

  spec?: {
    /** 描述 */
    displayName?: string;
    /** 仓库名称 */
    name?: string;
    tenantID?: string;
    visibility?: string;
  };

  status?: {
    locked?: boolean;
    /** 仓库数量 */
    repoCount?: number;
  };
}

export interface RepoFilter {
  /** 仓库名称 */
  name?: string;
  /** 描述 */
  displayName?: string;
}

export interface RepoCreation extends Identifiable {
  /** 描述 */
  displayName?: 'string';
  /** 仓库名称 */
  name?: 'string';
  v_name?: Validation;
  /** 公开或私有 */
  visibility?: 'Public' | 'Private';
}
