import { extend, createFFListActions } from '@tencent/ff-redux';
import { RootState, Project } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { router } from '../../router';
import { projectNamespaceActions } from '../namespace';
import { setProjectName } from '../../../../../helpers';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchProjectActions = createFFListActions<Project, void>({
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
      let { route } = getState();
      const projectId = route.queries['projectId'];
      let exist = record.data.records.find(r => {
        return r.metadata.name === projectId;
      });
      if (exist) {
        dispatch(listActions.selectProject(exist.metadata.name));
      } else {
        dispatch(listActions.selectProject(record.data.records[0].metadata.name));
      }
    }
  }
});

const restActions = {
  selectProject: (projectId: string) => {
    return async (dispatch, getState: GetState) => {
      setProjectName(projectId);

      let { route } = getState(),
        urlParams = router.resolve(route);
      router.navigate(urlParams, Object.assign({}, route.queries, { projectId: projectId }));

      dispatch(listActions.selectByValue(projectId));
      dispatch(
        projectNamespaceActions.list.applyFilter({
          projectId: projectId
        })
      );
    };
  }
};

export const listActions = extend({}, fetchProjectActions, restActions);
