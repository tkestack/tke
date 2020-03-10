import { userActions } from './userActions';
import { strategyActions } from './strategyActions';
import { associateActions } from './associatedActions';
import { roleActions } from './role';
import { commonUserActions } from './user';
import { policyActions } from './policy';
import { groupActions } from './group';

export const allActions = {
  user: userActions,
  strategy: strategyActions,
  associateActions: associateActions,

  role: roleActions,
  commonUser: commonUserActions,
  policy: policyActions,
  group: groupActions,
};
