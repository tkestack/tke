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

import { extend, ReduxAction, uuid } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { cloneDeep } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import { initSecretData } from '../constants/initState';
import { Namespace, RootState, SecretData } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

/** ================  start fetchNamespaceList =========================== */
const fetchNamespaceActions = generateFetcherActionCreator({
  actionType: ActionType.Sec_FetchNsList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { clusterVersion } = getState();
    // 获取资源的配置
    let namespaceInfo = resourceConfig(clusterVersion)['ns'];
    let response = await WebAPI.fetchNamespaceList(getState().subRoot.secretEdit.nsQuery, namespaceInfo);
    return response;
  }
});

const queryNamespaceActions = generateQueryActionCreator({
  actionType: ActionType.Sec_QueryNsList,
  bindFetcher: fetchNamespaceActions
});

const nsActions = extend(fetchNamespaceActions, queryNamespaceActions);
/** ================  end fetchNamespaceList =========================== */

export const secretEditActions = {
  /** ns的相关操作 */
  ns: nsActions,

  /** 输入secret的名称 */
  inputSecretName: (name: string): ReduxAction<string> => {
    return {
      type: ActionType.Sec_Name,
      payload: name
    };
  },

  /** 初始化 命名空间的选择列表，dockercfg */
  selectNsList: (ns: Namespace[]) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.Sec_NamespaceSelection,
        payload: [...ns]
      });
    };
  },

  /** 选择当前的生效范围 */
  selectNsType: (nsType: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.Sec_NsType,
        payload: nsType
      });
    };
  },

  /** 选择secret的类型 */
  selectSecretType: (type: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.Sec_SecretType,
        payload: type
      });
    };
  },

  /** 增加secret data变量 */
  addSecretData: () => {
    return async (dispatch, getState: GetState) => {
      let { data } = getState().subRoot.secretEdit;
      let newData: SecretData[] = cloneDeep(data);
      newData.push(Object.assign({}, initSecretData, { id: uuid() }));

      dispatch({
        type: ActionType.Sec_UpdateData,
        payload: newData
      });
    };
  },

  /** 删除Secret Data变量 */
  deleteSecretData: (dataId: string) => {
    return async (dispatch, getState: GetState) => {
      let { data } = getState().subRoot.secretEdit,
        newData: SecretData[] = cloneDeep(data),
        dIndex = newData.findIndex(item => item.id === dataId);

      newData.splice(dIndex, 1);
      dispatch({
        type: ActionType.Sec_UpdateData,
        payload: newData
      });
    };
  },

  /** 更新secret data */
  updateSecretData: (obj: any, dataId: string) => {
    return async (dispatch, getState: GetState) => {
      let { data } = getState().subRoot.secretEdit,
        newData: SecretData[] = cloneDeep(data),
        dIndex = newData.findIndex(item => item.id === dataId);

      let objKeys = Object.keys(obj);
      objKeys.forEach(item => {
        newData[dIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.Sec_UpdateData,
        payload: newData
      });
    };
  },

  /** 输入第三方镜像仓库的域名 */
  inputThirdHubDomain: (domain: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.Sec_Domain,
        payload: domain
      });
    };
  },

  /** 输入第三方镜像仓库的用户名 */
  inputThirdHubUserName: (username: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.Sec_Username,
        payload: username
      });
    };
  },

  /** 输入第三方镜像仓库的密码 */
  inputThirdHubPassword: (psw: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.Sec_Password,
        payload: psw
      });
    };
  },

  /** 清空Secret的编辑 */
  clearSecretEdit: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearSecretEdit
    };
  }
};
