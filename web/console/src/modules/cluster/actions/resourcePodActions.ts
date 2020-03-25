import {
    extend, FetchOptions, generateFetcherActionCreator, RecordSet, ReduxAction
} from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { apiVersion } from '../../../../config/resource/common';
import { ResourceConfigVersionMap } from '../../../../config/resourceConfig';
import { cloneDeep } from '../../../../src/modules/common';
import { ResourceInfo } from '../../common/models';
import { IsInNodeManageDetail } from '../components/resource/resourceDetail/ResourceDetail';
import {
    IsPodShowLoadingIcon, reduceContainerId
} from '../components/resource/resourceDetail/ResourcePodPanel';
import * as ActionType from '../constants/ActionType';
import { PollEventName } from '../constants/Config';
import { Pod, PodFilterInNode, ResourceFilter, RootState } from '../models';
import { TappGrayUpdateEditItem } from '../models/ResourceDetailState';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { resourceDetailActions } from './resourceDetailActions';
import { resourcePodLogActions } from './resourcePodLogActions';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** 获取Pod列表的action */
const fetchPodActions = generateFetcherActionCreator({
  actionType: ActionType.FetchPodList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, route, clusterVersion } = getState(),
      urlParams = router.resolve(route),
      { resourceDetailState, resourceInfo } = subRoot,
      { podQuery, podFilterInNode } = resourceDetailState;

    let isInNodeManager = IsInNodeManageDetail(urlParams['type']);
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response: RecordSet<any>;
    /**
     * workload里面拉取pods，因为workload集成了子资源，所以直接拉取workload的pods资源，即调用fetchExtraResourceList
     * 但，node详情里面，需要通过fieldSelector当中的
     */
    if (isInNodeManager) {
      // pods的apiVersion的配置
      let podVersionInfo = apiVersion[ResourceConfigVersionMap(clusterVersion)]['pods'];
      let { podName, phase, namespace } = podFilterInNode;
      // pods的resourceInfo的配置
      let podResourceInfo: ResourceInfo = {
        basicEntry: podVersionInfo.basicEntry,
        version: podVersionInfo.version,
        group: podVersionInfo.group,
        namespaces: '',
        requestType: {
          list: 'pods'
        }
      };
      // 过滤条件
      let k8sQueryObj = {
        fieldSelector: {
          'spec.nodeName': route.queries['resourceIns'] ? route.queries['resourceIns'] : undefined,
          'metadata.namespace': namespace ? namespace : undefined,
          'metadata.name': podName ? podName : undefined,
          'status.phase': phase ? phase : undefined
        }
      };
      k8sQueryObj = JSON.parse(JSON.stringify(k8sQueryObj));
      response = await WebAPI.fetchResourceList(podQuery, podResourceInfo, isClearData, k8sQueryObj, false, false);
    } else {
      // 这里是workload里面的拉取pod列表的逻辑
      response = await WebAPI.fetchExtraResourceList(podQuery, resourceInfo, isClearData, 'pods');
    }
    // 原因为 Evicted的pod没有必要再进行展示，直接进行过滤
    response.records = response.records.filter(item => item.status.reason !== 'Evicted');
    response.recordCount = response.records.length;
    return response;
  },
  finish: async (dispatch, getState: GetState) => {
    let { route, subRoot } = getState(),
      urlParams = router.resolve(route),
      { podList, logOption } = subRoot.resourceDetailState,
      { podName } = logOption;

    let isInNodeManager = IsInNodeManageDetail(urlParams['type']);

    // 这里去初始化containerList的列表
    let containerList = [];
    for (let pod of podList.data.records) {
      for (let container of pod.spec.containers) {
        container.id = reduceContainerId(pod.status.containerStatuses, container.name);
        containerList.push(container);
      }
    }
    dispatch(resourcePodActions.getContainerList(containerList));

    // 这里是去判断需不需要轮询，如果不在pod 的tab页面，直接停止podList的轮询，否则根据轮询条件判断
    if (urlParams['tab'] !== 'pod' && urlParams['tab'] !== '') {
      dispatch(resourcePodActions.clearPollEvent());
    } else {
      if (podList.data.recordCount) {
        if (podList.data.records.filter(item => IsPodShowLoadingIcon(item)).length === 0) {
          dispatch(resourcePodActions.clearPollEvent());
        }
      } else {
        dispatch(resourcePodActions.clearPollEvent());
      }
    }

    /**
     * pre: 不在Node详情页内的pod列表的展示
     * 拉取完之后，需要去触发一下 详情-日志页面 的选择，自动选择第一个pod，同时，得判断已经选择过podName的话，不需要继续选择
     */
    if (!isInNodeManager && podList.data.recordCount && podName === '') {
      let podName = podList.data.records[0].metadata.name;
      dispatch(resourcePodLogActions.selectPod(podName));
    }
  }
});

