import {
    extend, FetchOptions, generateFetcherActionCreator, RecordSet, ReduxAction
} from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { ResourceInfo } from '../../common/models/ResourceInfo';
import * as ActionType from '../constants/ActionType';
import { Computer, ComputerFilter, RootState } from '../models';
import { ResourceFilter } from '../models/ResourceOption';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/**
 * 获取节点 pod列表action
 */
const fetchComputerPodActions = generateFetcherActionCreator({
  actionType: ActionType.FetchComputerPodList,
  fetcher: async (getState: GetState) => {
    let {
      clusterVersion,
      subRoot: {
        computerState: { computerPodQuery }
      }
    } = getState();
    // pods的apiVersion的配置
    let podVersionInfo = resourceConfig(clusterVersion)['pods'];
    let { specificName, clusterId } = computerPodQuery.filter;
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
        'spec.nodeName': specificName
      }
    };

    k8sQueryObj = JSON.parse(JSON.stringify(k8sQueryObj));
    let response = await WebAPI.fetchResourceList(
      { filter: { clusterId: clusterId }, search: '' },
      podResourceInfo,
      false,
      k8sQueryObj
    );
    return response;
  }
});

/**
 * 查询节点 pod列表Action
 */
const queryComputerPodActions = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.QueryComputerPodList,
  bindFetcher: fetchComputerPodActions
});

/**
 * 其他Action
 */
const restActions = {};

export const computerPodActions = extend({}, queryComputerPodActions, fetchComputerPodActions, restActions);
