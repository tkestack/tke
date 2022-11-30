/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { deepClone } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { AlarmPolicyEdition } from '../models/AlarmPolicy';

type GetState = () => RootState;

export const validatorActions = {
  /**
   * 校验集群名称
   */
  _validateAlarmPolicyName(name) {
    let status = 0,
      message = '';

    //验证集群名称
    if (!name) {
      status = 2;
      message = t('告警策略名称不能为空');
    } else if (name.length > 60) {
      status = 2;
      message = t('告警策略名称不能超过60个字符');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateAlarmPolicyName() {
    return async (dispatch, getState: GetState) => {
      const name = getState().alarmPolicyEdition.alarmPolicyName;
      const result = validatorActions._validateAlarmPolicyName(name);
      dispatch({
        type: ActionType.ValidateAlarmPolicyName,
        payload: result
      });
    };
  },
  /**
   * 校验集群描述
   */
  _validateDescription(name) {
    let status = 0,
      message = '';

    //验证告警设置描述
    if (name && name.length > 100) {
      status = 2;
      message = t('告警设置描述不能超过100字符');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateDescription() {
    return async (dispatch, getState: GetState) => {
      const name = getState().alarmPolicyEdition.alarmPolicyDescription;
      const result = validatorActions._validateDescription(name);

      dispatch({
        type: ActionType.ValidateAlarmPolicyDescription,
        payload: result
      });
    };
  },
  validatePolicyTime() {
    return async (dispatch, getState: GetState) => {
      const start = getState().alarmPolicyEdition.shieldTimeStart,
        end = getState().alarmPolicyEdition.shieldTimeEnd;
      const result = validatorActions._validatePolicyTime(start, end);

      dispatch({
        type: ActionType.ValidateAlarmpolicyTime,
        payload: result
      });
    };
  },
  _validatePolicyTime(start, end) {
    let status = 0,
      message = '';
    if (Date.parse('1970-01-01 ' + end) - Date.parse('1970-01-01 ' + start) <= 0) {
      status = 2;
      message = t('告警设置有效时间结束时间不能早于开始时间');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  validateGroupSelection() {
    return async (dispatch, getState: GetState) => {
      const group = getState().receiverGroup.selections;

      const result = validatorActions._validateGroupSelection(group);

      dispatch({
        type: ActionType.ValidateSelectGroup,
        payload: result
      });
    };
  },
  _validateGroupSelection(group) {
    let status = 0,
      message = '';

    if (group.length === 0) {
      status = 2;
      message = t('接收组不能为空');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  validateAlarmPolicyType() {
    return async (dispatch, getState: GetState) => {
      const type = getState().alarmPolicyEdition.alarmPolicyType;

      const result = validatorActions._validateAlarmPolicyType(type);

      dispatch({
        type: ActionType.ValidateAlarmPolicyType,
        payload: result
      });
    };
  },
  _validateAlarmPolicyType(group) {
    let status = 0,
      message = '';

    if (!group) {
      status = 2;
      message = t('策略类型不能为空');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  _validateEvaluatorValue(value: string, canOver100 = false) {
    let status = 0,
      message = '';
    if (!value) {
      status = 2;
      message = t('阈值不能为空');
    } else if (isNaN(+value)) {
      status = 2;
      message = t('阈值的格式错误');
    } else if ((+value > 100 && !canOver100) || +value <= 0) {
      status = 2;
      message = t('阈值的范围错误');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  _validateEvaluatorSumValue(value: string) {
    let reg = /^\d+(\.\d{1,3})?$/,
      status = 0,
      message = '';
    if (!value) {
      status = 2;
      message = t('阈值不能为空');
    } else if (!reg.test(value)) {
      status = 2;
      message = t('数据格式不正确，使用量限制只能是小数，且只能精确到0.01');
    } else if (+value < 0.01) {
      status = 2;
      message = t('使用量限制最小为0.01');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  validateEvaluatorValue(id: string) {
    return async (dispatch, getState: GetState) => {
      const { alarmPolicyEdition } = getState(),
        { alarmMetrics } = alarmPolicyEdition;
      const newAlarmMetrics = deepClone(alarmMetrics),
        index = newAlarmMetrics.findIndex(e => e.id === id);
      if (newAlarmMetrics[index].type === 'sumCpu' || newAlarmMetrics[index].type === 'sumMem') {
        newAlarmMetrics[index].v_evaluatorValue = validatorActions._validateEvaluatorSumValue(
          newAlarmMetrics[index].evaluatorValue
        );
      } else {
        newAlarmMetrics[index].v_evaluatorValue = validatorActions._validateEvaluatorValue(
          newAlarmMetrics[index].evaluatorValue,
          newAlarmMetrics?.[index]?.metricName === 'vm_cpu_usage_rate'
        );
      }

      dispatch({
        type: ActionType.InputAlarmMetrics,
        payload: newAlarmMetrics
      });
    };
  },

  _validateAlarmPolicyEdition(alarmPolicyEdition: AlarmPolicyEdition, receiverGroup) {
    let isOk = true;
    isOk =
      isOk &&
      validatorActions._validateAlarmPolicyName(alarmPolicyEdition.alarmPolicyName).status === 1 &&
      validatorActions._validateDescription(alarmPolicyEdition.alarmPolicyDescription).status === 1 &&
      validatorActions._validateAlarmPolicyType(alarmPolicyEdition.alarmPolicyType).status === 1 &&
      validatorActions._validateGroupSelection(receiverGroup.selections).status === 1;
    alarmPolicyEdition.alarmMetrics.forEach(item => {
      if (item.enable) {
        if (item.type === 'sumCpu' || item.type === 'sumMem') {
          isOk = isOk && validatorActions._validateEvaluatorSumValue(item.evaluatorValue + '').status === 1;
        } else if (item.type === 'times' || item.type === 'percent') {
          isOk =
            isOk &&
            validatorActions._validateEvaluatorValue(item.evaluatorValue + '', item?.metricName === 'vm_cpu_usage_rate')
              .status === 1;
        }
      }
    });
    // isOk = isOk && alarmPolicyEdition.notifyWay.length !== 0;
    isOk = isOk && alarmPolicyEdition.alarmMetrics.some(item => item.enable);
    if (alarmPolicyEdition.alarmPolicyType === 'pod' && alarmPolicyEdition.alarmObjectsType === 'part') {
      isOk = isOk && alarmPolicyEdition.alarmObjects.length !== 0;
    }
    return isOk;
  },

  validateAlarmPolicyEdition() {
    return async (dispatch, getState: GetState) => {
      const { alarmPolicyEdition } = getState(),
        { alarmMetrics } = alarmPolicyEdition;
      dispatch(validatorActions.validateAlarmPolicyType());
      dispatch(validatorActions.validateAlarmPolicyName());
      dispatch(validatorActions.validateDescription());
      dispatch(validatorActions.validateGroupSelection());
      alarmMetrics.forEach(item => {
        if (item.enable) {
          dispatch(validatorActions.validateEvaluatorValue(item.id + ''));
        }
      });
      //当选择pod类型，按工作负载类型时校验告警对象
      if (alarmPolicyEdition.alarmPolicyType === 'pod' && alarmPolicyEdition.alarmObjectsType === 'part') {
        if (alarmPolicyEdition.alarmObjects.length === 0) {
          dispatch({
            type: ActionType.ValidateAlarmPolicyObjects,
            payload: {
              message: t('告警对象不能为空'),
              status: 2
            }
          });
        } else {
          dispatch({
            type: ActionType.ValidateAlarmPolicyObjects,
            payload: {
              message: '',
              status: 1
            }
          });
        }
      }
      //校验指标是否为空
      if (alarmPolicyEdition.alarmMetrics.some(item => item.enable) === false) {
        dispatch({
          type: ActionType.ValidateAlarmMetrics,
          payload: {
            message: t('告警设置指标不能为空'),
            status: 2
          }
        });
      } else {
        dispatch({
          type: ActionType.ValidateAlarmMetrics,
          payload: {
            message: '',
            status: 1
          }
        });
      }
      /**校验告警渠道是否为空 */
      if (alarmPolicyEdition.notifyWays.length === 0) {
        dispatch({
          type: ActionType.ValidatetAlarmNotifyWay,
          payload: {
            message: t('通知方式不能为空'),
            status: 2
          }
        });
      } else if (alarmPolicyEdition.notifyWays.find(({ channel, template }) => !channel || !template)) {
        dispatch({
          type: ActionType.ValidatetAlarmNotifyWay,
          payload: {
            message: t('渠道与模板不能为空'),
            status: 2
          }
        });
      } else {
        dispatch({
          type: ActionType.ValidatetAlarmNotifyWay,
          payload: {
            message: '',
            status: 1
          }
        });
      }
    };
  }
};
