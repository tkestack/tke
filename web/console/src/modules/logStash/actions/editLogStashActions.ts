import { ReduxAction, uuid } from '@tencent/ff-redux';

import { initValidator, Namespace } from '../../common/models';
import { cloneDeep, remove } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import { logModeList } from '../constants/Config';
import {
    initContainerFilePath, initContainerInputOption, initMetadata
} from '../constants/initState';
import { ContainerFilePathItem, ContainerLogs, MetadataItem, RootState } from '../models';
import { ResourceTarget, WorkLoadList } from '../models/Resource';
import { podActions } from './podActions';
import { resourceActions } from './resourceActions';
import { validatorActions } from './validatorActions';

type GetState = () => RootState;

export const editLogStashActions = {
  /** 输入日志采集器名称 */
  inputStashName: (name: string): ReduxAction<string> => {
    return {
      type: ActionType.LogStashName,
      payload: name
    };
  },

  /** 选择当前的日志采集的类型 指定容器日志 | 指定主机文件 */
  changeLogMode: (mode: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { logStashName } = getState().logStashEdit;
      dispatch({
        type: ActionType.ChangeLogMode,
        payload: mode
      });
      if (mode === logModeList.container.value) {
        dispatch(editLogStashActions.updateResourceTarget(false, true));
      } else if (mode === logModeList.containerFile.value) {
        dispatch(editLogStashActions.updateResourceTarget(true, false));
      }

      if (logStashName) {
        dispatch(validatorActions.validateStashName());
      }
    };
  },

  /**
   * pre: 当前的日志采集类型为 指定容器日志
   */
  selectAllNamespace: (mode: string): ReduxAction<string> => {
    return {
      type: ActionType.IsSelectedAllNamespace,
      payload: mode
    };
  },

  /** 选择每一个containerLog的namespace */
  selectContainerLogNamespace: (namespace: string, cIndex: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { logStashEdit, route } = getState(),
        { containerLogs } = logStashEdit;
      let containerLogsArr: ContainerLogs[] = cloneDeep(containerLogs);
      containerLogsArr[cIndex].namespaceSelection = namespace;
      dispatch({
        type: ActionType.UpdateContainerLogs,
        payload: containerLogsArr
      });
      // 进行namespace的校验
      if (getState().namespaceList.data.recordCount > 0) {
        dispatch(validatorActions.validateNamespace(cIndex));
      }

      // 清空workloadList的相关数据
      let { clusterId } = route.queries;
      dispatch(editLogStashActions.initContainerLog(true, cIndex));
      // 切換namespace之後进行resourceList列表的拉取
      const { namespaceSelection, workloadType } = getState().logStashEdit.containerLogs[cIndex];
      let isCanFetchResourceList = namespaceSelection && workloadType ? true : false;
      dispatch(
        resourceActions.applyFilter({
          clusterId,
          regionId: +route.queries['rid'],
          namespace: namespaceSelection,
          workloadType: workloadType,
          isCanFetchResourceList: isCanFetchResourceList
        })
      );
    };
  },

  /**
   * pre: 切换集群的，需要初始化所有的containerLogs的配置信息
   * @params clusterId: string  当前的集群Id
   * @params isChangeNamespace: boolean 是否为切换命名空间
   * @params cIndex:如果是切换命名空间则需要传入初始化哪个containerLogs的数组小标
   */
  initContainerLog: (isChangeNamespace: boolean = false, cIndex: number = 0) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { logStashEdit } = getState();

      let containerLog: ContainerLogs = Object.assign({}, initContainerInputOption, { id: uuid() });
      // 这里需要去判断是否为切换namespace，切换nameespace只需要清空workloadSelection、workloadList、workloadListFetch这些信息
      if (isChangeNamespace) {
        let containerLogs = cloneDeep(logStashEdit.containerLogs);
        let editingContainerLog: ContainerLogs = containerLogs.find(item => item.status === 'editing');
        containerLog.namespaceSelection = editingContainerLog.namespaceSelection;
        containerLog.collectorWay = editingContainerLog.collectorWay;
        containerLog.v_namespaceSelection = editingContainerLog.v_namespaceSelection;
        containerLog.v_workloadSelection = editingContainerLog.v_namespaceSelection;
        containerLogs[cIndex] = containerLog;
        dispatch({
          type: ActionType.UpdateContainerLogs,
          payload: containerLogs
        });
      } else {
        dispatch({
          type: ActionType.UpdateContainerLogs,
          payload: [containerLog]
        });
      }
    };
  },

  /** 更新 ContainerLog */
  updateContainerLog: (obj: any, logIndex: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerLogsArr: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
      let objKeys = Object.keys(obj);

      objKeys.forEach(keyName => {
        containerLogsArr[logIndex][keyName] = obj[keyName];
      });
      dispatch({
        type: ActionType.UpdateContainerLogs,
        payload: containerLogsArr
      });
    };
  },

  /** 删除指定容器当中的配置项 containerLog */
  deleteContainerLog: (logIndex: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerLogs: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
      containerLogs.splice(logIndex, 1);
      dispatch({
        type: ActionType.UpdateContainerLogs,
        payload: containerLogs
      });
    };
  },

  /**
   * pre: 能出发添加，说明是已经判断过是否能添加
   * 添加指定容器的配置项 containerLog
   */
  addContainerLog: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route, logStashEdit, namespaceList } = getState(),
        { containerLogs } = logStashEdit;
      let containerLogsArr: ContainerLogs[] = cloneDeep(containerLogs),
        namespaceArr: Namespace[] = cloneDeep(namespaceList.data.records);
      containerLogsArr.forEach(item => {
        remove(namespaceArr, (namespace: Namespace) => item.namespaceSelection === namespace.namespace);
      });

      // 上面剔除了已经选择namespace
      let editingIndex = containerLogsArr.findIndex(item => item.status === 'editing');
      if (editingIndex >= 0) {
        containerLogsArr[editingIndex]['status'] = 'edited';
      }

      // 添加新的containerLog到 containerLogs当中
      let newContainerLog = Object.assign({}, initContainerInputOption, {
        id: uuid(),
        namespaceSelection: namespaceArr[0].namespace
      });
      containerLogsArr.push(newContainerLog);

      dispatch({
        type: ActionType.UpdateContainerLogs,
        payload: containerLogsArr
      });

      // 添加新的containerLog之后，自动获取其对应的namespace下的resourceList
      let { clusterId, rid } = route.queries;
      dispatch(editLogStashActions.updateResourceTarget(false, true));
      dispatch(
        resourceActions.applyFilter({
          clusterId,
          regionId: +rid,
          namespace: newContainerLog.namespaceSelection,
          workloadType: newContainerLog.workloadType,
          isCanFetchResourceList: true
        })
      );
    };
  },

  /** 输入主机文件的收集路径 */
  inputNodeLogPath: (path: string): ReduxAction<string> => {
    return {
      type: ActionType.NodeLogPath,
      payload: path
    };
  },

  /** 编辑metadata */
  updateMetadata: (obj: { [props: string]: string }, mIndex: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let metadataArr: MetadataItem[] = cloneDeep(getState().logStashEdit.metadatas);
      let objKeys = Object.keys(obj);
      objKeys.forEach(item => {
        metadataArr[mIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.UpdateMetadata,
        payload: metadataArr
      });
    };
  },

  /** 删除metadata选项 */
  deleteMetadata: (mIndex: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let metadataArr: MetadataItem[] = cloneDeep(getState().logStashEdit.metadatas);
      metadataArr.splice(mIndex, 1);
      dispatch({
        type: ActionType.UpdateMetadata,
        payload: metadataArr
      });
    };
  },

  /** 新增metadata选项 */
  addMetadata: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let metadataArr: MetadataItem[] = cloneDeep(getState().logStashEdit.metadatas);
      metadataArr.push(Object.assign({}, initMetadata, { id: uuid() }));
      dispatch({
        type: ActionType.UpdateMetadata,
        payload: metadataArr
      });
    };
  },

  /** 变更消费端类型 */
  changeConsumerMode: (mode: string): ReduxAction<string> => {
    return {
      type: ActionType.ChangeConsumerMode,
      payload: mode
    };
  },

  /**
   * pre: 消费端类型为kafka，并且不使用ckafka
   */
  inputAddressIP: (addressIp: string): ReduxAction<string> => {
    return {
      type: ActionType.AddressIP,
      payload: addressIp
    };
  },

  /**
   * pre: 消费者类型为kafka，并且不使用ckafka
   */
  inputAddressPort: (addressPort: string): ReduxAction<string> => {
    return {
      type: ActionType.AddressPort,
      payload: addressPort
    };
  },

  /**
   * pre: 消费者类型为kafka，并且不使用ckafka
   */
  inputTopic: (topic: string): ReduxAction<string> => {
    return {
      type: ActionType.Topic,
      payload: topic
    };
  },

  /**
   * pre: 消费端类型为ES
   */
  inputEsAddress: (esAddress: string): ReduxAction<string> => {
    return {
      type: ActionType.EsAddress,
      payload: esAddress
    };
  },

  /**
   * pre: 消费端类型为ES
   */
  inputIndexName: (indexName: string): ReduxAction<string> => {
    return {
      type: ActionType.IndexName,
      payload: indexName
    };
  },

  /** 清空logStashEdit的内容 */
  clearLogStashEdit: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearLogStashEdit
    };
  },

  /**选择容器文件路径下的namespace */
  selectContainerFileNamespace: (namespace: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectContainerFileNamespace,
        payload: namespace
      });
      //校验规则
      if (getState().namespaceList.data.recordCount > 0) {
        dispatch(validatorActions.validateContainerFileNamespace());
      }

      // 进行resourceList列表的拉取
      let { logStashEdit, route, clusterSelection } = getState();
      let { name } = clusterSelection[0].metadata;
      let { containerFileNamespace, containerFileWorkloadType } = logStashEdit;

      let isCanFetchResourceList = containerFileNamespace && containerFileWorkloadType ? true : false;

      if (isCanFetchResourceList) {
        dispatch(
          resourceActions.applyFilter({
            clusterId: name,
            regionId: +route.queries['rid'],
            namespace: containerFileNamespace,
            workloadType: containerFileWorkloadType
          })
        );
      }
    };
  },
  /**选择容器文件路径下的workloadType */
  selectContainerFileWorkloadType: (workloadType: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectContainerFileWorkloadType,
        payload: workloadType
      });

      //校验规则
      dispatch(validatorActions.validateContainerFileWorkloadType());

      // 进行resourceList列表的拉取

      let { logStashEdit, route, clusterSelection } = getState();
      let { containerFileNamespace, containerFileWorkloadType } = logStashEdit;
      if (clusterSelection[0]) {
        let { name } = clusterSelection[0].metadata;
        let isCanFetchResourceList = containerFileNamespace && containerFileWorkloadType ? true : false;
        dispatch(
          resourceActions.applyFilter({
            clusterId: name,
            regionId: +route.queries['rid'],
            namespace: containerFileNamespace,
            workloadType: containerFileWorkloadType,
            isCanFetchResourceList: isCanFetchResourceList
          })
        );
      }
    };
  },

  /**选择容器文件路径下的workload */
  selectContainerFileWorkload: (workload: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectContainerFileWorkload,
        payload: workload
      });

      if (getState().logStashEdit.containerFileWorkloadList.length > 0) {
        //校验规则
        dispatch(validatorActions.validateContainerFileWorkload());
      }

      //拉取podList
      let { logStashEdit, route, clusterSelection } = getState();
      if (clusterSelection[0]) {
        let { name } = clusterSelection[0].metadata;
        const { containerFileNamespace, containerFileWorkload } = logStashEdit;
        if (containerFileNamespace && containerFileWorkload) {
          dispatch(
            podActions.applyFilter({
              namespace: containerFileNamespace,
              clusterId: name,
              regionId: +route.queries['rid'],
              specificName: containerFileWorkload,
              isCanFetchPodList: containerFileNamespace && containerFileWorkload ? true : false
            })
          );
        } else {
          dispatch(
            podActions.fetch({
              noCache: true
            })
          );
        }
        dispatch(editLogStashActions.initContainerFileContainerFilePaths());
      }
    };
  },

  /**选择容器文件路径下的容器名*/
  selectContainerFileContainerName: (containerName: string, index: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { logStashEdit } = getState();
      let containerFilePathArr: ContainerFilePathItem[] = cloneDeep(logStashEdit.containerFilePaths);
      containerFilePathArr[index].containerName = containerName;
      dispatch({
        type: ActionType.UpdateContainerFilePaths,
        payload: containerFilePathArr
      });
      //检验规则
      dispatch(validatorActions.validateContainerFileContainerName(index));
      dispatch(validatorActions.validateContainerFileContainerFilePath(index));
    };
  },

  /**输入容器文件路径下的文件路径*/
  inputContainerFileContainerFilePath: (containerFilePath: string, index: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerFilePathArr: ContainerFilePathItem[] = cloneDeep(getState().logStashEdit.containerFilePaths);
      containerFilePathArr[index].containerFilePath = containerFilePath;
      dispatch({
        type: ActionType.UpdateContainerFilePaths,
        payload: containerFilePathArr
      });
    };
  },

  /**添加容器文件路径下的文件路径 */
  addContainerFileContainerFilePath: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerFilePathArr: ContainerFilePathItem[] = cloneDeep(getState().logStashEdit.containerFilePaths);
      containerFilePathArr.push(Object.assign({}, initContainerFilePath, { id: uuid() }));
      dispatch({
        type: ActionType.UpdateContainerFilePaths,
        payload: containerFilePathArr
      });
    };
  },

  /**删除容器文件路径下的文件路径 */
  deleteContainerFileContainerFilePath: (index: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerFilePathArr: ContainerFilePathItem[] = cloneDeep(getState().logStashEdit.containerFilePaths);
      containerFilePathArr.splice(index, 1);
      dispatch({
        type: ActionType.UpdateContainerFilePaths,
        payload: containerFilePathArr
      });
    };
  },
  /**切换集群，切换workload需要初始化 containerFilePaths*/
  initContainerFileContainerFilePaths: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerFilePathArr = [cloneDeep(initContainerFilePath)];
      dispatch({
        type: ActionType.UpdateContainerFilePaths,
        payload: containerFilePathArr
      });
    };
  },

  /**更換resourseList的資源使用對象 */
  updateResourceTarget: (isForContainerFile: boolean, isForContainerLogs: boolean): ReduxAction<ResourceTarget> => {
    return {
      type: ActionType.UpdateResourceTarget,
      payload: { isForContainerLogs: isForContainerLogs, isForContainerFile: isForContainerFile }
    };
  },
  /**更換containerFileWorkloadList的列表 */
  updateContainerFileWorkloadList: (workloadList: WorkLoadList[]): ReduxAction<WorkLoadList[]> => {
    return {
      type: ActionType.UpdateContaierFileWorkloadList,
      payload: workloadList
    };
  },
  /**是否是第一次拉取 */
  ifFirstFetchResource: (ifFirstFetchResource: boolean): ReduxAction<boolean> => {
    return {
      type: ActionType.isFirstFetchResource,
      payload: ifFirstFetchResource
    };
  },

  reInitEdit: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(editLogStashActions.initContainerFileContainerFilePaths());
      dispatch(editLogStashActions.initContainerLog());
      dispatch(editLogStashActions.updateContainerFileWorkloadList([]));
      dispatch(editLogStashActions.ifFirstFetchResource(true));
      dispatch({
        type: ActionType.SelectContainerFileWorkloadType,
        payload: 'deployment'
      });

      dispatch({
        type: ActionType.V_ContainerFileNamespace,
        payload: initValidator
      });
      dispatch({
        type: ActionType.V_ContainerFileWorkload,
        payload: initValidator
      });
      dispatch({
        type: ActionType.V_ContainerFileWorkloadType,
        payload: initValidator
      });
    };
  }
};
