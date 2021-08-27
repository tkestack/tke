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
import { createFFListActions } from '@tencent/ff-redux';

import { Project, ProjectEdition, ProjectFilter, RootState } from '../models';
import { Manager } from '../models/Manager';
import { ProjectResourceLimit, ProjectUserMap } from '../models/Project';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
type GetState = () => RootState;

const FFModelProjectActions = createFFListActions<Project, ProjectFilter>({
  actionName: 'detailProject',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchProjectList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().detailProject;
  },
  keepLastSelection: true,
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState(),
      urlParams = router.resolve(route);
    if (record.data.records.filter(item => item.status.phase !== 'Active').length === 0) {
      dispatch(FFModelProjectActions.clearPolling());
    }
  }
});

export const detailActions = {
  project: FFModelProjectActions
};
