/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import {
  ReduxAction,
  extend,
  generateWorkflowActionCreator,
  OperationTrigger,
  isSuccessWorkflow,
  createFFObjectActions
} from '@tencent/ff-redux';
import { App, AppCreation, RootState, ChartInfo, ChartInfoFilter } from '../../models/index';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { initAppCreationState } from '../../constants/initState';
import { AppValidateSchema } from '../../constants/AppValidateConfig';
import { router } from '../../router';
import { createValidatorActions, getValidatorActionType } from '@tencent/ff-validator';
import { listActions } from './listActions';
type GetState = () => RootState;

/**
 * 增加仓库
 */
const addAppWorkflow = generateWorkflowActionCreator<App, void>({
  actionType: ActionTypes.AddApp,
  workflowStateLocator: (state: RootState) => state.appAddWorkflow,
  operationExecutor: WebAPI.addApp,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      let { appAddWorkflow, route } = getState();
      if (isSuccessWorkflow(appAddWorkflow)) {
        //判断是否是dryrun
        if (appAddWorkflow.targets.length > 0 && appAddWorkflow.targets[0].spec.dryRun) {
          if (
            appAddWorkflow.results &&
            appAddWorkflow.results.length > 0 &&
            appAddWorkflow.results[0].success &&
            appAddWorkflow.results[0].target
          ) {
            let { appDryRun } = getState();
            dispatch({
              type: ActionTypes.UpdateAppDryRunState,
              payload: Object.assign({}, appDryRun, appAddWorkflow.results[0].target)
            });
          }
        } else {
          router.navigate({ mode: '', sub: 'app' }, route.queries);
          //进入列表时自动加载
          //退出状态页面时自动清理状态
        }
      }
      /** 结束工作流 */
      dispatch(createActions.addAppWorkflow.reset());
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
    let { appCreation } = getState();
    let values = Object.assign({}, appCreation.spec.values);
    if (record.data && record.data.spec.values && record.data.spec.values['values.yaml']) {
      values.rawValues = record.data.spec.values['values.yaml'];
    } else {
      values.rawValues = '';
    }
    dispatch(createActions.updateCreationState({ spec: Object.assign({}, appCreation.spec, { values: values }) }));
  }
});

const restActions = {
  addAppWorkflow,

  validator: createValidatorActions({
    userDefinedSchema: AppValidateSchema,
    validateStateLocator: (store: RootState) => {
      return store.appCreation;
    },
    validatorStateLocation: (store: RootState) => {
      return store.appValidator;
    }
  }),

  /** 更新状态 */
  updateCreationState: obj => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { appCreation } = getState();
      dispatch({
        type: ActionTypes.UpdateAppCreationState,
        payload: Object.assign({}, appCreation, obj)
      });
    };
  },

  /** 离开创建页面，清除Creation当中的内容 */
  clearCreationState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateAppCreationState,
      payload: initAppCreationState
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

export const createActions = extend({}, { chart: fetchChartInfoActions }, restActions);
