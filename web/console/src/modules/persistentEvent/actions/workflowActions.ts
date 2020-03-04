import {
    generateWorkflowActionCreator, isSuccessWorkflow, OperationTrigger
} from '@tencent/ff-redux';

import { CommonAPI } from '../../common/webapi';
import * as ActionType from '../constants/ActionType';
import { CreateResource, RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

/**
 * 操作流Action
 */
export const workflowActions = {
  /** 设置持久化事件 */
  modifyPeFlow: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.ModifyPeFlow,
    workflowStateLocator: (state: RootState) => state.modifyPeFlow,
    operationExecutor: WebAPI.modifyPeConfig,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { route, modifyPeFlow } = getState(),
          urlParams = router.resolve(route);

        if (isSuccessWorkflow(modifyPeFlow)) {
          // 跳转到 持久化事件的列表页
          let routeQueries = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { clusterId: undefined })));
          let newUrlParams = JSON.parse(JSON.stringify(Object.assign({}, urlParams, { mode: undefined })));
          router.navigate(newUrlParams, routeQueries);
          // dispatch(clusterActions.applyFilter({ regionId: route.queries['rid'] }));
          // dispatch(peActions.applyFilter({ regionId: route.queries['rid'] }));
          dispatch(workflowActions.modifyPeFlow.reset());
        }
      }
    }
  }),

  /** 删除PersistentEvent */
  deletePeFlow: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.DeletePeFlow,
    workflowStateLocator: (state: RootState) => state.deletePeFlow,
    operationExecutor: CommonAPI.deleteResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { deletePeFlow, route } = getState(),
          urlParams = router.resolve(route);

        if (isSuccessWorkflow(deletePeFlow)) {
          // 跳转到 持久化事件的列表页
          let routeQueries = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { clusterId: undefined })));
          let newUrlParams = JSON.parse(JSON.stringify(Object.assign({}, urlParams, { mode: undefined })));
          router.navigate(newUrlParams, routeQueries);
          dispatch(workflowActions.deletePeFlow.reset());
        }
      }
    }
  })
};
