import {
  ReduxAction,
  extend,
  generateWorkflowActionCreator,
  OperationTrigger,
  isSuccessWorkflow,
  getWorkflowStatistics,
  createFFObjectActions
} from '@tencent/ff-redux';
import { RootState, AppFilter, AppEditor, App, AppDetailFilter, ChartInfo, ChartInfoFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { initAppEditorState } from '../../constants/initState';
import { AppValidateSchema } from '../../constants/AppValidateConfig';
import { router } from '../../router';
import { createValidatorActions, getValidatorActionType } from '@tencent/ff-validator';
type GetState = () => RootState;

/**
 * 修改应用
 */
const updateAppWorkflow = generateWorkflowActionCreator<App, void>({
  actionType: ActionTypes.UpdateApp,
  workflowStateLocator: (state: RootState) => state.appUpdateWorkflow,
  operationExecutor: WebAPI.updateApp,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      let { appUpdateWorkflow, appEditor, route } = getState();
      if (isSuccessWorkflow(appUpdateWorkflow)) {
        //判断是否是dryrun
        if (appUpdateWorkflow.targets.length > 0 && appUpdateWorkflow.targets[0].spec.dryRun) {
          if (
            appUpdateWorkflow.results &&
            appUpdateWorkflow.results.length > 0 &&
            appUpdateWorkflow.results[0].success &&
            appUpdateWorkflow.results[0].target
          ) {
            let { appDryRun } = getState();
            dispatch({
              type: ActionTypes.UpdateAppDryRunState,
              payload: Object.assign({}, appDryRun, appUpdateWorkflow.results[0].target)
            });
          }
        } else {
          //表示编辑模式结束
          dispatch({
            type: ActionTypes.UpdateAppEditorState,
            payload: Object.assign({}, appEditor, { v_editing: false })
          });
          /** 重新获取最新数据，从而Detail可以连续编辑且使用到最新的resourceVersion */
          // dispatch(
          //   detailActions.fetchApp({
          //     cluster: appEditor.spec.targetCluster,
          //     namespace: appEditor.metadata.namespace,
          //     name: appEditor.metadata.name
          //   })
          // );

          router.navigate({ mode: '', sub: 'app' }, route.queries);
        }
      }
      /** 结束工作流 */
      dispatch(detailActions.updateAppWorkflow.reset());
    }
  }
});

/**
 * 获取chart详情
 */
const fetchChartInfoActions = createFFObjectActions<ChartInfo, ChartInfoFilter>({
  actionName: ActionTypes.ChartInfo,
  fetcher: async query => {
    let response = await WebAPI.fetchChartInfo(query.filter);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chartInfo;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { appEditor } = getState();
    let values = Object.assign({}, appEditor.spec.values);
    if (record.data && record.data.spec.values && record.data.spec.values['values.yaml']) {
      values.rawValues = record.data.spec.values['values.yaml'];
    } else {
      values.rawValues = '';
    }
    dispatch(detailActions.updateEditorState({ spec: Object.assign({}, appEditor.spec, { values: values }) }));
  }
});

const restActions = {
  updateAppWorkflow,

  validator: createValidatorActions({
    userDefinedSchema: AppValidateSchema,
    validateStateLocator: (store: RootState) => {
      return store.appEditor;
    },
    validatorStateLocation: (store: RootState) => {
      return store.appValidator;
    }
  }),

  fetchApp: (filter: AppDetailFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let response = await WebAPI.fetchApp(filter);
      let editor: AppEditor = response;
      if (editor.spec.values && editor.spec.values.rawValues) {
        try {
          editor.spec.values.rawValues = window.atob(editor.spec.values.rawValues);
        } catch (error) {
          console.log(error);
        }
      }
      dispatch({
        type: ActionTypes.UpdateAppEditorState,
        payload: editor
      });
    };
  },

  /** 更新状态 */
  updateEditorState: obj => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { appEditor } = getState();
      dispatch({
        type: ActionTypes.UpdateAppEditorState,
        payload: Object.assign({}, appEditor, obj)
      });
    };
  },

  /** 离开更新页面，清除Editor当中的内容 */
  clearEditorState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateAppEditorState,
      payload: initAppEditorState
    };
  },

  /** 离开创建页面，清除Validator当中的内容 */
  clearValidatorState: (): ReduxAction<any> => {
    return {
      type: getValidatorActionType(AppValidateSchema.formKey),
      payload: {}
    };
  },

  /** 清除DryRun当中的内容 */
  clearDryRunState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateAppDryRunState,
      payload: {}
    };
  }
};
export const detailActions = extend({}, { chart: fetchChartInfoActions }, restActions);
