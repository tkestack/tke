import { projectActions } from './projectActions';
import { namespaceActions } from './namespaceActions';
import { regionActions } from './regionActions';
import { clusterActions } from './clusterActions';
import { managerActions } from './managerActions';

export const allActions = {
  project: projectActions,
  manager: managerActions,
  namespace: namespaceActions,
  region: regionActions,
  cluster: clusterActions
};
