import { projectNamespaceActions } from './projectNamespaceActions.project';
import { namespaceActions } from './namespaceActions.project';
import { clusterActions } from './clusterActions';
import { regionActions } from './regionActions';
import { helmActions } from './helmActions';

import { createActions } from './createActions';
import { detailActions } from './detailActions';

export const allActions = {
  cluster: clusterActions,
  region: regionActions,
  helm: helmActions,
  create: createActions,
  detail: detailActions,
  namespace: namespaceActions,
  projectNamespace: projectNamespaceActions
};
