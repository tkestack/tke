import { userActions } from './userActions';
import { strategyActions } from './strategyActions';
import { associateActions } from './associatedActions';
export const allActions = {
  user: userActions,
  strategy: strategyActions,
  associateActions: associateActions
};
