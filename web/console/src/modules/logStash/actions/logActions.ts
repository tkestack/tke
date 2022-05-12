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
import { extend, generateQueryActionCreator, QueryState, RecordSet, ReduxAction, uuid } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';

import { resourceConfig } from '../../../../config';
import { CommonAPI } from '../../common';
import { ResourceInfo } from '../../common/models';
import { cloneDeep } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import { inputTypeMap, outputTypeMap } from '../constants/Config';
import { initContainerFilePath, initContainerInputOption, initMetadata } from '../constants/initState';
import { ContainerLogs, Log, LogFilter, ResourceFilter, RootState } from '../models';
import { Resource } from '../models/Resource';
import { editLogStashActions } from './editLogStashActions';
import { podActions } from './podActions';
import { resourceActions } from './resourceActions';
import { Base64 } from 'js-base64';
import { HOST_LOG_INPUT_PATH_PREFIX } from '../constants/Config';

type GetState = () => RootState;

/** 获取Log采集器的列表的Action */
const fetchLogListActions = generateFetcherActionCreator({
  actionType: ActionType.FetchLogList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    const { logQuery, clusterVersion } = getState();
    const resourceInfo = resourceConfig(clusterVersion)['logcs'];
    const isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    const response = await CommonAPI.fetchResourceList({ resourceInfo, query: logQuery, isClearData });
    return response;
  }
});

/** 查询log采集器列表的Action */
const queryLogListActions = generateQueryActionCreator<LogFilter>({
  actionType: ActionType.QueryLogList,
  bindFetcher: fetchLogListActions
});

