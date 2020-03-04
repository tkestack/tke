import {
    extend, FetchOptions, generateFetcherActionCreator, ReduxAction, uuid
} from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { cloneDeep } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import { SessionAffinity } from '../constants/Config';
import { initPortsMap, initSelector } from '../constants/initState';
import { CLB, PortMap, Resource, RootState, Selector, ServicePorts } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { validateServiceActions } from './validateServiceActions';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** ===== start workload的相关选择 ============================= */
const fetchWorkloadActions = generateFetcherActionCreator({
  actionType: ActionType.S_FetchWorkloadList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { serviceEdit } = subRoot,
      { workloadQuery, workloadType } = serviceEdit;

    let workloadResourceInfo = resourceConfig(clusterVersion)[workloadType];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await WebAPI.fetchResourceList(workloadQuery, workloadResourceInfo, isClearData);
    return response;
  },
  finish: async (dispatch, getState: GetState) => {
    let { workloadList } = getState().subRoot.serviceEdit;

    // 如果拉回来的列表有数据，则默认选择第一项
    if (workloadList.data.recordCount) {
      dispatch(serviceEditActions.workload.selectWorkload([workloadList.data.records[0]]));
    } else {
      dispatch(serviceEditActions.workload.selectWorkload([]));
    }
  }
});

const queryWorkloadActions = generateQueryActionCreator({
  actionType: ActionType.S_QueryWorkloadList,
  bindFetcher: fetchWorkloadActions
});

const restActions = {
  /** 选择某个具体的workload */
  selectWorkload: (workload: Resource[]): ReduxAction<any> => {
    return {
      type: ActionType.S_WorkloadSelection,
      payload: workload
    };
  },

  /** 选择当前的workloadType */
  selectWorkloadType: (type: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.S_WorkloadType,
        payload: type
      });
      // 变更了类型之后需要重新拉取当前的资源列表，在拉取的时候，判断当前的workloadType
      dispatch(serviceEditActions.workload.fetch());
    };
  },

  /** 是否展示workloadDialog */
  toggleIsShowWorkloadDialog: () => {
    return async (dispatch, getState: GetState) => {
      let { isShowWorkloadDialog } = getState().subRoot.serviceEdit;
      dispatch({
        type: ActionType.S_IsShowWorkloadDialog,
        payload: !isShowWorkloadDialog
      });
    };
  }
};

const workloadActions = extend(fetchWorkloadActions, queryWorkloadActions, restActions);
/** ===== end workload的相关选择 ============================= */

