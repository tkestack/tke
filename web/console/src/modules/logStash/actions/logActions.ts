import { extend, ReduxAction, uuid, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator, QueryState } from '@tencent/qcloud-redux-query';
import * as ActionType from '../constants/ActionType';
import { RootState, LogFilter, Log, ResourceFilter, ContainerLogs } from '../models';
import * as WebAPI from '../WebAPI';
import { Validation, ResourceInfo, initValidator } from '../../common/models';
import { inputTypeMap, outputTypeMap } from '../constants/Config';
import { editLogStashActions } from './editLogStashActions';
import { initContainerInputOption, initContainerFilePath, initMetadata } from '../constants/initState';
import { cloneDeep } from '../../common/utils';
import { resourceConfig } from '../../../../config';
import { Resource } from '../models/Resource';
import { podActions } from './podActions';
import { CommonAPI } from '../../common';
import { resourceActions } from './resourceActions';

type GetState = () => RootState;

/** 获取Log采集器的列表的Action */
const fetchLogListActions = generateFetcherActionCreator({
  actionType: ActionType.FetchLogList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { logQuery, clusterVersion } = getState();
    let resourceInfo = resourceConfig(clusterVersion)['logcs'];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await CommonAPI.fetchResourceList({ resourceInfo, query: logQuery, isClearData });
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
      let { clusterVersion, route } = getState();
      let resourceInfo = resourceConfig(clusterVersion)['logcs'];
      let { clusterId, regionId } = route.queries;
      let result = await CommonAPI.fetchResourceList({
        query: {
          filter: {
            namespace,
            clusterId,
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

      let log = getState().logSelection[0];
      if (mode === 'update' || mode === 'detail') {
        let inputOption; //输入端选项
        let outputOption; //输出端选项
        let inputMode = log.spec.input.type ? inputTypeMap[log.spec.input.type] : ''; //输入端类型
        let consumerMode = log.spec.output.type ? outputTypeMap[log.spec.output.type] : ''; //输出端类型
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
            //选择了所有容器，则需要帮忙选择指定容器项目,
            //只有在update的时候才需要去获取resource、
            if (mode === 'update') {
              let containerLogsArr: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
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
            let containerArr = [];
            inputOption.namespaces.forEach(item => {
              let tmp = cloneDeep(initContainerInputOption);
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
                  let resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[workloadType];

                  let resourceQuery: QueryState<ResourceFilter> = {
                    filter: {
                      clusterId,
                      isCanFetchResourceList: true,
                      namespace: item.namespaceSelection,
                      workloadType: workloadType,
                      regionId: +route.queries['rid']
                    }
                  };
                  let response = await CommonAPI.fetchResourceList({
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
          let { namespace } = log.metadata,
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
          let resourceInfo = resourceConfig(clusterVersion)[workloadType];
          let responseWorkloadList: RecordSet<Resource> = await CommonAPI.fetchResourceList({
            resourceInfo,
            query: {
              filter: {
                clusterId,
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
          let containerFilePathArr = [];
          Object.keys(containerLogFiles).forEach(item => {
            containerLogFiles[item].forEach(element => {
              let containerFilePath = cloneDeep(initContainerFilePath);
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
          dispatch(editLogStashActions.inputNodeLogPath(inputOption.path));
          let metedataKeys = Object.keys(inputOption.labels);
          if (metedataKeys.length) {
            let labels = metedataKeys.map((key, index) => {
              let newLabel = Object.assign({}, initMetadata, {
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
        }
      }
    };
  }
};

//需要写一个函数获取全部的资源resource
export const logActions = extend({}, fetchLogListActions, queryLogListActions, restActions);
