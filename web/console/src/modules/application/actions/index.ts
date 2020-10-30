import { appActions } from './app';
import { clusterActions } from './cluster';
import { namespaceActions, projectNamespaceActions } from './namespace';
import { chartActions } from './chart';
import { chartGroupActions } from './chartGroup';
import { projectActions } from './project';

export const allActions = {
  app: appActions,
  cluster: clusterActions,
  namespace: namespaceActions,
  projectNamespace: projectNamespaceActions,
  chart: chartActions,
  chartGroup: chartGroupActions,
  project: projectActions
};
