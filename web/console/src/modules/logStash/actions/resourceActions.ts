import { extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { resourceConfig } from '../../../../config';
import { CommonAPI } from '../../../../src/modules/common';
import { ResourceInfo } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { ContainerLogs, ResourceFilter, RootState } from '../models';
import { WorkLoadList } from '../models/Resource';
import { editLogStashActions } from './editLogStashActions';

type GetState = () => RootState;

/** 获取workloadList */
const fetchWorkloadList = generateFetcherActionCreator({
  actionType: ActionType.FetchResourceList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { logStashEdit, clusterVersion } = getState(),
      { resourceQuery } = logStashEdit;

    let { workloadType, namespace } = resourceQuery.filter;
    let resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[workloadType];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    if (namespace === '' || workloadType === '') {
      isClearData = true;
    }
    let response = await CommonAPI.fetchResourceList({
      query: resourceQuery,
      resourceInfo,
      isClearData
    });
    return response;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    //在container模式下或者containerFile模式下都需要获取，需要分情况处理获取之后的结果
    const { isForContainerFile, isForContainerLogs } = getState().logStashEdit.resourceTarget;
    const { isFirstFetchResource } = getState().logStashEdit;
    if (isForContainerLogs || isFirstFetchResource) {
      let { resourceQuery, resourceList, containerLogs } = getState().logStashEdit,
        { workloadType } = resourceQuery.filter;

      let cIndex = containerLogs.findIndex(item => item.status === 'editing');
      if (cIndex !== -1) {
        // 将拉取的列表更新到对应的containerLogs当中
        let containerLog: ContainerLogs = containerLogs[cIndex];
        let workloadList = Object.assign({}, containerLog.workloadList, {
          [workloadType]: resourceList.data.records
        });
        let workloadListFetch = Object.assign({}, containerLog.workloadListFetch, {
          [workloadType]: true
        });
        dispatch(editLogStashActions.updateContainerLog({ workloadList, workloadListFetch }, cIndex));
      }
    }

    if (isForContainerFile || isFirstFetchResource) {
      const { resourceList } = getState().logStashEdit;
      let workloadList: WorkLoadList[] = resourceList.data.records.map(item => {
        return {
          name: item.metadata.name,
          value: item.metadata.name
        };
      });
      dispatch(editLogStashActions.updateContainerFileWorkloadList(workloadList));

      //默认选择containerFile下的workload
      const defaultWorkload = workloadList.length > 0 ? workloadList[0].value : '';
      dispatch(editLogStashActions.selectContainerFileWorkload(defaultWorkload));
    }

    //已經不是第一次獲取資源對象，需要更新字段
    if (isFirstFetchResource) {
      dispatch(editLogStashActions.ifFirstFetchResource(false));
    }
  }
});

/** 获取workloadList的查询 */
const queryWorkloadList = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.QueryResourceList,
  bindFetcher: fetchWorkloadList
});

const restActions = {};

export const resourceActions = extend({}, fetchWorkloadList, queryWorkloadList, restActions);
