import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import {
  createFFListActions,
  createFFObjectActions,
  deepClone,
  extend,
  generateWorkflowActionCreator,
  isSuccessWorkflow,
  OperationTrigger,
  uuid
} from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { initValidator } from '../../common/models/Validation';
import * as ActionType from '../constants/ActionType';
import {
  FFReduxActionName,
  initProjectEdition,
  initProjectResourceLimit,
  resourceTypeToUnit
} from '../constants/Config';
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