export const restActions = {
  /** 选择某个具体的日志采集器 */
  selectLogStash: (log: Log[]): ReduxAction<Log[]> => {
    return {
      type: ActionType.SelectLog,
      payload: log
    };
  },

  /** 拉取单个日志采集规则的 */
  fetchSpecificLog: (name: string, clusterId: string, namespace: string, mode: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { clusterVersion, route, clusterSelection } = getState();
      const logAgentName = (clusterSelection && clusterSelection[0] && clusterSelection[0].spec.logAgentName) || '';
      const resourceInfo = resourceConfig(clusterVersion)['logcs'];
      const { clusterId, regionId } = route.queries;
      const result = await CommonAPI.fetchResourceList({
        query: {
          filter: {
            namespace,
            clusterId,
            logAgentName,
            specificName: name
          }
        },
        resourceInfo
      });
      dispatch({
        type: ActionType.IsFetchDoneSpecificLog,
        payload: true
      });
      dispatch({
        type: ActionType.SelectLog,
        payload: result.records
      });

      const log = getState().logSelection[0];
      if (mode === 'update' || mode === 'detail') {
        let inputOption; //输入端选项
        let outputOption; //输出端选项
        const inputMode = log.spec.input.type ? inputTypeMap[log.spec.input.type] : ''; //输入端类型
        const consumerMode = log.spec.output.type ? outputTypeMap[log.spec.output.type] : ''; //输出端类型
        dispatch(editLogStashActions.inputStashName(log.metadata.name));
        dispatch(editLogStashActions.changeLogMode(inputMode)); //选择mode
        dispatch(editLogStashActions.changeConsumerMode(consumerMode));
        dispatch({
          type: ActionType.isFirstFetchResource,
          payload: false
        });
        //处理输入端选项

        //如果为容器标准模式输出
        if (inputMode === inputTypeMap['container-log']) {
          //如果选择了所有的容器
          inputOption = log.spec.input.container_log_input;
          if (inputOption.all_namespaces) {
            dispatch(editLogStashActions.selectAllNamespace('selectAll'));
            //选择了所有容器，则需要帮忙选择指定容器业务,
            //只有在update的时候才需要去获取resource、
            if (mode === 'update') {
              const containerLogsArr: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
              containerLogsArr[0].namespaceSelection = 'default';
              dispatch({
                type: ActionType.UpdateContainerLogs,
                payload: containerLogsArr
              });
              //进行资源的拉取
              dispatch(
                resourceActions.applyFilter({
                  clusterId: route.queries['clusterId'],
                  namespace: 'default',
                  workloadType: 'deployment',
                  regionId: +route.queries['rid']
                })
              );
            }
          } else {
            //选择了指定容器
            dispatch(editLogStashActions.selectAllNamespace('selectOne'));
            const containerArr = [];
            inputOption.namespaces.forEach(item => {
              const tmp = cloneDeep(initContainerInputOption);
              tmp.id = uuid();
              tmp.collectorWay = item.all_containers ? 'container' : 'workload';
              tmp.namespaceSelection = item.namespace;
              tmp.status = 'edited';
              if (!item.all_containers) {
                //选择了指定容器 兼容旧版本日志
                if (item.workloads) {
                  item.workloads.forEach(item => {
                    tmp.workloadSelection[item.type].push(item.name);
                  });
                } else if (item.services) {
                  item.services.forEach(service => {
                    tmp.workloadSelection['deployment'].push(service);
                  });
                }
              }
              containerArr.push(tmp);
            });
            //只有在update情况下才需要获取请求
            if (mode === 'update') {
              containerArr.forEach(item => {
                Object.keys(item.workloadList).forEach(async workloadType => {
                  const resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[workloadType];

                  const resourceQuery: QueryState<ResourceFilter> = {
                    filter: {
                      clusterId,
                      logAgentName,
                      isCanFetchResourceList: true,
                      namespace: item.namespaceSelection,
                      workloadType: workloadType,
                      regionId: +route.queries['rid']
                    }
                  };
                  const response = await CommonAPI.fetchResourceList({
                    query: resourceQuery,
                    resourceInfo,
                    isClearData: false
                  });
                  item.workloadList[workloadType] = response.records;
                  item.workloadListFetch[workloadType] = true;
                });
              });
            }
            dispatch({
              type: ActionType.UpdateContainerLogs,
              payload: containerArr
            });
          }
        } else if (inputMode === inputTypeMap['pod-log']) {
          //如果为容器文件路径
          inputOption = log.spec.input.pod_log_input;
          const { namespace } = log.metadata,
            containerLogFiles = inputOption.container_log_files,
            workload = inputOption.workload.name,
            workloadType = inputOption.workload.type;
          dispatch({
            type: ActionType.SelectContainerFileNamespace,
            payload: namespace
          });
          dispatch({
            type: ActionType.SelectContainerFileWorkloadType,
            payload: workloadType
          });
          //拉取数据
          const resourceInfo = resourceConfig(clusterVersion)[workloadType];
          const responseWorkloadList: RecordSet<Resource> = await CommonAPI.fetchResourceList({
            resourceInfo,
            query: {
              filter: {
                clusterId,
                logAgentName,
                namespace
              }
            }
          });
          dispatch({
            type: ActionType.UpdateContaierFileWorkloadList,
            payload: responseWorkloadList.records.map(item => {
              return {
                name: item.metadata.name,
                value: item.metadata.name
              };
            })
          });
          dispatch({
            type: ActionType.SelectContainerFileWorkload,
            payload: workload
          });
          //拉取podList数据
          dispatch(
            podActions.applyFilter({
              namespace: namespace,
              clusterId: clusterId,
              regionId: +route.queries['rid'],
              specificName: workload,
              isCanFetchPodList: namespace && workload ? true : false
            })
          );
          const containerFilePathArr = [];
          Object.keys(containerLogFiles).forEach(item => {
            containerLogFiles[item].forEach(element => {
              const containerFilePath = cloneDeep(initContainerFilePath);
              containerFilePath.containerName = item;
              containerFilePath.containerFilePath = element.path;
              containerFilePathArr.push(containerFilePath);
            });
          });
          dispatch({
            type: ActionType.UpdateContainerFilePaths,
            payload: containerFilePathArr
          });
        } else if (inputMode === inputTypeMap['host-log']) {
          //如果为主机文件路径
          inputOption = log.spec.input.host_log_input;

          let inputPath = inputOption?.path ?? '';

          const nodeInputPathType = inputPath.includes(HOST_LOG_INPUT_PATH_PREFIX) ? 'container' : 'host';

          dispatch(editLogStashActions.setNodeLogPathType(nodeInputPathType));

          if (nodeInputPathType === 'container') {
            inputPath = inputPath.replace(HOST_LOG_INPUT_PATH_PREFIX, '');
          }

          dispatch(editLogStashActions.inputNodeLogPath(inputPath));
          const metedataKeys = Object.keys(inputOption.labels);
          if (metedataKeys.length) {
            const labels = metedataKeys.map((key, index) => {
              const newLabel = Object.assign({}, initMetadata, {
                id: uuid(),
                metadataKey: key,
                metadataValue: inputOption.labels[key]
              });
              return newLabel;
            });
            dispatch({ type: ActionType.UpdateMetadata, payload: labels });
          }
        }

        if (consumerMode === outputTypeMap.kafka) {
          // kafka当中分为两种模式 ckafka 和 自建 kafka
          outputOption = log.spec.output.kafka_output;
          dispatch(editLogStashActions.inputAddressIP(outputOption.host));
          dispatch(editLogStashActions.inputAddressPort(outputOption.port));
          dispatch(editLogStashActions.inputTopic(outputOption.topic));
        } else if (consumerMode === outputTypeMap.elasticsearch) {
          outputOption = log.spec.output.elasticsearch_output;
          dispatch(editLogStashActions.inputEsAddress('http://' + outputOption.hosts[0]));
          dispatch(editLogStashActions.inputIndexName(outputOption.index));
          dispatch(editLogStashActions.inputEsUsername(outputOption.user));
          dispatch(editLogStashActions.inputEsPassword(Base64.decode(outputOption.password)));
        }
      }
    };
  }
};

//需要写一个函数获取全部的资源resource
export const logActions = extend({}, fetchLogListActions, queryLogListActions, restActions);
