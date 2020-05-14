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
