
import { ReduxAction, extend, generateWorkflowActionCreator, OperationTrigger, isSuccessWorkflow } from '@tencent/ff-redux';
import { ChartGroup, RootState } from '../../models/index';
import * as ActionTypes from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { initChartGroupCreationState } from '../../constants/initState';
import { ChartGroupValidateSchema } from '../../constants/ChartGroupValidateConfig';
import { router } from '../../router';
import { createValidatorActions, getValidatorActionType } from '@tencent/ff-validator';
import { listActions } from './listActions';
type GetState = () => RootState;

/**
 * 增加仓库
 */
const addChartGroupWorkflow = generateWorkflowActionCreator<ChartGroup, void>({
  actionType: ActionTypes.AddChartGroup,
  workflowStateLocator: (state: RootState) => state.chartGroupAddWorkflow,
  operationExecutor: WebAPI.addChartGroup,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      let { chartGroupAddWorkflow, route } = getState();
      if (isSuccessWorkflow(chartGroupAddWorkflow)) {
        router.navigate({ mode: '', sub: 'chartgroup' }, route.queries);
        //进入列表时自动加载
        //退出状态页面时自动清理状态
      }
      /** 结束工作流 */
      dispatch(createActions.addChartGroupWorkflow.reset());
    }
  }
});

const restActions = {
  addChartGroupWorkflow,

  validator: createValidatorActions({
    userDefinedSchema: ChartGroupValidateSchema,
    validateStateLocator: (store: RootState) => {
      return store.chartGroupCreation;
    },
    validatorStateLocation: (store: RootState) => {
      return store.chartGroupValidator;
    },
    // used in extraStore, i.t. customFunc: (value, store, extraStore)
    extraValidateStateLocatorPath: ['userInfo']
  }),

  /** 更新状态 */
  updateCreationState: (obj) => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { chartGroupCreation } = getState();
      dispatch({
        type: ActionTypes.UpdateChartGroupCreationState,
        payload: Object.assign({}, chartGroupCreation, obj)
      });
    };
  },

  /** 离开创建页面，清除Creation当中的内容 */
  clearCreationState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateChartGroupCreationState,
      payload: initChartGroupCreationState
    };
  },

  /** 离开创建页面，清除Validator当中的内容 */
  clearValidatorState: (): ReduxAction<any> => {
    return {
      type: getValidatorActionType(ChartGroupValidateSchema.formKey),
      payload: {}
    };
  }
};

export const createActions = extend({}, restActions);