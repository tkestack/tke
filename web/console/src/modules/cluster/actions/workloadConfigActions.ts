import { uuid, extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { RootState, ResourceFilter, Resource, ConfigItems } from '../models';
import * as WebAPI from '../WebAPI';
import { cloneDeep } from '../../common/utils';
import { initConfigMapItem } from '../constants/initState';
import { resourceConfig } from '../../../../config/resourceConfig';
import * as ActionType from '../constants/ActionType';
import { workloadEditActions } from './workloadEditActions';

type GetState = () => RootState;

const fetchConfigActions = generateFetcherActionCreator({
  actionType: ActionType.W_FetchConfig,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { workloadEdit } = subRoot,
      { configEdit } = workloadEdit;

    let configMapResourceInfo = resourceConfig(clusterVersion)['configmap'];

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await WebAPI.fetchSpecificResourceList(
      configEdit.configQuery,
      configMapResourceInfo,
      isClearData,
      true
    );
    return response;
  },
  finish: (dispatch, getState: GetState) => {
    let { configList } = getState().subRoot.workloadEdit.configEdit;

    configList.data.recordCount && dispatch(workloadConfigActions.selectConfig([configList.data.records[0]]));
  }
});

const queryConfigActions = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.W_QueryConfig,
  bindFetcher: fetchConfigActions
});

const configRestActions = {
  /** 选择具体的config */
  selectConfig: (configMaps: Resource[]) => {
    return async (dispatch, getState: GetState) => {
      let { keyType } = getState().subRoot.workloadEdit.configEdit;

      dispatch({
        type: ActionType.W_ConfigSelect,
        payload: configMaps
      });

      // 这里去初始化configItems
      dispatch(workloadEditActions.config.initConfigItems(configMaps, keyType, true));
    };
  },

  /** 初始化configMapItems */
  initConfigItems: (configMaps: Resource[], keyType: string, isNeedUpdateConfigItems: boolean = false) => {
    return async (dispatch, getState: GetState) => {
      let dataKeys = configMaps.length && configMaps[0].data ? Object.keys(configMaps[0].data) : [];
      let configMapItems: ConfigItems[] = [];

      dispatch({
        type: ActionType.W_UpdateConfigItems,
        payload: configMapItems
      });

      // 更新configKeys
      isNeedUpdateConfigItems &&
        dispatch({
          type: ActionType.W_UpdateConfigKeys,
          payload: dataKeys
        });

      // 如果keyType为optional，默认加入一个key
      if (keyType === 'optional' && dataKeys.length) {
        dispatch(workloadConfigActions.addConfigItem());
      }
    };
  },

  /** initByVolumesConfigKey */
  initConfigItemsByVolumes: (configs: ConfigItems[]) => {
    return async (dispatch, getState: GetState) => {
      let configItems: ConfigItems[] = [];

      configs.forEach(item => {
        let tmp: ConfigItems = Object.assign({}, initConfigMapItem, {
          id: uuid(),
          configKey: item.configKey,
          path: item.path
        });
        configItems.push(tmp);
      });

      dispatch({
        type: ActionType.W_UpdateConfigItems,
        payload: configItems
      });
    };
  },

  /** 更新 configItems */
  updateConfigItems: (obj: any, cId: string) => {
    return async (dispatch, getState: GetState) => {
      let configItems: ConfigItems[] = cloneDeep(getState().subRoot.workloadEdit.configEdit.configItems),
        cIndex = configItems.findIndex(item => item.id === cId);

      let itemKeys = Object.keys(obj);
      itemKeys.forEach(key => {
        configItems[cIndex][key] = obj[key];
      });

      dispatch({
        type: ActionType.W_UpdateConfigItems,
        payload: configItems
      });
    };
  },

  /** 新增configItem */
  addConfigItem: () => {
    return async (dispatch, getState: GetState) => {
      let { configItems, configKeys } = getState().subRoot.workloadEdit.configEdit;
      let newConfigItems = cloneDeep(configItems);

      newConfigItems.push(
        Object.assign({}, initConfigMapItem, {
          id: uuid(),
          configKey: configKeys[newConfigItems.length]
        })
      );

      dispatch({
        type: ActionType.W_UpdateConfigItems,
        payload: newConfigItems
      });
    };
  },

  /** 删除configItem */
  deleteConfigMapItem: (cId: string) => {
    return async (dispatch, getState: GetState) => {
      let configMapItems: ConfigItems[] = cloneDeep(getState().subRoot.workloadEdit.configEdit.configItems),
        cIndex = configMapItems.findIndex(c => c.id === cId);

      configMapItems.splice(cIndex, 1);
      dispatch({
        type: ActionType.W_UpdateConfigItems,
        payload: configMapItems
      });
    };
  },

  /** 变换 configMap的key的选项 */
  changeKeyType: (type: string) => {
    return async (dispatch, getState: GetState) => {
      let { configSelection } = getState().subRoot.workloadEdit.configEdit;

      dispatch({
        type: ActionType.W_ChangeKeyType,
        payload: type
      });

      dispatch(workloadEditActions.config.initConfigItems(configSelection, type));
    };
  }
};

export const workloadConfigActions = extend({}, fetchConfigActions, queryConfigActions, configRestActions);
