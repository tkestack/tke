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

import {
  extend,
  FetchOptions,
  generateFetcherActionCreator,
  ReduxAction,
  createFFListActions
} from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { RootState, Resource, ResourceFilter } from '../models';
import * as WebAPI from '../WebAPI';
import { resourceDetailEventActions } from './resourceDetailEventActions';
import { resourcePodActions } from './resourcePodActions';
import { resourcePodLogActions } from './resourcePodLogActions';
import { resourceRsActions } from './resourceRsActions';
import { FFReduxActionName } from '../constants/Config';
import { router } from '../router';
import { ResourceInfo } from '@src/modules/common';
import { resourceConfig } from '@config';

const ReduceSecretDataForPsw = (dataInfo: string) => {
  let jsonData = JSON.parse(window.atob(dataInfo)),
    jsonKeys = Object.keys(jsonData)[0];
  return jsonData[jsonKeys]['password'] || '';
};

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** ========================= start 拉取集群详情信息 ========================= */
const FFModelResourceInfoActions = createFFListActions<Resource, ResourceFilter>({
  actionName: FFReduxActionName.Resource_Detail_Info,
  fetcher: async (query, getState: GetState) => {
    let { route, clusterVersion } = getState(),
      urlParams = router.resolve(route);

    let resourceName = urlParams['resourceName'];
    let resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[resourceName];

    let response = await WebAPI.fetchResourceList(query, {
      resourceInfo
    });
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().subRoot.resourceDetailState.resourceDetailInfo;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {
    if (record.data.recordCount) {
      dispatch(FFModelResourceInfoActions.select(record.data.records[0]));
    }
  }
});
/** ========================= end 拉取集群详情信息 ========================= */

export const resourceDetailActions = {
  /** 获取事件的相关操作 */
  event: resourceDetailEventActions,

  /** 获取修订历史版本的相关操作 */
  rs: resourceRsActions,

  /** 获取pod列表的相关操作 */
  pod: resourcePodActions,

  /** 获取日志的相关操作 */
  log: resourcePodLogActions,

  /** 拉取资源的详情信息，只拉取单独一个接口的数据 */
  resourceInfo: FFModelResourceInfoActions,

  /** 获取资源，如果deployment的 yaml文件 */
  fetchResourceYaml: generateFetcherActionCreator({
    actionType: ActionType.FetchYaml,
    fetcher: async (getState: GetState, fetchOptions, dispatch) => {
      let { route, subRoot, namespaceSelection } = getState(),
        {
          resourceInfo,
          detailResourceOption: { detailResourceInfo, detailResourceSelection }
        } = subRoot;

      let { clusterId, rid, resourceIns } = route.queries;
      if (resourceInfo.requestType.useDetailInfo) {
        let response = await WebAPI.fetchResourceYaml(
          detailResourceSelection,
          detailResourceInfo,
          namespaceSelection,
          clusterId,
          +rid
        );
        return response;
      } else {
        let { clusterId, rid, resourceIns } = route.queries;

        let response = await WebAPI.fetchResourceYaml(resourceIns, resourceInfo, namespaceSelection, clusterId, +rid);
        return response;
      }
    }
  }),

  /** 离开详情页，清除detail的详情 */
  clearDetail: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearResourceDetail
    };
  }
};
