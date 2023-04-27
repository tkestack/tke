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
import { createFFListActions, deepClone, extend, uuid } from '@tencent/ff-redux';

import { initValidator } from '../../common/models';
import { resourceActions } from '../../notify/actions/resourceActions';
import { userActions } from '../../uam/actions/userActions';
import * as ActionType from '../constants/ActionType';
import { AlarmPolicyMetrics, AlarmPolicyType } from '../constants/Config';
import { AlarmPolicy, AlarmPolicyFilter, RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { namespaceActions } from './namespaceActions';
import { validatorActions } from './validatorActions';
import { workloadActions } from './workloadActions';
import { reverseReduceNs } from '@helper/urlUtil';

type GetState = () => RootState;

const _alarmPolicyActions = createFFListActions<AlarmPolicy, AlarmPolicyFilter>({
  actionName: 'AlarmPolicy',
  fetcher: async (query, getstate: GetState) => {
    const response = await WebAPI.fetchAlarmPolicy(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().alarmPolicy;
  },
  onFinish: (record, dispatch, getState: GetState) => {
    const { sub } = router.resolve(getState().route);
    if (sub !== '') {
      dispatch(alarmPolicyActions.initAlarmPolicyData());
    }
  }
});

/**
 * 选择告警设置
 */
const restActions = {
  fetchAlarmPolicyDetail: alarmpolicy => {
    return (dispatch, getState) => {
      dispatch({
        type: ActionType.FetchalarmPolicyDetail,
        payload: alarmpolicy
      });
    };
  }
};

const editActions = {
  initAlarmPolicyEditionForCopy: (alarmpolicy: AlarmPolicy) => {
    return (dispatch, getState: GetState) => {
      dispatch(alarmPolicyActions.inputAlarmPolicyName(alarmpolicy.alarmPolicyName));
      dispatch(alarmPolicyActions.inputAlarmPolicyDescription(alarmpolicy.alarmPolicyDescription || ''));
      dispatch({
        type: ActionType.InputAlarmPolicyType,
        payload: alarmpolicy.alarmPolicyType
      });
      if (alarmpolicy.alarmPolicyType === 'pod' || alarmpolicy.alarmPolicyType === 'virtualMachine') {
        dispatch({
          type: ActionType.InputAlarmPolicyObjectsType,
          payload: 'part'
        });
      }

      dispatch(alarmPolicyActions.initAlarmMetricsForUpdate(alarmpolicy.alarmMetrics, alarmpolicy.alarmPolicyType));
      dispatch(
        alarmPolicyActions.inputAlarmNotifyWays(
          alarmpolicy.notifyWays.map(w => ({
            id: uuid(),
            channel: w.channel,
            template: w.template
          }))
        )
      );
    };
  },

  initAlarmPolicyEditionForUpdate: (alarmpolicy: AlarmPolicy) => {
    return (dispatch, getState: GetState) => {
      //初始化一部分,复用
      dispatch(alarmPolicyActions.initAlarmPolicyEditionForCopy(alarmpolicy));
      //将id赋值给edition
      dispatch({
        type: ActionType.InputAlarmPolicyId,
        payload: alarmpolicy.alarmPolicyId
      });

      if (alarmpolicy.alarmPolicyType !== 'cluster') {
        dispatch(alarmPolicyActions.inputAlarmPolicyObjectsType(alarmpolicy.alarmObjetcsType));
        //告警对象是workload且选择按工作负载选择初始化

        if (alarmpolicy.alarmPolicyType === 'pod') {
          dispatch(alarmPolicyActions.inputAlarmObjectWorkloadType(alarmpolicy.alarmObjectWorkloadType));
        }

        if (alarmpolicy.alarmPolicyType === 'pod' || alarmpolicy.alarmPolicyType === 'virtualMachine') {
          const hasNamespace = alarmpolicy?.alarmObjetcsType === 'part' || alarmpolicy?.alarmObjectNamespace;

          dispatch({
            type: ActionType.InputAlarmWorkLoadNameSpace,
            payload: hasNamespace ? reverseReduceNs(alarmpolicy.clusterId, alarmpolicy.alarmObjectNamespace) : 'ALL'
          });

          dispatch(alarmPolicyActions.inputAlarmPolicyObjects(alarmpolicy.alarmObjetcs));
        }
      }
    };
  },
  initAlarmPolicyData: () => {
    return (dispatch, getState: GetState) => {
      const { route, alarmPolicy } = getState(),
        urlParams = router.resolve(route),
        mode = urlParams['sub'];
      dispatch(resourceActions.channel.fetch());
      dispatch(resourceActions.template.fetch());
      dispatch(resourceActions.receiverGroup.fetch());
      if (mode === 'create') {
        /// #if tke
        dispatch(alarmPolicyActions.inputAlarmPolicyType('cluster'));
        dispatch(
          namespaceActions.applyFilter({
            regionId: route.queries['rid'],
            clusterId: route.queries['clusterId'],
            default: true
          })
        );
        /// #endif
        /// #if project
        dispatch(alarmPolicyActions.inputAlarmPolicyType('pod'));
        /// #endif
      } else if (mode === 'update' || mode === 'copy') {
        const alarmPolicyId = route.queries['alarmPolicyId'];
        const finder = alarmPolicy.list.data.records.find(item => item.id === alarmPolicyId);

        dispatch(
          resourceActions.receiverGroup.fetch({
            data: finder.receiverGroups
          })
        );
        //初始化workload列表不使用初始值
        /// #if tke
        if (mode === 'update' && ['pod', 'virtualMachine'].includes(finder?.alarmPolicyType)) {
          dispatch(
            namespaceActions.applyFilter({
              regionId: route.queries['rid'],
              clusterId: route.queries['clusterId'],
              default: false
            })
          );
        } else {
          dispatch(
            namespaceActions.applyFilter({
              regionId: route.queries['rid'],
              clusterId: route.queries['clusterId'],
              default: true
            })
          );
        }
        /// #endif

        if (mode === 'update') {
          dispatch(alarmPolicyActions.initAlarmPolicyEditionForUpdate(finder));
        } else {
          dispatch(alarmPolicyActions.initAlarmPolicyEditionForCopy(finder));
        }
      } else if (mode === 'detail') {
        const alarmPolicyId = route.queries['alarmPolicyId'];
        const finder = alarmPolicy.list.data.records.find(item => item.id === alarmPolicyId);
        dispatch(alarmPolicyActions.fetchAlarmPolicyDetail(finder));
      }
    };
  },

  inputAlarmPolicyChannel: payload => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPolicyChannel,
        payload
      });
    };
  },

  inputAlarmPolicyTemplate: payload => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPolicyTemplate,
        payload
      });
    };
  },

  inputAlarmPolicyName: alarmPolicyName => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPolicyName,
        payload: alarmPolicyName
      });
    };
  },
  inputAlarmPolicyDescription: alarmPolicyDescription => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPolicyDescription || '',
        payload: alarmPolicyDescription
      });
    };
  },
  //策略类型cluster//nodo//pod/vm
  inputAlarmPolicyType: alarmPolicyType => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPolicyType,
        payload: alarmPolicyType
      });
      const defaultAlarmPolicyObjectsType = 'all';
      dispatch({
        type: ActionType.InputAlarmPolicyObjectsType,
        payload: defaultAlarmPolicyObjectsType
      });
      dispatch(alarmPolicyActions.initAlarmMetrics(alarmPolicyType));
    };
  },

  inputAlarmPolicyStatisticsPeriod: statisticsPeriod => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPolicyStatisticsPeriod,
        payload: statisticsPeriod
      });
    };
  },

  inputAlarmPolicyObjects: alarmPolicyObjects => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPolicyObjects,
        payload: alarmPolicyObjects
      });
    };
  },

  //alarm选择告警对象类型
  inputAlarmPolicyObjectsType: objectType => {
    return (dispatch, getState: GetState) => {
      const { alarmPolicyEdition, namespaceList } = getState();
      dispatch({
        type: ActionType.InputAlarmPolicyObjectsType,
        payload: objectType
      });
      if (alarmPolicyEdition.alarmPolicyType === 'pod') {
        if (objectType === 'part') {
          if (namespaceList.data.recordCount > 0 && alarmPolicyEdition.alarmObjectNamespace === 'ALL') {
            dispatch({
              type: ActionType.InputAlarmWorkLoadNameSpace,
              payload: namespaceList.data.records[0].name
            });
          }
          dispatch(alarmPolicyActions.inputAlarmObjectWorkloadType('Deployment'));
        } else {
          dispatch({
            type: ActionType.InputAlarmObjectWorkloadType,
            payload: 'ALL'
          });
          dispatch({
            type: ActionType.InputAlarmWorkLoadNameSpace,
            payload: 'ALL'
          });
        }
      }
    };
  },

  inputAlarmMetrics: (id: string, Obj: Object) => {
    return (dispatch, getState: GetState) => {
      const { alarmPolicyEdition } = getState(),
        { alarmMetrics } = alarmPolicyEdition;
      const newAlarmMetrics = deepClone(alarmMetrics),
        index = newAlarmMetrics.findIndex(e => e.id === id),
        objKeys = Object.keys(Obj);
      objKeys.forEach(item => {
        newAlarmMetrics[index][item] = Obj[item];
      });
      dispatch({
        type: ActionType.InputAlarmMetrics,
        payload: newAlarmMetrics
      });
    };
  },
  //更新告警设置初始化
  //将配置中
  initAlarmMetricsForUpdate: (alarmMetrics, alarmPolicyType) => {
    return (dispatch, getState: GetState) => {
      const finalType = alarmPolicyType;

      const alarmPolicyMetricsConfig = deepClone(AlarmPolicyMetrics[finalType]);
      const initalarmMetrics = alarmMetrics.length
        ? alarmMetrics.map(item => {
            const index = alarmPolicyMetricsConfig.findIndex(metrics => metrics.metricName === item.metricName);
            if (index !== -1) {
              alarmPolicyMetricsConfig.splice(index, 1);
            }
            return Object.assign({}, item, {
              id: uuid(),
              v_evaluatorValue: initValidator,
              enable: true
            });
          })
        : [];
      //将返回的配置中没有启用的指标项加上
      alarmPolicyMetricsConfig.forEach(item => {
        initalarmMetrics.push(
          Object.assign({}, item, {
            id: uuid(),
            v_evaluatorValue: initValidator,
            enable: false
          })
        );
      });

      dispatch({
        type: ActionType.InputAlarmMetrics,
        payload: initalarmMetrics
      });
    };
  },
  //告警设置不同类型不同初始值
  initAlarmMetrics: (type: string) => {
    return (dispatch, getState: GetState) => {
      const finalType = type;
      if (type === AlarmPolicyType[0].value) {
        const { cluster } = getState(),
          finder = cluster.list.data.records.find(
            record => cluster.selection && record.metadata.name === cluster.selection.metadata.name
          );
        // if (finder && finder.clusterType === 'INDEPENDENT_CLUSTER') {
        //   finalType = 'independentClusetr';
        // }
      }

      const items = AlarmPolicyMetrics[finalType],
        alarmMetrics = items
          ? items.map(item => {
              return Object.assign({}, item, {
                id: uuid(),
                v_evaluatorValue: initValidator
              });
            })
          : [];
      dispatch({
        type: ActionType.InputAlarmMetrics,
        payload: alarmMetrics
      });
    };
  },

  inputAlarmNotifyWays: notifyWays => {
    return (dispatch, getState: GetState) => {
      const newNotifyWays = deepClone(notifyWays);

      dispatch({
        type: ActionType.InputAlarmNotifyWay,
        payload: newNotifyWays
      });
    };
  },

  inputAlarmNotifyWay: (id: string, obj: Object) => {
    return (dispatch, getState: GetState) => {
      const { alarmPolicyEdition } = getState(),
        { notifyWays } = alarmPolicyEdition;
      const newNotifyWays = deepClone(notifyWays),
        index = newNotifyWays.findIndex(e => e.id === id),
        objKeys = Object.keys(obj);
      objKeys.forEach(item => {
        newNotifyWays[index][item] = obj[item];
      });
      dispatch({
        type: ActionType.InputAlarmNotifyWay,
        payload: newNotifyWays
      });
    };
  },

  deleteAlarmNotifyWay: (id: string) => {
    return (dispatch, getState: GetState) => {
      const { alarmPolicyEdition } = getState(),
        { notifyWays } = alarmPolicyEdition;
      const newNotifyWays = deepClone(notifyWays),
        index = newNotifyWays.findIndex(e => e.id === id);
      newNotifyWays.splice(index, 1);
      dispatch({
        type: ActionType.InputAlarmNotifyWay,
        payload: newNotifyWays
      });
    };
  },

  addAlarmNotifyWay: () => {
    return (dispatch, getState: GetState) => {
      const { alarmPolicyEdition } = getState(),
        { notifyWays } = alarmPolicyEdition;
      const newNotifyWays = deepClone(notifyWays);
      newNotifyWays.push({ id: uuid(), channel: undefined, template: undefined });
      dispatch({
        type: ActionType.InputAlarmNotifyWay,
        payload: newNotifyWays
      });
    };
  },

  inputAlarmShieldTimeStart: time => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmShieldTimeStart,
        payload: time
      });
      dispatch(validatorActions.validatePolicyTime());
    };
  },

  inputAlarmShieldTimeEnd: time => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmShieldTimeEnd,
        payload: time
      });
      dispatch(validatorActions.validatePolicyTime());
    };
  },

  inputAlarmPhoneCircleTimes: value => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPhoneCircleTimes,
        payload: value
      });
    };
  },
  inputAlarmPhoneInnerInterval: value => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPhoneInnerInterval,
        payload: value
      });
    };
  },
  inputAlarmPhoneCircleInterval: value => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPhoneCircleInterval,
        payload: value
      });
    };
  },
  inputAlarmPhoneArriveNotice: value => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPhoneArriveNotice,
        payload: value
      });
    };
  },
  selectsWorkLoadNamespace: namespace => {
    return (dispatch, getState: GetState) => {
      const { regionSelection, cluster, alarmPolicyEdition } = getState();
      dispatch({
        type: ActionType.InputAlarmWorkLoadNameSpace,
        payload: namespace
      });
      dispatch({
        type: ActionType.InputAlarmPolicyObjects,
        payload: []
      });
      if (
        (alarmPolicyEdition.alarmPolicyType === 'pod' && alarmPolicyEdition.alarmObjectsType === 'all') ||
        alarmPolicyEdition.alarmPolicyType === 'virtualMachine'
      ) {
        //
      } else {
        dispatch(
          workloadActions.applyFilter({
            regionId: +regionSelection.value,
            clusterId: cluster.selection ? cluster.selection.metadata.name : '',
            namespace: namespace,
            workloadType: alarmPolicyEdition.alarmObjectWorkloadType
          })
        );
      }
    };
  },
  inputAlarmObjectWorkloadType: (type: string) => {
    return (dispatch, getState: GetState) => {
      const { regionSelection, cluster, alarmPolicyEdition } = getState();
      dispatch({
        type: ActionType.InputAlarmObjectWorkloadType,
        payload: type
      });
      dispatch({
        type: ActionType.InputAlarmPolicyObjects,
        payload: []
      });
      if (alarmPolicyEdition.alarmPolicyType === 'pod' && alarmPolicyEdition.alarmObjectsType === 'all') {
        //
      } else {
        dispatch(
          workloadActions.applyFilter({
            regionId: +regionSelection.value,
            clusterId: cluster.selection ? cluster.selection.metadata.name : '',
            namespace: alarmPolicyEdition.alarmObjectNamespace,
            workloadType: type
          })
        );
      }
    };
  },

  clearAlarmPolicyEdit: () => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.ClearAlarmPolicyEdit,
        payload: {}
      });
      dispatch(resourceActions.receiverGroup.selects([]));
      // groupActions.selectGroup([]);
    };
  }
};

export const alarmPolicyActions = extend({}, _alarmPolicyActions, restActions, editActions);
