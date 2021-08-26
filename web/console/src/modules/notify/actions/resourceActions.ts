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

import { createFFListActions } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config';
import { ResourceInfo } from '../../common/models';
import { Resource, ResourceFilter, RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
let rc = resourceConfig();

export const resourceActions = {
  channel: createFFListActionsFactory('channel'),
  template: createFFListActionsFactory('template'),
  receiver: createFFListActionsFactory('receiver'),
  receiverGroup: createFFListActionsFactory('receiverGroup')
};

function createFFListActionsFactory(resourceName) {
  return createFFListActions<Resource, ResourceFilter>({
    actionName: resourceName,
    fetcher: async (query, getState: GetState, fetchOptions) => {
      let resourceInfo: ResourceInfo = rc[resourceName];
      let resourceItems = await WebAPI.fetchResourceList(query, resourceInfo);

      let { route, receiverGroup } = getState();
      let urlParams = router.resolve(route);
      if (resourceName === 'receiver' && urlParams.resourceName === 'receiverGroup' && urlParams.mode === 'detail') {
        let rg = receiverGroup.list.data.records.find(rg => rg.metadata.name === route.queries.resourceIns);
        if (rg) {
          resourceItems.records = resourceItems.records.filter(item =>
            rg.spec.receivers.find(r => r === item.metadata.name)
          );
        }
      }

      if (resourceName === 'channel') {
        resourceItems.records = resourceItems.records.filter(item => item.status.phase !== 'Terminating');
      }

      resourceItems.recordCount = resourceItems.records.length;
      let response = resourceItems;

      // 告警编辑页里有receiver group 根据id过滤，选中对应id
      if (fetchOptions && fetchOptions.data) {
        fetchOptions.data.forEach(item => {
          let finder = response.records.find(group => group.metadata.name === item);
          finder && (finder.selected = true);
        });
      }
      return response;
    },
    getRecord: (getState: GetState) => {
      return getState()[resourceName];
    },
    onFinish: (record, dispatch) => {
      let selects = record.data.records.filter(r => r.selected);
      // 告警编辑页里有receiver group 根据id过滤，选中对应id
      if (selects) {
        dispatch(resourceActions[resourceName].selects(selects));
      }
    }
  });
}
