import { extend, generateWorkflowActionCreator, OperationResult, OperationTrigger, createFFListActions } from '@tencent/ff-redux';
import { RootState, Role, RoleFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchRoleActions = createFFListActions<Role, RoleFilter>({
  actionName: ActionTypes.RoleList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchRoleList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().roleList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {
    let isNotNeedPoll = true;
    if (record.data.recordCount) {
      isNotNeedPoll = record.data.records.filter(item => item.status && item.status['phase'] && item.status['phase'] === 'Terminating').length === 0;
    }
    if (isNotNeedPoll) {
      dispatch(fetchRoleActions.clearPolling());
    }
  }
});

/**
 * 删除角色
 */
const removeRoleWorkflow = generateWorkflowActionCreator<Role, void>({
  actionType: ActionTypes.RemoveRole,
  workflowStateLocator: (state: RootState) => state.roleRemoveWorkflow,
  operationExecutor: WebAPI.deleteRole,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      //在列表页删除的动作，因此直接重新拉取一次数据
      dispatch(listActions.poll());
      /** 结束工作流 */
      dispatch(listActions.removeRoleWorkflow.reset());
    }
  }
});

const restActions = {
  removeRoleWorkflow,

  /** 轮询操作 */
  poll: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(
        listActions.polling({
            delayTime: 3000
          }
        )
      );
    };
  },

};

export const listActions = extend({}, fetchRoleActions, restActions);
