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
import { extend, createFFListActions } from '@tencent/ff-redux';
import { RootState, Project, ChartInfoFilter, ProjectFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { projectNamespaceActions } from '../namespace';
import { router } from '../../router';
import { setProjectName } from '../../../../../helpers';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchProjectActions = createFFListActions<Project, ProjectFilter, ChartInfoFilter>({
  actionName: ActionTypes.ProjectList,
  fetcher: async (query, getState: GetState) => {
    // let response = await WebAPI.fetchManagedProjectList(query);
    let response = await WebAPI.fetchPortalProjectList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().projectList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      dispatch(listActions.selectProject(record.data.records[0].metadata.name, record.data.data));
    }
  }
});

const restActions = {
  selectProject: (projectId: string, chartInfoFilter?: ChartInfoFilter) => {
    return async (dispatch, getState: GetState) => {
      setProjectName(projectId);

      dispatch(listActions.selectByValue(projectId));
      dispatch(
        projectNamespaceActions.list.applyFilter({
          projectId: projectId,
          chartInfoFilter: chartInfoFilter
        })
      );
    };
  }
};

export const listActions = extend({}, fetchProjectActions, restActions);
