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

import { combineReducers } from 'redux';

import { reduceToPayload } from '@tencent/ff-redux';

import { initValidator } from '../../common/models/Validation';
import * as ActionType from '../constants/ActionType';

let initAlarmPolicyType = 'cluster';
let initAlarmObjectsType = 'all';
/// #if project
initAlarmPolicyType = 'pod';
initAlarmObjectsType = 'part';
/// #endif
const TempReducer = combineReducers({
  alarmPolicyId: reduceToPayload(ActionType.InputAlarmPolicyId, ''),

  alarmPolicyName: reduceToPayload(ActionType.InputAlarmPolicyName, ''),

  v_alarmPolicyName: reduceToPayload(ActionType.ValidateAlarmPolicyName, initValidator),

  /**策略备注 */
  alarmPolicyDescription: reduceToPayload(ActionType.InputAlarmPolicyDescription, ''),
  v_alarmPolicyDescription: reduceToPayload(ActionType.ValidateAlarmPolicyDescription, initValidator),

  /**策略对象类型 */
  alarmPolicyType: reduceToPayload(ActionType.InputAlarmPolicyType, initAlarmPolicyType),

  v_alarmPolicyType: reduceToPayload(ActionType.ValidateAlarmPolicyType, initValidator),

  statisticsPeriod: reduceToPayload(ActionType.InputAlarmPolicyStatisticsPeriod, '1'),
  alarmPolicyChannel: reduceToPayload(ActionType.InputAlarmPolicyChannel, ''),
  alarmPolicyTemplate: reduceToPayload(ActionType.InputAlarmPolicyTemplate, ''),

  /**策略对象数组 */
  alarmObjects: reduceToPayload(ActionType.InputAlarmPolicyObjects, []),

  v_alarmObjects: reduceToPayload(ActionType.ValidateAlarmPolicyObjects, initValidator),

  alarmObjectsType: reduceToPayload(ActionType.InputAlarmPolicyObjectsType, initAlarmObjectsType),

  alarmObjectNamespace: reduceToPayload(ActionType.InputAlarmWorkLoadNameSpace, ''),

  alarmObjectWorkloadType: reduceToPayload(ActionType.InputAlarmObjectWorkloadType, 'Deployment'),

  /**策略指标 */
  alarmMetrics: reduceToPayload(ActionType.InputAlarmMetrics, []),
  v_alarmMetrics: reduceToPayload(ActionType.ValidateAlarmMetrics, initValidator),

  /**是否有生效时间 */

  enableShield: reduceToPayload(ActionType.InputAlarmenAbleShield, false),

  /**策略生效时间 */
  shieldTimeStart: reduceToPayload(ActionType.InputAlarmShieldTimeStart, '00:00:00'),

  shieldTimeEnd: reduceToPayload(ActionType.InputAlarmShieldTimeEnd, '00:00:01'),

  v_policyTime: reduceToPayload(ActionType.ValidateAlarmpolicyTime, initValidator),

  /**策略告警通道方式  eg：SMS，EMAIL..*/
  notifyWays: reduceToPayload(ActionType.InputAlarmNotifyWay, [{ id: 0, channel: undefined, template: undefined }]),
  v_notifyWay: reduceToPayload(ActionType.ValidatetAlarmNotifyWay, initValidator),
  /**手机告警设置 */
  phoneNotifyOrder: reduceToPayload(ActionType.InputAlarmPhoneNotifyOrder, []),
  phoneCircleTimes: reduceToPayload(ActionType.InputAlarmPhoneCircleTimes, 3),
  phoneInnerInterval: reduceToPayload(ActionType.InputAlarmPhoneInnerInterval, 3),
  phoneCircleInterval: reduceToPayload(ActionType.InputAlarmPhoneCircleInterval, 3),
  phoneArriveNotice: reduceToPayload(ActionType.InputAlarmPhoneArriveNotice, false),

  // 告警组选择
  groupSelection: reduceToPayload(ActionType.SelectGroup, []),
  v_groupSelection: reduceToPayload(ActionType.ValidateSelectGroup, initValidator)
});

export const AlarmPolicyEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建Ingress页面
  if (action.type === ActionType.ClearAlarmPolicyEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
