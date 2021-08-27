/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import { clusterCreationAction } from './clusterCreationAction';
import { clusterActions } from './clusterActions';
import { computerActions } from './computerActions';
import { computerPodActions } from './computerPodActions';
import { modeActions } from './modeActions';
import { namespaceActions } from './namespaceActions';
import { namespaceEditActions } from './namespaceEditActions';
import { regionActions } from './regionActions';
import { resourceActions } from './resourceActions';
import { resourceDetailActions } from './resourceDetailActions';
import { serviceEditActions } from './serviceEditActions';
import { subRouterActions } from './subRouterActions';
import { validatorActions } from './validatorActions';
import { workflowActions } from './workflowActions';
import { workloadEditActions } from './workloadEditActions';
import { cmEditActions } from './cmEditActions';
import { resourceLogActions } from './resourceLogActions';
import { resourceEventActions } from './resourceEventActions';
import { secretEditActions } from './secretEditActions';
import { dialogActions } from './dialogActions';
import { createICAction } from './createICAction';
import { lbcfEditActions } from './lbcfEditActions';

import { projectNamespaceActions } from './projectNamespaceActions.project';

export const allActions = {
  dialog: dialogActions,
  region: regionActions,
  cluster: clusterActions,
  computer: computerActions,
  resource: resourceActions,
  namespace: namespaceActions,
  workflow: workflowActions,
  subRouter: subRouterActions,
  resourceDetail: resourceDetailActions,
  validate: validatorActions,
  mode: modeActions,
  computerPod: computerPodActions,
  editSerivce: serviceEditActions,
  editNamespace: namespaceEditActions,
  editWorkload: workloadEditActions,
  editCM: cmEditActions,
  resourceLog: resourceLogActions,
  resourceEvent: resourceEventActions,
  editSecret: secretEditActions,
  validator: validatorActions,
  clusterCreation: clusterCreationAction,
  projectNamespace: projectNamespaceActions,
  createIC: createICAction,
  lbcf: lbcfEditActions
};
