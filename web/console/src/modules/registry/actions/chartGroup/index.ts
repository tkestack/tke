import { listActions } from './listActions';
import { createActions } from './createActions';
import { detailActions } from './detailActions';

export const chartGroupActions = {
  list: listActions,
  create: createActions,
  detail: detailActions
};
