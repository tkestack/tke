import { clusterCreationAction } from './clusterCreationAction';
import { clusterActions } from './clusterActions';
import { computerActions } from './computerActions';
import { computerPodActions } from './computerPodActions';
import { modeActions } from './modeActions';
import { namespaceActions } from './namespaceActions';
import { namespaceEditActions } from './namespaceEditActions';
import { regionActions } from './regionActions';
import { resourceActions } from './resourceActions';
import { resourceDetailActions } from './resourceDetailActions';
import { serviceEditActions } from './serviceEditActions';
import { subRouterActions } from './subRouterActions';
import { validatorActions } from './validatorActions';
import { workflowActions } from './workflowActions';
import { workloadEditActions } from './workloadEditActions';
import { cmEditActions } from './cmEditActions';
import { resourceLogActions } from './resourceLogActions';
import { resourceEventActions } from './resourceEventActions';
import { secretEditActions } from './secretEditActions';
import { dialogActions } from './dialogActions';
import { createICAction } from './createICAction';
import { lbcfEditActions } from './lbcfEditActions';

import { projectNamespaceActions } from './projectNamespaceActions.project';

export const allActions = {
  dialog: dialogActions,
  region: regionActions,
  cluster: clusterActions,
  computer: computerActions,
  resource: resourceActions,
  namespace: namespaceActions,
  workflow: workflowActions,
  subRouter: subRouterActions,
  resourceDetail: resourceDetailActions,
  validate: validatorActions,
  mode: modeActions,
  computerPod: computerPodActions,
  editSerivce: serviceEditActions,
  editNamespace: namespaceEditActions,
  editWorkload: workloadEditActions,
  editCM: cmEditActions,
  resourceLog: resourceLogActions,
  resourceEvent: resourceEventActions,
  editSecret: secretEditActions,
  validator: validatorActions,
  clusterCreation: clusterCreationAction,
  projectNamespace: projectNamespaceActions,
  createIC: createICAction,
  lbcf: lbcfEditActions
};
