import {
    extend, FetchOptions, generateFetcherActionCreator, ReduxAction, uuid
} from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config/resourceConfig';
import { cloneDeep } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import { initImagePullSecrets } from '../constants/initState';
import { ImagePullSecrets, Resource, ResourceFilter, RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

const fetchSecretActions = generateFetcherActionCreator({
  actionType: ActionType.W_FetchSecret,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { workloadEdit } = subRoot,
      { configEdit } = workloadEdit;

    let secretResourceInfo = resourceConfig(clusterVersion)['secret'];

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await WebAPI.fetchSpecificResourceList(
      configEdit.secretQuery,
      secretResourceInfo,
      isClearData,
      true
    );
    return response;
  },
  finish: (dispatch, getState: GetState) => {
    let { secretList } = getState().subRoot.workloadEdit.configEdit;
    dispatch(workloadSecretActions.initImagePullSecrets(secretList.data.records));
  }
});

const querySecretActions = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.W_QuerySecret,
  bindFetcher: fetchSecretActions
});

const restActions = {
  /** 这里去初始化imagePullSecrets的列表 */
  initImagePullSecrets: (secretList: Resource[]) => {
    return async (dispatch, getState: GetState) => {
      let imagePullSecrets: ImagePullSecrets[];

      dispatch({
        type: ActionType.ImagePullSecrets,
        payload: imagePullSecrets || []
      });
    };
  },

  /** 新增imagePull */
  addImagePullSecret: () => {
    return async (dispatch, getState: GetState) => {
      let newList = cloneDeep(getState().subRoot.workloadEdit.imagePullSecrets);
      let newImagePullSecret: ImagePullSecrets = Object.assign({}, initImagePullSecrets, {
        id: uuid()
      });
      newList.push(newImagePullSecret);
      dispatch({
        type: ActionType.ImagePullSecrets,
        payload: newList
      });
    };
  },

  /** 删除imgaePull */
  deleteImagePullSecret: (sId: string) => {
    return async (dispatch, getState: GetState) => {
      let newList: ImagePullSecrets[] = cloneDeep(getState().subRoot.workloadEdit.imagePullSecrets),
        sIndex = newList.findIndex(item => item.id === sId);

      newList.splice(sIndex, 1);
      dispatch({
        type: ActionType.ImagePullSecrets,
        payload: newList
      });
    };
  },

  /** 选择imagePullSecret的操作 */
  updateImagePullSecret: (obj: any = {}, sId: string) => {
    return async (dispatch, getState: GetState) => {
      let newList: ImagePullSecrets[] = cloneDeep(getState().subRoot.workloadEdit.imagePullSecrets),
        sIndex = newList.findIndex(item => item.id === sId);
      let objKeys = Object.keys(obj);
      objKeys.forEach(item => {
        newList[sIndex][item] = obj[item];
      });

      dispatch({
        type: ActionType.ImagePullSecrets,
        payload: newList
      });
    };
  }
};

export const workloadSecretActions = extend({}, fetchSecretActions, querySecretActions, restActions);
