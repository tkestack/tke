import { createFFListActions, extend, ReduxAction } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { t } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../config';
import { assureRegion } from '../../../../helpers';
import { Region, RegionFilter, Resource, ResourceFilter, ResourceInfo } from '../../common/models';
import { CommonAPI } from '../../common/webapi';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName, HelmResource, OtherType, TencentHubType } from '../constants/Config';
import {
    HelmCreationValid, HelmKeyValue, RootState, TencenthubChart, TencenthubChartVersion
} from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { helmActions } from './helmActions';

type GetState = () => RootState;

export const regionActions = createFFListActions<Region, RegionFilter>({
  id: 'HelmCreate',
  actionName: FFReduxActionName.REGION,
  selectFirst: true,
  fetcher: async query => {
    let response = await CommonAPI.fetchRegionList(query);
    return response;
  },
  onFinish: (record, dispatch, getState: GetState) => {
    let {
      helmCreation: { region },
      route
    } = getState();
    if (record.data.recordCount) {
      let defaultRegion = assureRegion(record.data.records, route.queries['rid'], 1);
      let rs = region.list.data.records.find(r => r.value + '' === defaultRegion + '');
      rs && dispatch(regionActions.select(rs));
    }
  },
  getRecord: (getState: GetState) => {
    return getState().helmCreation.region;
  },
  onSelect: (region: Region, dispatch: Redux.Dispatch) => {
    dispatch(clusterActions.applyFilter({ regionId: region.value }));
  }
});

/** 获取集群列表 */
const ListClusterActions = createFFListActions<Resource, ResourceFilter>({
  id: 'HelmCreate',
  actionName: FFReduxActionName.CLUSTER,
  fetcher: async query => {
    let clusterInfo: ResourceInfo = resourceConfig()['cluster'];
    let response = await CommonAPI.fetchResourceList({ query, resourceInfo: clusterInfo });
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().helmCreation.cluster;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState();
    if (record.data.recordCount) {
      let defaultClusterId = route.queries['clusterId'];
      let defaultCluster = record.data.records.find(item => item.metadata.name === defaultClusterId);
      dispatch(clusterActions.selectCluster(defaultCluster || record.data.records[0]));
    }
  }
});

const restClusterActions = {
  selectCluster: (cluster: Resource) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);
      dispatch(clusterActions.select(cluster));
      router.navigate(urlParams, Object.assign({}, route.queries, { clusterId: cluster.metadata.name }));
    };
  }
};

const clusterActions = extend({}, ListClusterActions, restClusterActions);

