import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { initValidator } from '../../common/models/Validation';

const TempReducer = combineReducers({
  alarmPolicyId: reduceToPayload(ActionType.InputAlarmPolicyId, ''),

  alarmPolicyName: reduceToPayload(ActionType.InputAlarmPolicyName, ''),

  v_alarmPolicyName: reduceToPayload(ActionType.ValidateAlarmPolicyName, initValidator),

  /**策略备注 */
  alarmPolicyDescription: reduceToPayload(ActionType.InputAlarmPolicyDescription, ''),
  v_alarmPolicyDescription: reduceToPayload(ActionType.ValidateAlarmPolicyDescription, initValidator),

  /**策略对象类型 */
  alarmPolicyType: reduceToPayload(ActionType.InputAlarmPolicyType, 'cluster'),

  v_alarmPolicyType: reduceToPayload(ActionType.ValidateAlarmPolicyType, initValidator),

  statisticsPeriod: reduceToPayload(ActionType.InputAlarmPolicyStatisticsPeriod, '1'),
  alarmPolicyChannel: reduceToPayload(ActionType.InputAlarmPolicyChannel, ''),
  alarmPolicyTemplate: reduceToPayload(ActionType.InputAlarmPolicyTemplate, ''),

  /**策略对象数组 */
  alarmObjects: reduceToPayload(ActionType.InputAlarmPolicyObjects, []),

  v_alarmObjects: reduceToPayload(ActionType.ValidateAlarmPolicyObjects, initValidator),

  alarmObjectsType: reduceToPayload(ActionType.InputAlarmPolicyObjectsType, 'all'),

  alarmObjectNamespace: reduceToPayload(ActionType.InputAlarmWorkLoadNameSpace, ''),

  alarmObjectWorkloadType: reduceToPayload(ActionType.InputAlarmObjectWorkloadType, 'deployment'),

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
