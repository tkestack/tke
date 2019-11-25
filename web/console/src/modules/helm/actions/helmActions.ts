import { extend, ReduxAction } from '@tencent/qcloud-lib';
import { FetchOptions, generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import * as ActionType from '../constants/ActionType';
import { ClusterHelmStatus, HelmResource, InstallingStatus, OtherType } from '../constants/Config';
import { Helm, HelmKeyValue, InstallingHelm, RootState, TencenthubChartVersion } from '../models';
import { HelmListUpdateValid } from '../models/ListState';
import * as WebAPI from '../WebAPI';
import { helm } from 'config/resource/k8sConfig';

type GetState = () => RootState;

const refreshDelay = 5000;

/** 获取列表 */
const fetchHelmActions = generateFetcherActionCreator({
  actionType: ActionType.FetchHelmList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let {
      listState: { helmQuery, region, cluster }
    } = getState();
    let response = await WebAPI.fetchHelmList(helmQuery, region.selection.value, cluster.selection.metadata.name);

    return response;
  },
  finish: (dispatch, getState: GetState) => {
    // let {
    //   listState: { helmList },
    //   route
    // } = getState();
    // if (helmList.data.recordCount) {
    //   dispatch(restActions.select(helmList.data.records[0]));
    // }
  }
});

/** 查询列表action */
const queryHelmActions = generateQueryActionCreator({
  actionType: ActionType.QueryHelmList,
  bindFetcher: fetchHelmActions
});

