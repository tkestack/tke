import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';
import { Manager } from './Manager';

export interface Project extends Identifiable {
  /** metadata */
  metadata: ProjectMetadata;

  /** spec */
  spec: ProjectSpec;

  /** status */
  status: ProjectStatus;
}

interface ProjectMetadata {
  /** projectId */
  name?: string;

  /** 创建时间 */
  creationTimestamp?: string;

  /** 其余属性 */
  [props: string]: any;
}

interface ProjectSpec {
  /** project名称 */
  displayName: string;

  /** Project成员 */
  members: string[];

  /** project分配的quota */
  clusters: {
    [props: string]: {
      hard?: StatusResource;
    };
  };

  parentProjectName?: string;
}

interface ProjectStatus {
  /** 可分配 */
  used?: StatusResource;
  /** 未分配 */
  calculatedChildProjects?: string;

  calculatedNamespaces?: string[];

  /** 状态 */
  phase?: string;
}

interface StatusResource {
  [props: string]: string;
}

export interface ProjectEdition extends Identifiable {
  id: string;

  resourceVersion: string;

  members: Manager[];

  displayName: string;
  v_displayName: Validation;

  clusters: {
    name: string;
    v_name: Validation;
    resourceLimits: ProjectResourceLimit[];
  }[];

  parentProject: string;

  status: any;
}

export interface ProjectResourceLimit extends Identifiable {
  type: string;
  v_type: Validation;
  value: string;
  v_value: Validation;
}

export interface ProjectFilter {
  /** 业务id */
  ProjectId?: string;

  /**业务名称 */
  displayName?: string;
}
