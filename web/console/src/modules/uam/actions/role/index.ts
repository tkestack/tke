import { listActions } from './listActions';
import { createActions } from './createActions';
import { detailActions } from './detailActions';
import { associateActions } from './associateActions';

export const roleActions = {
  associate: associateActions,
  list: listActions,
  create: createActions,
  detail: detailActions,
};
