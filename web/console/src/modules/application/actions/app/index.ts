import { listActions } from './listActions';
import { createActions } from './createActions';
import { detailActions } from './detailActions';
import { resourceActions } from './resourceActions';
import { historyActions } from './historyActions';

export const appActions = {
  list: listActions,
  create: createActions,
  detail: detailActions,
  resource: resourceActions,
  history: historyActions
};
