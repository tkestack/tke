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

import { extend, createFFObjectActions, uuid } from '@tencent/ff-redux';
import * as JsYAML from 'js-yaml';
import { RootState, AppResource, AppResourceFilter, Resource } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
type GetState = () => RootState;
const tips = seajs.require('tips');

/**
 * 获取资源详情
 */

const fetchAppResourceActions = createFFObjectActions<AppResource, AppResourceFilter>({
  actionName: ActionTypes.AppResource,
  fetcher: async (query, getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let response = await WebAPI.fetchAppResource(query.filter);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().appResource;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let resources = new Map<string, Resource[]>();
    if (record.data) {
      try {
        Object.keys(record.data.spec.resources).forEach(k => {
          record.data.spec.resources[k].forEach(item => {
            let json = JsYAML.safeLoad(item);
            if (!resources.get(k)) {
              resources.set(k, []);
            }
            resources.get(k).push({
              id: uuid(),
              metadata: {
                namespace: json.metadata.namespace,
                name: json.metadata.name
              },
              kind: json.kind,
              cluster: record.data.spec.targetCluster,
              yaml: JsYAML.safeDump(json),
              object: json
            });
          });
        });
      } catch (e) {
        // console.log(e);
        tips.error(t('资源列表读取失败'), 2000);
      }
    }
    dispatch({
      type: ActionTypes.ResourceList,
      payload: {
        resources: resources
      }
    });
  }
});

const restActions = {
  /** 轮询操作 */
  poll: (filter: AppResourceFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(
        resourceActions.polling({
          delayTime: 5000,
          filter: filter
        })
      );
    };
  }
};

export const resourceActions = extend({}, fetchAppResourceActions, restActions);
