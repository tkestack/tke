import { AlarmPolicyEdition, AlarmPolicyOperator, AlarmPolicy } from './../models/AlarmPolicy';
import { RecordSet, extend, ReduxAction } from '@tencent/qcloud-lib';
import { generateWorkflowActionCreator, OperationHooks, OperationTrigger, isSuccessWorkflow } from '@tencent/ff-redux';
import { router } from '../router';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../WebAPI';
import { RootState } from '../models';
import { alarmPolicyActions } from './alarmPolicyActions';

type GetState = () => RootState;

/**
 * 操作流Action
 */
export const workflowActions = {
  editAlarmPolicy: generateWorkflowActionCreator<AlarmPolicyEdition, AlarmPolicyOperator>({
    actionType: ActionType.CreateAlarmPolicy,
    workflowStateLocator: (state: RootState) => state.alarmPolicyCreateWorkflow,
    operationExecutor: (
      targets: AlarmPolicyEdition[],
      params: AlarmPolicyOperator,
      dispatch: Redux.Dispatch,
      getState: GetState
    ) => {
      return WebAPI.editAlarmPolicy(targets, params, getState().receiverGroup);
    },
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        // 清空更改名称的表单
        let { alarmPolicyCreateWorkflow, route } = getState();
        if (isSuccessWorkflow(alarmPolicyCreateWorkflow)) {
          dispatch(workflowActions.editAlarmPolicy.reset());
          dispatch(alarmPolicyActions.clearAlarmPolicyEdit());
          router.navigate({}, { rid: route.queries['rid'], clusterId: route.queries['clusterId'] });
        }
      }
    }
  }),
  deleteAlarmPolicy: generateWorkflowActionCreator<AlarmPolicy, AlarmPolicyOperator>({
    actionType: ActionType.DeleteAlarmPolicy,
    workflowStateLocator: (state: RootState) => state.alarmPolicyDeleteWorkflow,
    operationExecutor: WebAPI.deleteAlarmPolicy,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        // 清空更改名称的表单
        let { alarmPolicyDeleteWorkflow, route } = getState();
        if (isSuccessWorkflow(alarmPolicyDeleteWorkflow)) {
          dispatch(workflowActions.deleteAlarmPolicy.reset());
          dispatch(alarmPolicyActions.fetch());
          dispatch(alarmPolicyActions.selects([]));
        }
      }
    }
  })
};