/** 查询Pod列表action */
const queryPodActions = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.QueryPodList,
  bindFetcher: fetchPodActions
});

/** 剩余的pod的操作 */
const restActions = {
  getContainerList: payload => {
    return {
      type: ActionType.FetchContainerList,
      payload
    };
  },

  /** pod的选择 */
  podSelect: (pods: Pod[]): ReduxAction<Pod[]> => {
    return {
      type: ActionType.PodSelection,
      payload: pods
    };
  },

  /**Tapp灰度升级编辑项 */
  initTappGrayUpdate: (items: TappGrayUpdateEditItem[]): ReduxAction<TappGrayUpdateEditItem[]> => {
    return {
      type: ActionType.W_TappGrayUpdate,
      payload: items
    };
  },
  updateTappGrayUpdate: (index_out, index_in, imageName, imageTag) => {
    return async (dispatch, getState: GetState) => {
      let { editTappGrayUpdate } = getState().subRoot.resourceDetailState;
      let target: TappGrayUpdateEditItem[] = cloneDeep(editTappGrayUpdate);
      target[index_out].containers[index_in].imageName = imageName;
      target[index_out].containers[index_in].imageTag = imageTag;

      dispatch({
        type: ActionType.W_TappGrayUpdate,
        payload: target
      });
    };
  },
  /** 是否展示 登录弹框 */
  toggleLoginDialog: () => {
    return async (dispatch, getState: GetState) => {
      let { isShowLoginDialog } = getState().subRoot.resourceDetailState;
      dispatch({
        type: ActionType.IsShowLoginDialog,
        payload: !isShowLoginDialog
      });
    };
  },

  /** 轮询拉取条件 */
  poll: (queryObj: ResourceFilter) => {
    return async (dispatch, getState: GetState) => {
      // 每次轮询之前先清空之前的轮询
      dispatch(resourcePodActions.clearPollEvent());
      // 触发列表的查询
      dispatch(resourcePodActions.applyFilter(queryObj));

      window[PollEventName['resourcePodList']] = setInterval(() => {
        dispatch(resourcePodActions.poll(queryObj));
      }, 10000);
    };
  },

  /** 清空轮询条件 */
  clearPollEvent: () => {
    return async (dispatch, getState: GetState) => {
      clearInterval(window[PollEventName['resourcePodList']]);
    };
  },

  /** 选择pod的筛选项 */
  updatePodFilterInNode: (podFilter: PodFilterInNode) => {
    return async (dispatch, getState: GetState) => {
      let { podQuery } = getState().subRoot.resourceDetailState;

      dispatch({
        type: ActionType.PodFilterInNode,
        payload: podFilter
      });

      // 根据筛选项，进行pod列表的查询，namespace、metadata.name等过滤条件
      dispatch(resourceDetailActions.pod.poll(podQuery.filter));
    };
  }
};

export const resourcePodActions = extend(fetchPodActions, queryPodActions, restActions);
