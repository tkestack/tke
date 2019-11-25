import { regionActions } from './regionActions';
import { clusterActions } from './clusterActions';
import { addonActions } from './addonActions';
import { addonEditActions } from './addonEditActions';
import { workflowActions } from './workflowActions';
import { validatorActions } from './validatorActions';

export const allActions = {
  region: regionActions,
  cluster: clusterActions,
  addon: addonActions,
  editAddon: addonEditActions,
  workflow: workflowActions,
  validator: validatorActions
};
