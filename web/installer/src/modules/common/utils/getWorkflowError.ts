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

import { WorkflowState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export const getWorkflowError = (workflow: WorkflowState<any, any>) => {
  let error = workflow && workflow.results && workflow.results.length && workflow.results[0].error;

  return error || {};
};

/**获取错误码 */
export const getWorkflowErrorCode = (workflow: WorkflowState<any, any>) => {
  let reg = /\(-\w+\)/g,
    code = 0;
  if (workflow && workflow.results && workflow.results.length && workflow.results[0].error) {
    let msg = workflow.results[0].error.message || '',
      matches = msg.match(reg);
    if (matches && matches[0]) {
      let codeStr = matches[0].substring(1, matches[0].length - 1);
      code = parseInt(codeStr);
    }
  }

  return code;
};
