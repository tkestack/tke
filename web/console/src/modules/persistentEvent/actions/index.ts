import { clusterActions } from './clusterActions';
import { peActions } from './peActions';
import { regionActions } from './regionActions';
import { workflowActions } from './workflowActions';
import { peEditActions } from './peEditActions';
import { validatorActions } from './validatorActions';

export const allActions = {
  cluster: clusterActions,
  pe: peActions,
  region: regionActions,
  workflow: workflowActions,
  editPE: peEditActions,
  validate: validatorActions
};
