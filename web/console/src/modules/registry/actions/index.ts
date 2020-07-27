import { apiKeyActions } from './apiKeyActions';
import { repoActions } from './repoActions';
import { imageActions } from './imageActions';
import { chartActions as chart } from './chartActions';
import { chartInsActions } from './chartInsActions';
import { chartGroupActions } from './chartGroup';
import { projectActions } from './project';
import { userActions } from './user';
import { chartActions } from './chart';
import { appActions } from './app';
import { clusterActions } from './cluster';
import { namespaceActions, projectNamespaceActions } from './namespace';

export const allActions = {
  apiKey: apiKeyActions,
  repo: repoActions,
  image: imageActions,
  charts: chart,
  chartIns: chartInsActions,

  chartGroup: chartGroupActions,
  chart: chartActions,
  project: projectActions,
  user: userActions,
  app: appActions,
  cluster: clusterActions,
  namespace: namespaceActions,
  projectNamespace: projectNamespaceActions
};
