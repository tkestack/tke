import { projectActions } from './projectActions';
import { namespaceActions } from './namespaceActions';
import { regionActions } from './regionActions';
import { clusterActions } from './clusterActions';
import { managerActions } from './managerActions';
import { userActions } from './userActions';
import { policyActions } from './policy';
import { roleActions } from './role';
import { groupActions } from './group';

export const allActions = {
  project: projectActions,
  manager: managerActions,
  namespace: namespaceActions,
  region: regionActions,
  cluster: clusterActions,
  user: userActions,
  policy: policyActions,
  role: roleActions,
  group: groupActions
};