export const serviceEditActions = {
  /** workload的相关操作 */
  workload: workloadActions,

  /** 更新服务名称 */
  inputServiceName: (name: string): ReduxAction<string> => {
    return {
      type: ActionType.S_ServiceName,
      payload: name
    };
  },

  /** 更新服务的描述 */
  inputServiceDesp: (desp: string): ReduxAction<string> => {
    return {
      type: ActionType.S_Description,
      payload: desp
    };
  },

  /** 选择命名空间 */
  selectNamespace: (namespace: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route } = getState();

      dispatch({
        type: ActionType.S_Namespace,
        payload: namespace
      });

      // 验证命名空间的选择是否有效
      dispatch(validateServiceActions.validateNamespace());

      // 重新拉取workload
      dispatch(
        serviceEditActions.workload.applyFilter({
          clusterId: route.queries['clusterId'],
          namespace,
          regionId: +route.queries['rid']
        })
      );
    };
  },

  /** 选择访问的方式 */
  selectCommunicationType: (communication: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { isOpenHeadless } = getState().subRoot.serviceEdit;

      dispatch({
        type: ActionType.S_CommunicationType,
        payload: communication
      });

      // 这里需要去判断，如果切换的时候，headless 为 true，切换到别的模式，则需要设置为false，保留true，别的方式创建的时候可能会导致错误
      if (isOpenHeadless) {
        dispatch(serviceEditActions.isOpenHeadless(isOpenHeadless));
      }
    };
  },

  /** 更新端口映射 */
  updatePortMap: (obj: any, portMapId: string) => {
    return async (dispatch, getState: GetState) => {
      let { portsMap } = getState().subRoot.serviceEdit;
      let newPortsMap: PortMap[] = cloneDeep(portsMap);

      let portsIndex = newPortsMap.findIndex(item => item.id === portMapId);
      let keyArr = Object.keys(obj);
      keyArr.forEach(item => {
        newPortsMap[portsIndex][item] = obj[item];
      });

      dispatch({
        type: ActionType.S_UpdatePortsMap,
        payload: newPortsMap
      });
    };
  },

  /** 删除端口映射 */
  deletePortMap: (portMapId: string) => {
    return async (dispatch, getState: GetState) => {
      let { portsMap } = getState().subRoot.serviceEdit;
      let newPortsMap: PortMap[] = cloneDeep(portsMap);

      let portsIndex = newPortsMap.findIndex(item => item.id === portMapId);
      newPortsMap.splice(portsIndex, 1);

      dispatch({
        type: ActionType.S_UpdatePortsMap,
        payload: newPortsMap
      });
    };
  },

  /** 增加端口映射 */
  addPortMap: () => {
    return async (dispatch, getState: GetState) => {
      let { portsMap } = getState().subRoot.serviceEdit;
      let newPortsMap: PortMap[] = cloneDeep(portsMap);
      newPortsMap.push(
        Object.assign({}, initPortsMap, {
          id: uuid()
        })
      );

      dispatch({
        type: ActionType.S_UpdatePortsMap,
        payload: newPortsMap
      });
    };
  },

  /** 是否开启headless service */
  isOpenHeadless: (isOpen: boolean) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.S_IsOpenHeadless,
        payload: !isOpen
      });
    };
  },

  /** 增加selector 配置项 */
  addSelector: () => {
    return async (dispatch, getState: GetState) => {
      let selectors = cloneDeep(getState().subRoot.serviceEdit.selector);

      selectors.push(
        Object.assign({}, initSelector, {
          id: uuid()
        })
      );
      dispatch({
        type: ActionType.S_Selector,
        payload: selectors
      });
    };
  },

  /** 更新Selectors的操作 */
  updateSelectorConfig: (obj: any, sId: string) => {
    return async (dispatch, getState: GetState) => {
      let selectors: Selector[] = cloneDeep(getState().subRoot.serviceEdit.selector),
        sIndex = selectors.findIndex(s => s.id === sId);

      let keyName = Object.keys(obj)[0];

      selectors[sIndex][keyName] = obj[keyName];
      dispatch({
        type: ActionType.S_Selector,
        payload: selectors
      });
    };
  },

  /** 初始化selectors */
  initSelectorFromWorkload: (selectors: Selector[]): ReduxAction<any> => {
    return {
      type: ActionType.S_Selector,
      payload: selectors
    };
  },

  /** 删除selector的操作 */
  deleteSelectorContent: (sId: string) => {
    return async (dispatch, getState: GetState) => {
      let selectors: Selector[] = cloneDeep(getState().subRoot.serviceEdit.selector),
        sIndex = selectors.findIndex(s => s.id === sId);

      selectors.splice(sIndex, 1);
      dispatch({
        type: ActionType.S_Selector,
        payload: selectors
      });
    };
  },

  /** 初始化portsMap */
  initPortsMapForUpdate: (portsMap: ServicePorts[]) => {
    return async (dispatch, getState: GetState) => {
      let newPortsMap: PortMap[] = [];

      newPortsMap = portsMap.map(item => {
        let tmp: PortMap = Object.assign({}, initPortsMap, {
          id: uuid(),
          protocol: item.protocol,
          targetPort: item.targetPort + '',
          port: item.port + '',
          nodePort: item.nodePort ? item.nodePort + '' : ''
        });

        return tmp;
      });

      dispatch({
        type: ActionType.S_UpdatePortsMap,
        payload: newPortsMap
      });
    };
  },

  /** 更新访问方式的时候，初始化一些数据 */
  initServiceEditForUpdate: (resource: Resource) => {
    return async (dispatch, getState: GetState) => {
      let annotations = resource['metadata']['annotations'];

      let resourceType = resource.spec.type;

      // 如果是ClusterIP的方式，并且clusterIP为None，则开启headless的设置
      if (resourceType === 'ClusterIP' && resource.spec.clusterIP === 'None') {
        dispatch(serviceEditActions.isOpenHeadless(false));
      }

      dispatch(serviceEditActions.selectCommunicationType(resourceType));
      dispatch(serviceEditActions.initPortsMapForUpdate(resource.spec.ports));

      if (resource.spec.externalTrafficPolicy) {
        dispatch(serviceEditActions.chooseExternalTrafficPolicyMode(resource.spec.externalTrafficPolicy));
      }

      if (resource.spec.sessionAffinity) {
        dispatch(serviceEditActions.chooseChoosesessionAffinityMode(resource.spec.sessionAffinity));
        resource.spec.sessionAffinity === SessionAffinity.ClientIP &&
          dispatch(
            serviceEditActions.inputsessionAffinityTimeout(
              (resource.spec.sessionAffinityConfig &&
                resource.spec.sessionAffinityConfig.clientIP &&
                resource.spec.sessionAffinityConfig.clientIP.timeoutSeconds) ||
                30
            )
          );
      }
    };
  },

  /**设置创建service高级设置 */

  chooseExternalTrafficPolicyMode: (mode: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.S_ChooseExternalTrafficPolicy,
        payload: mode
      });
    };
  },

  chooseChoosesessionAffinityMode: (mode: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.S_ChoosesessionAffinity,
        payload: mode
      });
    };
  },

  inputsessionAffinityTimeout: (seconds: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.S_InputsessionAffinityTimeout,
        payload: seconds
      });
    };
  },

  /** 离开创建页面，清除 serviceEdit当中的内容 */
  clearServiceEdit: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearServiceEdit
    };
  }
};