const restActions = {
  select: (helm: Helm) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectHelm,
        payload: helm
      });

      if (helm.chart_metadata.repo === HelmResource.TencentHub) {
        dispatch(helmActions.fetchTencenthubChartVersionList.fetch());
        const {
          listState: { token }
        } = getState();
        if (!token) {
          dispatch(helmActions.getToken());
        }
      }
    };
  },

  timer4checkClusterHelmStatus: null,
  checkClusterHelmStatus: () => {
    return async (dispatch, getState: GetState) => {
      let {
        listState: { region, cluster }
      } = getState();
      let status = await WebAPI.checkClusterHelmStatus(region.selection.value, cluster.selection.metadata.name);
      dispatch({
        type: ActionType.ClusterHelmStatus,
        payload: status
      });

      if (status.code === ClusterHelmStatus.RUNNING) {
        //如果已经开通了，就查询一下helm应用列表
        dispatch(helmActions.fetch());
        dispatch(helmActions.fetchInstallingHelmList.fetch());
      } else {
        //如果切换到一个未开通的集群，把之前的列表轮询清理掉
        clearTimeout(helmActions.timer4fetchInstallingHelmList);
        helmActions.timer4fetchInstallingHelmList = null;
      }

      if (
        status.code === ClusterHelmStatus.CHECKING ||
        status.code === ClusterHelmStatus.INIT ||
        status.code === ClusterHelmStatus.REINIT
      ) {
        //如果正在开通，轮询一下
        if (helmActions.timer4checkClusterHelmStatus) {
          clearTimeout(helmActions.timer4checkClusterHelmStatus);
          helmActions.timer4checkClusterHelmStatus = null;
        }
        helmActions.timer4checkClusterHelmStatus = setTimeout(() => {
          dispatch(helmActions.checkClusterHelmStatus());
        }, refreshDelay);
      }
    };
  },

  setupHelm: () => {
    return async (dispatch, getState: GetState) => {
      let {
        listState: { region, cluster }
      } = getState();
      await WebAPI.setupHelm(region.selection.value, cluster.selection.metadata.name);
      dispatch(helmActions.checkClusterHelmStatus());
    };
  },

  fetchTencenthubChartVersionList: generateFetcherActionCreator({
    actionType: ActionType.TableFetchTencenthubChartVersionList,
    fetcher: async (getState: GetState, fetchOptions, dispatch) => {
      const {
        listState: { helmSelection }
      } = getState();
      let response = await WebAPI.fetchTencenthubChartVersionList(
        helmSelection.chart_metadata.chart_ns,
        helmSelection.chart_metadata.name
      );
      return response;
    },
    finish: (dispatch, getState: GetState) => {
      let {
        listState: { tencenthubChartVersionList },
        route
      } = getState();
      if (tencenthubChartVersionList.data.recordCount) {
        dispatch(helmActions.selectTencenthubChartVersion(tencenthubChartVersionList.data.records[0]));
      }
    }
  }),
  selectTencenthubChartVersion: (version: TencenthubChartVersion) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.TableTencenthubChartVersionSelection,
        payload: version
      });
    };
  },
  getToken: () => {
    return async (dispatch, getState: GetState) => {
      let token = await WebAPI.getTencentHubToken();
      dispatch({
        type: ActionType.TableTencenthubToken,
        payload: token
      });
    };
  },
  updateHelm: () => {
    return async (dispatch, getState: GetState) => {
      let {
        listState: {
          helmSelection,
          tencenthubChartVersionSelection,
          region,
          cluster,
          otherChartUrl,
          otherUserName,
          otherPassword,
          otherTypeSelection,
          kvs,
          token
        }
      } = getState();
      // await WebAPI.deleteHelm(helm, regionSelection.value, clusterSelection.clusterId);
      // dispatch(
      //     helmActions.fetch()
      // );
      for (let i = 0; i < kvs.length; i++) {
        kvs[i].value = kvs[i].value.replace(/,/g, '\\,');
        if (!kvs[i].key || !kvs[i].key.trim()) {
          kvs.splice(i--, 1);
        }
      }
      if (helmSelection.chart_metadata.repo === HelmResource.TencentHub) {
        let data = {
          helmName: helmSelection.name,
          chart_url: tencenthubChartVersionSelection.download_url,
          token: token
        };

        if (kvs && kvs.length) {
          data['kvs'] = kvs;
        }
        await WebAPI.updateHelm(data, region.selection.value, cluster.selection.metadata.name);
      } else {
        let data = {
          helmName: helmSelection.name,
          chart_url: otherChartUrl
        };
        if (otherTypeSelection === OtherType.Private) {
          data['username'] = otherUserName;
          data['password'] = otherPassword;
        }
        if (kvs && kvs.length) {
          data['kvs'] = kvs;
          /// #if project
          kvs.push({ key: 'NAMESPACE', value: getState().namespaceSelection });
          /// #endif
        }
        await WebAPI.updateHelmByOther(data, region.selection.value, cluster.selection.metadata.name);
      }

      dispatch(helmActions.fetchInstallingHelmList.fetch());
      dispatch(helmActions.fetch());
    };
  },
  delete: (helm: Helm) => {
    return async (dispatch, getState: GetState) => {
      let {
        listState: { region, cluster }
      } = getState();
      await WebAPI.deleteHelm(
        {
          helmName: helm.name
        },
        region.selection.value,
        cluster.selection.metadata.name
      );
      dispatch(helmActions.fetch());
    };
  },

  timer4fetchInstallingHelmList: null,
  fetchInstallingHelmList: generateFetcherActionCreator({
    actionType: ActionType.FetchInstallingHelmList,
    fetcher: async (getState: GetState, fetchOptions, dispatch) => {
      let {
        listState: { region, cluster }
      } = getState();
      let list = await WebAPI.fetchInstallingHelmList(region.selection.value, cluster.selection.metadata.name);
      return list;
    },
    finish: (dispatch, getState: GetState) => {
      let {
        listState: { installingHelmList },
        route
      } = getState();
      if (installingHelmList.data.recordCount) {
        // dispatch(
        //   restActions.selectInstallingHelm(installingHelmList.data.records[0])
        // );

        //如果有正在安装的应用，轮询一下
        if (installingHelmList.data.records.find(item => item.status === InstallingStatus.INSTALLING)) {
          if (helmActions.timer4fetchInstallingHelmList) {
            clearTimeout(helmActions.timer4fetchInstallingHelmList);
            helmActions.timer4fetchInstallingHelmList = null;
          }
          helmActions.timer4fetchInstallingHelmList = setTimeout(() => {
            dispatch(helmActions.fetch());
            dispatch(helmActions.fetchInstallingHelmList.fetch());
          }, refreshDelay);
        }
      }
    }
  }),
  selectInstallingHelm: (helm: InstallingHelm) => {
    return async (dispatch, getState: GetState) => {
      if (helm) {
        dispatch({
          type: ActionType.SelectInstallingHelm,
          payload: helm
        });
        dispatch(helmActions.fetchInstallingHelm(helm.name));
      } else {
        dispatch({
          type: ActionType.FetchInstallingHelm,
          payload: {
            code: 0,
            message: ''
          }
        });
        dispatch({
          type: ActionType.SelectInstallingHelm,
          payload: null
        });
      }
    };
  },
  fetchInstallingHelm: (helmName: string) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState();

      let helm = await WebAPI.fetchInstallingHelm({ helmName }, +route.queries['rid'], route.queries['clusterId']);

      dispatch({
        type: ActionType.FetchInstallingHelm,
        payload: helm
      });
    };
  },
  ignoreInstallingHelm: (helm: InstallingHelm) => {
    return async (dispatch, getState: GetState) => {
      let {
        listState: { region, cluster, installingHelmSelection }
      } = getState();
      await WebAPI.ignoreInstallingHelm(
        {
          helmName: helm.name
        },
        region.selection.value,
        cluster.selection.metadata.name
      );

      if (installingHelmSelection && installingHelmSelection.name === helm.name) {
        //如果删除了当前选中的helm，应该把右侧的内容清理一下
        dispatch(helmActions.selectInstallingHelm(null));
      }

      dispatch(helmActions.fetchInstallingHelmList.fetch());
    };
  },

  inputOtherChartUrl: (chart_url: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.ListOtherChartUrl,
        payload: chart_url
      });
      dispatch(helmActions.validOtherChartUrl());
    };
  },

  validOtherChartUrl: () => {
    return async (dispatch, getState: GetState) => {
      let {
        listState: { otherChartUrl, isValid }
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
      dispatch(helmActions.setIsValid(isValid));
    };
  },

  selectOtherType: (type: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.ListOtherType,
        payload: type
      });
    };
  },
  inputOtherUserName: (username: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.ListOtherUserName,
        payload: username
      });
      dispatch(helmActions.validOtherUserName());
    };
  },

  validOtherUserName: () => {
    return async (dispatch, getState: GetState) => {
      let {
        listState: { otherUserName, isValid }
      } = getState();
      if (otherUserName === '') {
        isValid.otherUserName = t('请输入用户名');
      } else {
        isValid.otherUserName = '';
      }
      dispatch(helmActions.setIsValid(isValid));
    };
  },
  inputOtherPassword: (password: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.ListOtherPassword,
        payload: password
      });
      dispatch(helmActions.validOtherPassword());
    };
  },

  validOtherPassword: () => {
    return async (dispatch, getState: GetState) => {
      let {
        listState: { otherPassword, isValid }
      } = getState();
      if (otherPassword === '') {
        isValid.otherPassword = t('请输入密码');
      } else {
        isValid.otherPassword = '';
      }
      dispatch(helmActions.setIsValid(isValid));
    };
  },

  inputKeyValue: (kvs: HelmKeyValue[]) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.ListKeyValue,
        payload: kvs
      });
    };
  },

  setIsValid: (isValid: HelmListUpdateValid) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.ListUpdateIsValid,
        payload: Object.assign({}, isValid)
      });
    };
  },
  validAll: () => {
    return async (dispatch, getState: GetState) => {
      dispatch(helmActions.validOtherChartUrl());
      dispatch(helmActions.validOtherUserName());
      dispatch(helmActions.validOtherPassword());
    };
  }
};

export const helmActions = extend({}, fetchHelmActions, queryHelmActions, restActions);
