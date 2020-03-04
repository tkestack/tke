import { extend, Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models/Validation';
import { Resource } from '../../notify/models';
import { Group } from './Group';

export interface AlarmPolicy extends Identifiable {
  alarmPolicyId: string;
  clusterId: string;
  alarmPolicyName: string;
  alarmPolicyDescription: string;
  alarmPolicyType: string;
  statisticsPeriod: number;
  alarmMetrics: MetricsObject[];
  alarmObjetcs: string[];
  alarmObjetcsType: string;
  alarmObjectWorkloadType?: string;
  alarmObjectNamespace?: string;
  enableShield: boolean;
  shieldTimeStart: string;
  shieldTimeEnd: string;
  receiverGroups: string[];
  notifyWays: { id: string; channel: string; template: string }[];

  // notifyWay: string[];
  // phoneNotifyOrder: number[];
  // phoneCircleTimes: number;
  // phoneInnerInterval: number;
  // phoneCircleInterval: number;
  // phoneArriveNotice: number;
}

export interface AlarmPolicyEdition extends Identifiable {
  alarmPolicyId?: string;
  /**策略名称 */
  alarmPolicyName?: string;
  v_alarmPolicyName?: Validation;

  /**策略备注 */
  alarmPolicyDescription?: string;
  v_alarmPolicyDescription?: Validation;

  /**策略对象类型 */
  alarmPolicyType?: string;

  statisticsPeriod?: number;

  v_alarmPolicyType?: Validation;

  /**策略对象*/
  alarmObjects?: string[];
  v_alarmObjects?: Validation;

  /**策略对象类型*/
  alarmObjectsType?: string;

  /**策略对象 Pod设置*/
  alarmObjectNamespace?: string;

  alarmObjectWorkloadType?: string;

  /**策略指标 */
  alarmMetrics?: MetricsObjectEdition[];
  v_alarmMetrics?: Validation;

  /**是否有生效时间 */

  enableShield: boolean;

  /**策略生效时间 */
  shieldTimeStart: string;

  shieldTimeEnd: string;

  v_policyTime?: Validation;

  /**策略告警通道方式  eg：SMS，EMAIL..*/
  // notifyWay: string[];
  // v_notifyWay?: Validation;

  notifyWays: { channel: string; template: string; id: string }[];

  // alarmPolicyChannel: string;

  // alarmPolicyTemplate: string;

  /**手机告警设置 */
  // phoneNotifyOrder: number[];
  // phoneCircleTimes: number;
  // phoneInnerInterval: number;
  // phoneCircleInterval: number;
  // phoneArriveNotice: boolean;

  // 告警组选择
  // groupSelection?: Resource[];
  v_groupSelection?: Validation;
}

export interface AlarmPolicyFilter {
  /**
   * 根据集群Id进行筛选
   */
  clusterId?: string;

  /**
   * 地域
   */
  regionId?: number;
}

export interface AlarmPolicyOperator {
  /**
   * 集群Id
   */
  clusterId?: string;

  /**
   * 地域
   */
  regionId?: number;
}

export interface MetricsObject {
  /**策略ID */
  metricId?: string;
  /*指标名称*/
  measurement: string;
  /**统计周期 */
  // statisticsPeriod: number;

  /**指标类型 */
  metricName: string;

  /**指标描述 */
  metricDisplayName: string;

  /**指标表达式 string里面放一个json eg:{\"type\":\"gt\",\"value\": 95}*/
  evaluatorType: string;

  evaluatorValue: string;

  /**持续周期 */
  continuePeriod: number;

  status?: boolean;

  //指标类型
  type: string;

  //指标提示
  tip: string;

  //单位
  unit: string;
}

export interface MetricsObjectEdition extends MetricsObject, Identifiable {
  /**是否启用该指标 */
  enable?: boolean;

  v_evaluatorValue: Validation;
}
