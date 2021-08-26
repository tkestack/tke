/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { projectNamespaceActions } from './projectNamespaceActions.project';
import { workloadActions } from './workloadActions';
// import { groupActions } from './groupActions';
import { alarmPolicyActions } from './alarmPolicyActions';
// import { regionActions } from './regionActions';
import { workflowActions } from './workflowActions';
import { validatorActions } from './validatorActions';
import { clusterActions } from './clusterActions';
import { namespaceActions } from './namespaceActions';
import { userActions } from '../../uam/actions/userActions';
import { resourceActions } from '../../notify/actions/resourceActions';

export const allActions = {
  // region: regionActions,
  workflow: workflowActions,
  validator: validatorActions,
  cluster: clusterActions,
  alarmPolicy: alarmPolicyActions,
  // group: groupActions,
  resourceActions: resourceActions,
  user: userActions,
  namespace: namespaceActions,
  workload: workloadActions,
  projectNamespace: projectNamespaceActions
};