export const createActions = {
  clear: (): ReduxAction<void> => {
    return { type: ActionType.ClearCreation };
  },
  fetchRegionList: () => {
    return async (dispatch, getState: GetState) => {
      dispatch(regionActions.fetch());
    };
  },

  createHelm: () => {
    return async (dispatch, getState: GetState) => {
      const {
        region,
        cluster,
        name,
        token,
        resourceSelection,
        tencenthubNamespaceSelection,
        tencenthubChartVersionSelection,
        otherChartUrl,
        otherTypeSelection,
        otherUserName,
        otherPassword,
        kvs
      } = getState().helmCreation;

      for (let i = 0; i < kvs.length; i++) {
        kvs[i].value = kvs[i].value.replace(/,/g, '\\,');
        if (!kvs[i].key || !kvs[i].key.trim()) {
          kvs.splice(i--, 1);
        }
      }

      let clusterSelection = cluster.selection.metadata.name;

      if (resourceSelection === HelmResource.TencentHub) {
        let response = await WebAPI.createHelm(
          {
            helmName: name,
            resource: resourceSelection,
            namespace: tencenthubNamespaceSelection,
            chart_url: tencenthubChartVersionSelection.download_url,
            password: token,
            kvs
          },
          region.selection.value,
          clusterSelection
        );
      } else if (resourceSelection === HelmResource.Other) {
        let options = {
          helmName: name,
          resource: resourceSelection,
          chart_url: otherChartUrl,
          kvs
        };
        if (otherTypeSelection === OtherType.Private) {
          options['username'] = otherUserName;
          options['password'] = otherPassword;
        }
        /// #if project
        kvs.push({ key: 'NAMESPACE', value: getState().namespaceSelection });
        /// #endif
        let response = await WebAPI.createHelmByOther(options, region.selection.value, clusterSelection);
      }

      dispatch(helmActions.fetch());
      dispatch(helmActions.fetchInstallingHelmList.fetch());

      router.navigate(
        {},
        {
          rid: region.selection.value + '',
          clusterId: clusterSelection
        }
      );
    };
  },

  inputName: (name: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.C_CreateionName,
        payload: name
      });
      dispatch(createActions.validName());
    };
  },
  validName: () => {
    return async (dispatch, getState: GetState) => {
      let {
        helmCreation: { name, isValid }
      } = getState();
      if (name === '') {
        isValid.name = t('请输入应用名');
      } else {
        if (!/^[a-z]([a-z0-9-]{0,61}[0-9a-z])?$/.test(name)) {
          isValid.name = t('应用名格式不正确');
        } else {
          isValid.name = '';
        }
      }
      dispatch(createActions.setIsValid(isValid));
    };
  },
  setIsValid: (isValid: HelmCreationValid) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.IsValid,
        payload: Object.assign({}, isValid)
      });
    };
  },
  getToken: () => {
    return async (dispatch, getState: GetState) => {
      let token = await WebAPI.getTencentHubToken();
      dispatch({
        type: ActionType.TencenthubToken,
        payload: token
      });
    };
  },
  selectResource: (resource: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.ResourceSelection,
        payload: resource
      });
    };
  },

  selectTencenthubType: (type: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.TencenthubTypeSelection,
        payload: type
      });
      if (type === TencentHubType.Private) {
        dispatch(createActions.fetchTencenthubNamespaceList.fetch());
      } else {
        dispatch(createActions.selectTencenthubNamespace('tencenthub'));
      }
    };
  },
  fetchTencenthubNamespaceList: generateFetcherActionCreator({
    actionType: ActionType.FetchTencenthubNamespaceList,
    fetcher: async (getState: GetState, fetchOptions, dispatch) => {
      let response = await WebAPI.fetchTencenthubNamespaceList();
      return response;
    },
    finish: (dispatch, getState: GetState) => {
      let {
        helmCreation: { tencenthubNamespaceList },
        route
      } = getState();
      if (tencenthubNamespaceList.data.recordCount) {
        dispatch(createActions.selectTencenthubNamespace(tencenthubNamespaceList.data.records[0].name));
      }
    }
  }),
  selectTencenthubNamespace: (namespace: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.TencenthubNamespaceSelection,
        payload: namespace
      });

      dispatch({
        type: ActionType.TencenthubChartReadMe,
        payload: ''
      });
      dispatch(createActions.fetchTencenthubChartList.fetch());
    };
  },
  fetchTencenthubChartList: generateFetcherActionCreator({
    actionType: ActionType.FetchTencenthubChartList,
    fetcher: async (getState: GetState, fetchOptions, dispatch) => {
      const {
        helmCreation: { tencenthubNamespaceSelection }
      } = getState();
      let response = await WebAPI.fetchTencenthubChartList(tencenthubNamespaceSelection);
      return response;
    },
    finish: (dispatch, getState: GetState) => {
      let {
        helmCreation: { tencenthubChartList },
        route
      } = getState();
      if (tencenthubChartList.data.recordCount) {
        dispatch(createActions.selectTencenthubChart(tencenthubChartList.data.records[0]));
      }
    }
  }),
  selectTencenthubChart: (chart: TencenthubChart) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.TencenthubChartSelection,
        payload: chart
      });
      dispatch(
        createActions.fetchTencenthubChartVersionList.update({
          recordCount: 0,
          records: []
        })
      );
      dispatch(createActions.fetchTencenthubChartVersionList.fetch());
    };
  },
  fetchTencenthubChartVersionList: generateFetcherActionCreator({
    actionType: ActionType.FetchTencenthubChartVersionList,
    fetcher: async (getState: GetState, fetchOptions, dispatch) => {
      const {
        helmCreation: { tencenthubNamespaceSelection, tencenthubChartSelection }
      } = getState();
      let response = await WebAPI.fetchTencenthubChartVersionList(
        tencenthubNamespaceSelection,
        tencenthubChartSelection.name
      );
      return response;
    },
    finish: (dispatch, getState: GetState) => {
      let {
        helmCreation: { tencenthubChartVersionList },
        route
      } = getState();
      if (tencenthubChartVersionList.data.recordCount) {
        dispatch(createActions.selectTencenthubChartVersion(tencenthubChartVersionList.data.records[0]));
      }
    }
  }),
  selectTencenthubChartVersion: (version: TencenthubChartVersion) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.TencenthubChartVersionSelection,
        payload: version
      });
      dispatch({
        type: ActionType.TencenthubChartReadMe,
        payload: null
      });
      dispatch(createActions.fetchTencenthubChartReadMe());
    };
  },
  fetchTencenthubChartReadMe: () => {
    return async (dispatch, getState: GetState) => {
      const {
        helmCreation: { tencenthubNamespaceSelection, tencenthubChartSelection, tencenthubChartVersionSelection }
      } = getState();
      let readme = await WebAPI.fetchTencenthubChartReadMe(
        tencenthubNamespaceSelection,
        tencenthubChartSelection.name,
        tencenthubChartVersionSelection.version
      );
      dispatch({
        type: ActionType.TencenthubChartReadMe,
        payload: readme
      });
    };
  },
  inputOtherChartUrl: (chart_url: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.OtherChartUrl,
        payload: chart_url
      });
      dispatch(createActions.validOtherChartUrl());
    };
  },
  validOtherChartUrl: () => {
    return async (dispatch, getState: GetState) => {
      let {
        helmCreation: { otherChartUrl, isValid }
      } = getState();

      if (otherChartUrl === '') {
        isValid.otherChartUrl = t('请输入Chart_url');
      } else {
        if (!/^http[\d\D]*tgz$/.test(otherChartUrl)) {
          isValid.otherChartUrl = t('Chart_url格式不正确');
        } else {
          isValid.otherChartUrl = '';
        }
      }
      dispatch(createActions.setIsValid(isValid));
    };
  },
  selectOtherType: (type: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.OtherType,
        payload: type
      });
    };
  },
  inputOtherUserName: (username: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.OtherUserName,
        payload: username
      });
      dispatch(createActions.validOtherUserName());
    };
  },
  validOtherUserName: () => {
    return async (dispatch, getState: GetState) => {
      let {
        helmCreation: { otherUserName, isValid }
      } = getState();
      if (otherUserName === '') {
        isValid.otherUserName = t('请输入用户名');
      } else {
        isValid.otherUserName = '';
      }
      dispatch(createActions.setIsValid(isValid));
    };
  },
  inputOtherPassword: (password: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.OtherPassword,
        payload: password
      });
      dispatch(createActions.validOtherPassword());
    };
  },
  validOtherPassword: () => {
    return async (dispatch, getState: GetState) => {
      let {
        helmCreation: { otherPassword, isValid }
      } = getState();
      if (otherPassword === '') {
        isValid.otherPassword = t('请输入密码');
      } else {
        isValid.otherPassword = '';
      }
      dispatch(createActions.setIsValid(isValid));
    };
  },
  inputKeyValue: (kvs: HelmKeyValue[]) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.KeyValue,
        payload: kvs
      });
    };
  },
  validAll: () => {
    return async (dispatch, getState: GetState) => {
      dispatch(createActions.validName());
      dispatch(createActions.validOtherChartUrl());
      dispatch(createActions.validOtherUserName());
      dispatch(createActions.validOtherPassword());
    };
  }
};
