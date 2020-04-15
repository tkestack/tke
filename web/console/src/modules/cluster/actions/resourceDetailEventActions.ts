import { createFFListActions, extend, ReduxAction } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config';
import { IsInNodeManageDetail } from '../components/resource/resourceDetail/ResourceDetail';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { Event, ResourceFilter, RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

/** 获取事件列表 FFRedux */
const FFModelEventActions = createFFListActions<Event, ResourceFilter>({
  actionName: FFReduxActionName.DETAILEVENT,
  fetcher: async (query, getState: GetState, fetchOptions) => {
    let { subRoot, route, clusterVersion } = getState(),
      urlParams = router.resolve(route),
      { resourceDetailState, resourceInfo } = subRoot,
      { event } = resourceDetailState;

    let isInNodeManager = IsInNodeManageDetail(urlParams['type']);
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;

    /**
     * workload里面拉取events，是因为workload集成了events的子资源，所以直接拉workload的events资源，即调用fetchExtraResourceList来进行子资源的拉取，类似pods
     * 但，node详情里面，需要通过fieldSelector当中的involvedObject.kind来拉取
     */
    if (isInNodeManager) {
      // event的resourceInfo的配置
      let eventResourceInfo = resourceConfig(clusterVersion)['event'];
      // 过滤条件
      let k8sQueryObj = {
        fieldSelector: {
          'involvedObject.kind': 'Node',
          'involvedObject.name': route.queries['resourceIns'] ? route.queries['resourceIns'] : undefined
        }
      };
      k8sQueryObj = JSON.parse(JSON.stringify(k8sQueryObj));
      let response = await WebAPI.fetchResourceList(event.query, {
        resourceInfo: eventResourceInfo,
        isClearData,
        k8sQueryObj,
        isNeedDes: true
      });
      return response;
    } else {
      let response;
      if (resourceInfo.requestType.useDetailInfo) {
        let {
          subRoot: {
            detailResourceOption: { detailResourceInfo }
          }
        } = getState();
        response = await WebAPI.fetchExtraResourceList(
          event.query,
          detailResourceInfo,
          isClearData,
          'events',
          {},
          true
        );
      } else {
        response = await WebAPI.fetchExtraResourceList(event.query, resourceInfo, isClearData, 'events', {}, true);
      }

      return response;
    }
  },
  getRecord: (getState: GetState) => {
    return getState().subRoot.resourceDetailState.event;
  }
});

/** 剩余的Action的操作 */
const restActions = {
  /** 轮询事件的列表 */
  poll: () => {
    return async (dispatch, getState: GetState) => {
      let {
          route,
          subRoot: { resourceInfo, detailResourceOption }
        } = getState(),
        urlParams = router.resolve(route);

      let { clusterId, np, resourceIns, rid } = route.queries;

      let currentIns;
      if (resourceInfo.requestType && resourceInfo.requestType.useDetailInfo) {
        currentIns = detailResourceOption.detailResourceSelection;
      } else if (!IsInNodeManageDetail(urlParams['type'])) {
        currentIns = resourceIns;
      }

      let eventFilter: ResourceFilter = {
        clusterId,
        namespace: np,
        regionId: +rid,
        specificName: currentIns
      };

      dispatch(
        FFModelEventActions.polling({
          filter: eventFilter,
          delayTime: 10000,
          onError: (dispatch: Redux.Dispatch) => {
            // 如果发生了错误，把自动刷新的按钮置灰
            dispatch(resourceDetailEventActions.switchAutoPolling(false));
          }
        })
      );
    };
  },

  /** 进行自动刷新按钮的切换 */
  switchAutoPolling: (isAuto: boolean): ReduxAction<boolean> => {
    return {
      type: ActionType.IsAutoPollingEvent,
      payload: isAuto
    };
  }
};

export const resourceDetailEventActions = extend(FFModelEventActions, restActions);
