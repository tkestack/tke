import { namespaceActions } from './namespaceActions';
import { router } from './../router';
import { AlarmPolicyType } from './../constants/Config';
import { extend, deepClone, uuid } from '@tencent/qcloud-lib';
import { RootState, AlarmPolicy, AlarmPolicyFilter } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../WebAPI';
import { AlarmPolicyMetrics } from '../constants/Config';
import { workloadActions } from './workloadActions';
import { validatorActions } from './validatorActions';
import { resourceActions } from '../../notify/actions/resourceActions';
import { initValidator } from '../../common/models';
import { userActions } from '../../uam/actions/userActions';
import { createListAction } from '@tencent/redux-list';

type GetState = () => RootState;

const _alarmPolicyActions = createListAction<AlarmPolicy, AlarmPolicyFilter>({
  actionName: 'AlarmPolicy',
  fetcher: async (query, getstate: GetState) => {
    const response = await WebAPI.fetchAlarmPolicy(query);

    //业务侧中过滤只有这个namepace下的AlarmPolicy
    /// #if project
    response.records = response.records.filter(
      item =>
        item.alarmObjetcsType === 'part' &&
        item.alarmObjectNamespace === getstate().namespaceSelection &&
        item.alarmPolicyType === 'pod'
    );
    response.recordCount = response.records.length;
    /// #endif

    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().alarmPolicy;
  },
  onFinish: (record, dispatch, getState: GetState) => {
    let { sub } = router.resolve(getState().route);
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

      // dispatch(groupActions.selectGroup(alarmpolicy.receiverGroups));
      dispatch(alarmPolicyActions.initAlarmMetricsForUpdate(alarmpolicy.alarmMetrics, alarmpolicy.alarmPolicyType));
      // dispatch(alarmPolicyActions.inputAlarmShieldTimeStart(alarmpolicy.shieldTimeStart || null));
      // dispatch(alarmPolicyActions.inputAlarmShieldTimeEnd(alarmpolicy.shieldTimeEnd));
      dispatch(
        alarmPolicyActions.inputAlarmNotifyWays(
          alarmpolicy.notifyWays.map(w => ({
            id: uuid(),
            channel: w.channel,
            template: w.template
          }))
        )
      );
      // dispatch({
      //   type: ActionType.InputAlarmNotifyWay,
      //   payload: alarmpolicy.notifyWay
      // });
      // //当告警渠道有电话时初始化
      // if (alarmpolicy.notifyWay.indexOf('CALL') !== -1) {
      //   dispatch(alarmPolicyActions.inputAlarmPhoneCircleInterval(alarmpolicy.phoneCircleInterval));
      //   dispatch(alarmPolicyActions.inputAlarmPhoneArriveNotice(alarmpolicy.phoneNotifyOrder));
      //   dispatch(alarmPolicyActions.inputAlarmPhoneCircleTimes(alarmpolicy.phoneCircleTimes));
      //   dispatch(alarmPolicyActions.inputAlarmPhoneInnerInterval(alarmpolicy.phoneInnerInterval));
      // }
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
        if (alarmpolicy.alarmObjetcsType === 'part') {
          dispatch({
            type: ActionType.InputAlarmWorkLoadNameSpace,
            payload: alarmpolicy.alarmObjectNamespace
          });
          dispatch(alarmPolicyActions.inputAlarmObjectWorkloadType(alarmpolicy.alarmObjectWorkloadType));
          dispatch(alarmPolicyActions.inputAlarmPolicyObjects(alarmpolicy.alarmObjetcs));
        }
      }
    };
  },
  initAlarmPolicyData: () => {
    return (dispatch, getState: GetState) => {
      let { route, alarmPolicy } = getState(),
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
        let alarmPolicyId = route.queries['alarmPolicyId'];
        let finder = alarmPolicy.list.data.records.find(item => item.id === alarmPolicyId);

        dispatch(
          resourceActions.receiverGroup.fetch({
            data: finder.receiverGroups
          })
        );
        //初始化workload列表不使用初始值
        /// #if tke
        if (mode === 'update' && finder.alarmPolicyType === 'pod' && finder.alarmObjetcsType === 'part') {
          namespaceActions.applyFilter({
            regionId: route.queries['rid'],
            clusterId: route.queries['clusterId'],
            default: false
          });
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
        let alarmPolicyId = route.queries['alarmPolicyId'];
        let finder = alarmPolicy.list.data.records.find(item => item.id === alarmPolicyId);
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
  //策略类型cluster//nodo//pod
  inputAlarmPolicyType: alarmPolicyType => {
    return (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.InputAlarmPolicyType,
        payload: alarmPolicyType
      });
      let defaultAlarmPolicyObjectsType = alarmPolicyType === 'pod' ? 'part' : 'all';
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
      dispatch({
        type: ActionType.InputAlarmPolicyObjectsType,
        payload: objectType
      });
    };
  },

  inputAlarmMetrics: (id: string, Obj: Object) => {
    return (dispatch, getState: GetState) => {
      let { alarmPolicyEdition } = getState(),
        { alarmMetrics } = alarmPolicyEdition;
      let newAlarmMetrics = deepClone(alarmMetrics),
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
      let finalType = alarmPolicyType;
      if (alarmPolicyType === 'cluster') {
        let { cluster } = getState(),
          finder = cluster.list.data.records.find(
            record => cluster.selection && record.metadata.name === cluster.selection.metadata.name
          );
        // if (finder.clusterType === 'INDEPENDENT_CLUSTER') {
        //   finalType = 'independentClusetr';
        // }
      }
      let alarmPolicyMetricsConfig = deepClone(AlarmPolicyMetrics[finalType]);
      let initalarmMetrics = alarmMetrics.length
        ? alarmMetrics.map(item => {
            let index = alarmPolicyMetricsConfig.findIndex(metrics => metrics.metricName === item.metricName);
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
      let finalType = type;
      if (type === AlarmPolicyType[0].value) {
        let { cluster } = getState(),
          finder = cluster.list.data.records.find(
            record => cluster.selection && record.metadata.name === cluster.selection.metadata.name
          );
        // if (finder && finder.clusterType === 'INDEPENDENT_CLUSTER') {
        //   finalType = 'independentClusetr';
        // }
      }

      let items = AlarmPolicyMetrics[finalType],
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
      let newNotifyWays = deepClone(notifyWays);

      dispatch({
        type: ActionType.InputAlarmNotifyWay,
        payload: newNotifyWays
      });
    };
  },

  inputAlarmNotifyWay: (id: string, obj: Object) => {
    return (dispatch, getState: GetState) => {
      let { alarmPolicyEdition } = getState(),
        { notifyWays } = alarmPolicyEdition;
      let newNotifyWays = deepClone(notifyWays),
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
      let { alarmPolicyEdition } = getState(),
        { notifyWays } = alarmPolicyEdition;
      let newNotifyWays = deepClone(notifyWays),
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
      let { alarmPolicyEdition } = getState(),
        { notifyWays } = alarmPolicyEdition;
      let newNotifyWays = deepClone(notifyWays);
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
      let { regionSelection, cluster, alarmPolicyEdition } = getState();
      dispatch({
        type: ActionType.InputAlarmWorkLoadNameSpace,
        payload: namespace
      });
      dispatch({
        type: ActionType.InputAlarmPolicyObjects,
        payload: []
      });
      dispatch(
        workloadActions.applyFilter({
          regionId: +regionSelection.value,
          clusterId: cluster.selection ? cluster.selection.metadata.name : '',
          namespace: namespace,
          workloadType: alarmPolicyEdition.alarmObjectWorkloadType
        })
      );
    };
  },
  inputAlarmObjectWorkloadType: (type: string) => {
    return (dispatch, getState: GetState) => {
      let { regionSelection, cluster, alarmPolicyEdition } = getState();
      dispatch({
        type: ActionType.InputAlarmObjectWorkloadType,
        payload: type
      });
      dispatch({
        type: ActionType.InputAlarmPolicyObjects,
        payload: []
      });
      dispatch(
        workloadActions.applyFilter({
          regionId: +regionSelection.value,
          clusterId: cluster.selection ? cluster.selection.metadata.name : '',
          namespace: alarmPolicyEdition.alarmObjectNamespace,
          workloadType: type
        })
      );
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
