import { regionActions } from './regionActions';
import { clusterActions } from './clusterActions';
import { logActions } from './logActions';
import { workflowActions } from './workflowActions';
import { editLogStashActions } from './editLogStashActions';
import { validatorActions } from './validatorActions';
import { resourceActions } from './resourceActions';
import { namespaceActions } from './namespaceActions';
import { podActions } from './podActions';
import { logDaemonsetActions } from './logDaemonsetActions';
import { projectNamespaceActions } from '@src/modules/helm/actions/projectNamespaceActions.project';
export const allActions = {
  region: regionActions,
  cluster: clusterActions,
  log: logActions,
  logDaemonset: logDaemonsetActions,
  workflow: workflowActions,
  editLogStash: editLogStashActions,
  validate: validatorActions,
  resource: resourceActions,
  namespace: namespaceActions,
  pod: podActions
};
