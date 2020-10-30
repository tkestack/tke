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
