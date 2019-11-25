import { workloadActions } from './workloadActions';
// import { groupActions } from './groupActions';
import { alarmPolicyActions } from './alarmPolicyActions';
// import { regionActions } from './regionActions';
import { workflowActions } from './workflowActions';
import { validatorActions } from './validatorActions';
import { clusterActions } from './clusterActions';
import { namespaceActions } from './namespaceActions';
import { userActions } from '../../uam/actions/userActions';
import { resourceActions } from '../../notify/actions/resourceActions';

export const allActions = {
  // region: regionActions,
  workflow: workflowActions,
  validator: validatorActions,
  cluster: clusterActions,
  alarmPolicy: alarmPolicyActions,
  // group: groupActions,
  resourceActions: resourceActions,
  user: userActions,
  namespace: namespaceActions,
  workload: workloadActions
};
