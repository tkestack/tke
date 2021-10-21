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
import { createFFListActions, extend } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config';
import { ResourceFilter, ResourceInfo } from '../../common';
import { CommonAPI } from '../../common/webapi';
import { FFReduxActionName } from '../constants/Config';
import { Addon, RootState } from '../models';

type GetState = () => RootState;

/** addon的相关操作 */
const ListAddonActions = createFFListActions<Addon, ResourceFilter>({
  actionName: FFReduxActionName.ADDON,
  fetcher: async (query, getState: GetState) => {
    let { clusterVersion, addon } = getState();
    let addonInfo: ResourceInfo = resourceConfig(clusterVersion)['addon'];
    let response = await CommonAPI.fetchResourceList({ query: addon.query, resourceInfo: addonInfo });

    // 对结果进行排序，保证每次的结果一样，后台是通过promise.all 并行的，返回结果顺序不确定
    response.records = response.records.sort((prev, next) => (prev.type < next.type ? 1 : -1));
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().addon;
  }
});

/** restActions */
const restActions = {};

export const addonActions = extend({}, ListAddonActions, restActions);
