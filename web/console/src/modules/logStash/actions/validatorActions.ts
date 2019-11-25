import { RootState, ContainerLogs, MetadataItem, WorkloadType, LogStashEdit, ContainerFilePathItem } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../WebAPI';
import { router } from '../router';
import { cloneDeep } from '../../common/utils';
import { initValidator } from 'src/modules/common';
import { isTemplateExpression } from 'typescript';
import { namespace } from 'config/resource/k8sConfig';

type GetState = () => RootState;

export const validatorActions = {
  /** 校验采集器的名称是否正确 */
  async _validateStashName(
    cluserVersion: string,
    namespace: string,
    name: string,
    clusterId: string,
    mode: string,
    regionId: number
  ) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    //验证集群名称
    if (!name) {
      status = 2;
      message = '日志收集规则名称不能为空';
    } else if (name.length > 63) {
      status = 2;
      message = '日志收集规则名称不能超过63个字符';
    } else if (!reg.test(name)) {
      status = 2;
      message = '日志收集规则名称格式不正确';
    } else {
      let res = false;
      if (mode === 'create') {
        res = await WebAPI.checkStashNameIsExist(cluserVersion, name, clusterId, regionId, namespace);
      }
      if (res) {
        status = 2;
        message = '日志收集器名称已存在';
      } else {
        status = 1;
        message = '';
      }
    }
    return { status, message };
  },

  validateStashName() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route, logStashEdit, clusterVersion } = getState(),
        { logStashName, logMode, containerFileNamespace } = logStashEdit,
        urlParams = router.resolve(route),
        { clusterId, rid } = route.queries;
      let namespace = logMode === 'containerFile' ? containerFileNamespace : 'kube-system';
      const result = await validatorActions._validateStashName(
        clusterVersion,
        namespace,
        logStashName,
        clusterId,
        urlParams['mode'],
        +rid
      );
      dispatch({
        type: ActionType.V_LogStashName,
        payload: result
      });
    };
  },

  /** 校验当前的containerLog的命名空间 */
  _validateNamespace(namespace: string, containerLogs: ContainerLogs[]) {
    let status = 0,
      message = '';

    if (!namespace) {
      status = 2;
      message = 'Namespace不能为空';
    } else if (containerLogs.filter(item => item.namespaceSelection === namespace).length > 1) {
      status = 2;
      message = 'Namespace不能重复';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateNamespace(logIndex: number) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { containerLogs } = getState().logStashEdit,
        containerLogsArr: ContainerLogs[] = cloneDeep(containerLogs);
      let containerLog: ContainerLogs = containerLogs[logIndex];
      let result = validatorActions._validateNamespace(containerLog.namespaceSelection, containerLogs);
      containerLogsArr[logIndex].v_namespaceSelection = result;
      dispatch({
        type: ActionType.UpdateContainerLogs,
        payload: containerLogsArr
      });
    };
  },

  /**
   * pre: 日志类型为容器日志
   * 验证所属服务的填写是否正确
   */
  _canAddContainerLog(containerLog: ContainerLogs, containerLogs: ContainerLogs[]) {
    let canAdd = true;
    canAdd = canAdd && validatorActions._validateNamespace(containerLog.namespaceSelection, containerLogs).status === 1;
    if (containerLog.collectorWay === 'workload') {
      canAdd = canAdd && validatorActions._validateWorkloadSelectedNumber(containerLog.workloadSelection).status === 1;
    }
    return canAdd;
  },

  /**
   * pre: 日志类型为容器日志, 并且为指定容器日志，采集对象为按工作负载选择
   * 校验工作负载选择的数量
   */
  _validateWorkloadSelectedNumber(workloadSelection: WorkloadType<string>) {
    let status = 0,
      message = '';
    let { deployment, statefulset, daemonset, job, cronjob } = workloadSelection;
    if (
      deployment.length === 0 &&
      statefulset.length === 0 &&
      daemonset.length === 0 &&
      job.length === 0 &&
      cronjob.length === 0
    ) {
      status = 2;
      message = '已选工作负载项为0个，请至少选择一个工作负载项或者选择全部容器';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateWorkloadSelectedNumber(containerLogIndex: number) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerLogArr: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
      let containerLog = containerLogArr[containerLogIndex];
      let result = validatorActions._validateWorkloadSelectedNumber(containerLog.workloadSelection);
      containerLogArr[containerLogIndex]['v_workloadSelection'] = result;
      dispatch({
        type: ActionType.UpdateContainerLogs,
        payload: containerLogArr
      });
    };
  },

  /**
   * pre: 日志类型为 容器日志
   * 验证当前容器日志
   */
  validateContainerLog(containerLogIndex: number) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { collectorWay } = getState().logStashEdit.containerLogs[containerLogIndex];
      if (collectorWay === 'container') {
        //如果是指定容器 则不需要校验 工作负载选项
        dispatch(validatorActions.validateNamespace(containerLogIndex));
      } else if (collectorWay === 'workload') {
        // 校验命名空间的选择
        dispatch(validatorActions.validateNamespace(containerLogIndex));
        // 校验工作负载的选择
        dispatch(validatorActions.validateWorkloadSelectedNumber(containerLogIndex));
      }
    };
  },

  /**
   * pre: 日志类型为 指定主机文件
   * 校验收集路径
   */
  async _validateNodeLogPath(path: string) {
    let firstStrReg = /^\/$/g,
      status = 0,
      message = '';

    if (path.length === 0) {
      status = 2;
      message = '日志收集路径不能为空';
    } else if (!firstStrReg.test(path.split('')[0])) {
      status = 2;
      message = '日志收集路径不正确';
    } else {
      status = 1;
      message = '';
      // let res = await WebAPI.checkNodeLogPathIsValid(path);

      // if (!res) {
      //   status = 2;
      //   message = '日志收集路径不正确..';
      // }
    }
    return { status, message };
  },

  validateNodeLogPath() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { logStashEdit } = getState(),
        { nodeLogPath } = logStashEdit;

      let result = await validatorActions._validateNodeLogPath(nodeLogPath);
      dispatch({
        type: ActionType.V_NodeLogPath,
        payload: result
      });
    };
  },

  /**
   * pre: 日志类型为主机日志
   * 验证metadata的填写是否正确
   */
  _validateMetadataItem(value: string, metadataArr: MetadataItem[], isKey: boolean) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    if (!value) {
      status = 2;
      message = isKey ? 'Key值不能为空' : 'Value值不能为空';
    } else if (value.length > 63) {
      status = 2;
      message = '不能超过63个字符';
    } else if (isKey && !reg.test(value)) {
      if (!reg.test(value)) {
        status = 2;
        message = '格式不正确';
      } else if (metadataArr.filter(item => item.metadataKey === value).length > 1) {
        status = 2;
        message = 'Key值不能重复';
      }
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateMetadataItem(obj: { [props: string]: string }, mIndex: number) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let metadataArr: MetadataItem[] = cloneDeep(getState().logStashEdit.metadatas),
        keyName = Object.keys(obj)[0];

      let result = validatorActions._validateMetadataItem(obj[keyName], metadataArr, keyName === 'metadataKey');
      metadataArr[mIndex]['v_' + keyName] = result;
      dispatch({
        type: ActionType.UpdateMetadata,
        payload: metadataArr
      });
    };
  },

  /** 校验所有的Metadata的配置 */
  _validateAllMetadataItem(metadatas: MetadataItem[]) {
    let result = true;
    metadatas.forEach((metadata, index) => {
      let resultKey = validatorActions._validateMetadataItem(metadata.metadataKey, metadatas, true),
        resultValue = validatorActions._validateMetadataItem(metadata.metadataValue, metadatas, false);

      result = result && resultKey.status === 1 && resultValue.status === 1;
    });
    return result;
  },

  validateAllMetadataItem() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let metadatas: MetadataItem[] = cloneDeep(getState().logStashEdit.metadatas);
      metadatas.forEach((metadata, index) => {
        let resultKey = validatorActions._validateMetadataItem(metadata.metadataKey, metadatas, true),
          resultValue = validatorActions._validateMetadataItem(metadata.metadataValue, metadatas, false);
        metadatas[index]['v_metadataKey'] = resultKey;
        metadatas[index]['v_metadataValue'] = resultValue;
      });
      dispatch({
        type: ActionType.UpdateMetadata,
        payload: metadatas
      });
    };
  },

  /**
   * pre: 消费端类型为kafka类型
   * 校验访问地址IP
   */
  _validateAddressIP(addressIP: string) {
    let reg = /^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$/,
      status = 0,
      message = '';

    if (!addressIP) {
      status = 2;
      message = 'IP地址不能为空';
    } else if (!reg.test(addressIP)) {
      status = 2;
      message = 'IP地址格式不正确';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateAddressIP() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let result = validatorActions._validateAddressIP(getState().logStashEdit.addressIP);
      dispatch({
        type: ActionType.V_AddressIP,
        payload: result
      });
    };
  },

  /**
   * pre: 消费端类型为kafka类型
   * 校验访问地址的端口
   */
  _validateAddressPort(addressPort: string) {
    let status = 0,
      message = '';

    if (!addressPort) {
      status = 2;
      message = 'IP端口不能为空';
    } else if (+addressPort < 0 || +addressPort > 65535) {
      status = 2;
      message = '您的端口设置超过了范围：0-65535';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateAddressPort() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let result = validatorActions._validateAddressPort(getState().logStashEdit.addressPort);
      dispatch({
        type: ActionType.V_AddressPort,
        payload: result
      });
    };
  },

  /**
   * pre: 消费端类型为kafka类型
   * 校验日志主题topic
   */
  _validateTopic(topic: string) {
    let reg = /^[a-zA-Z0-9\\._\\-]+$/,
      status = 0,
      message = '';

    if (!topic) {
      status = 2;
      message = 'topic名称不能为空';
    } else if (topic.length > 64) {
      status = 2;
      message = 'topic名称不能超过64个字符';
    } else if (!reg.test(topic)) {
      status = 2;
      message = 'topic名称格式不正确';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateTopic() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let result = validatorActions._validateTopic(getState().logStashEdit.topic);
      dispatch({
        type: ActionType.V_Topic,
        payload: result
      });
    };
  },

  /** 校验当前的es地址是否正确 */
  _validateEsAddress(address: string) {
    let status = 0,
      message = '',
      hostReg = /^((http|https):\/\/)((25[0-5]|2[0-4]\d|1?\d?\d)\.){3}(25[0-5]|2[0-4]\d|1?\d?\d):([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/;

    if (!address) {
      status = 2;
      message = 'Elasticsearch地址不能为空';
    } else if (!hostReg.test(address)) {
      status = 2;
      message = 'Elasticsearch地址格式不正确，{scheme}://{addr}:{port}';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateEsAddress() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { esAddress } = getState().logStashEdit;

      let result = validatorActions._validateEsAddress(esAddress);
      dispatch({
        type: ActionType.V_EsAddress,
        payload: result
      });
    };
  },

  /** 校验当前的索引名是否正确 */
  _validateIndexName(indexName: string) {
    let status = 0,
      message = '',
      reg = /^[a-z][0-9a-z_+-]+$/;
    if (!indexName) {
      status = 2;
      message = '索引名不能为空';
    } else if (!reg.test(indexName)) {
      status = 2;
      message = '索引名格式不正确';
    } else if (indexName.length > 60) {
      status = 2;
      message = '索引名不能超过60个字符';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateIndexName() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { indexName } = getState().logStashEdit;
      let result = validatorActions._validateIndexName(indexName);
      dispatch({
        type: ActionType.V_IndexName,
        payload: result
      });
    };
  },

  /** 校验日志采集的编辑项是否正确 */
  async _validateLogStashEdit(
    logStashEdit: LogStashEdit,
    namespace: string,
    clusterVersion: string,
    clusterId: string,
    mode: string,
    regionId: number
  ) {
    let {
      logStashName,
      logMode,
      isSelectedAllNamespace,
      containerLogs,
      nodeLogPath,
      metadatas,
      consumerMode,
      esAddress,
      indexName,
      addressIP,
      addressPort,
      topic,
      containerFileNamespace,
      containerFileWorkload,
      containerFileWorkloadType,
      containerFilePaths
    } = logStashEdit;

    let result = true;
    // 校验日志采集的名称
    result =
      result &&
      (await validatorActions._validateStashName(clusterVersion, namespace, logStashName, clusterId, mode, regionId))
        .status === 1;

    if (logMode === 'container') {
      // 判断是否为 指定容器，如果为所有容器，则不需要校验
      if (result && isSelectedAllNamespace === 'selectOne') {
        containerLogs.forEach(item => {
          result = result && validatorActions._canAddContainerLog(item, containerLogs);
        });
      }
    } else if (logMode === 'node') {
      result = result && (await validatorActions._validateNodeLogPath(nodeLogPath)).status === 1;
      if (result && metadatas.length) {
        result = result && validatorActions._validateAllMetadataItem(metadatas);
      }
    } else if (logMode === 'containerFile') {
      result = result && validatorActions._validateContainerFileNamespace(containerFileNamespace).status === 1;
      result = result && validatorActions._validateContainerFileWorkloadType(containerFileWorkloadType).status === 1;
      result = result && validatorActions._validateContainerFileWorkload(containerFileWorkload).status === 1;
      result = result && validatorActions._validateAllContainerFilePath(containerFilePaths);
    }

    // 校验消费端的相关合法性
    if (result) {
      if (consumerMode === 'kafka') {
        result =
          result &&
          validatorActions._validateAddressIP(addressIP).status === 1 &&
          validatorActions._validateAddressPort(addressPort).status === 1 &&
          validatorActions._validateTopic(topic).status === 1;
      } else if (consumerMode === 'es') {
        result =
          result &&
          validatorActions._validateEsAddress(esAddress).status === 1 &&
          validatorActions._validateIndexName(indexName).status === 1;
      }
    }

    return result;
  },

  validateLogStashEdit() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { logStashEdit, clusterVersion } = getState(),
        { logMode, isSelectedAllNamespace, containerLogs, metadatas, consumerMode } = logStashEdit;
      // 校验日志采集的名称
      dispatch(validatorActions.validateStashName());

      // 当日志类型为 指定容器日志，验证内容如下
      if (logMode === 'container') {
        // 判断是否为 指定容器，如果为所有容器，则不需要校验
        if (isSelectedAllNamespace === 'selectOne') {
          containerLogs.forEach((item, index) => {
            dispatch(validatorActions.validateContainerLog(index));
          });
        }
      } else if (logMode === 'node') {
        dispatch(validatorActions.validateNodeLogPath());
        if (metadatas.length) {
          dispatch(validatorActions.validateAllMetadataItem());
        }
      } else if (logMode === 'containerFile') {
        // 当前日志类型为 容器文件路径
        dispatch(validatorActions.validateAllContainerFilePath());
        dispatch(validatorActions.validateContainerFileNamespace());
        dispatch(validatorActions.validateContainerFileWorkloadType());
        dispatch(validatorActions.validateContainerFileWorkload());
      }

      // 校验消费端的相关合法性
      if (consumerMode === 'kafka') {
        dispatch(validatorActions.validateAddressIP());
        dispatch(validatorActions.validateAddressPort());
        dispatch(validatorActions.validateTopic());
      } else if (consumerMode === 'es') {
        dispatch(validatorActions.validateEsAddress());
        dispatch(validatorActions.validateIndexName());
      }
    };
  },

  /**
   * pre:日志类型为 指定容器日志文件
   * 校验namespace
   */
  validateContainerFileNamespace() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const namespace = getState().logStashEdit.containerFileNamespace;
      let result = validatorActions._validateContainerFileNamespace(namespace);
      dispatch({
        type: ActionType.V_ContainerFileNamespace,
        payload: result
      });
    };
  },

  _validateContainerFileNamespace(namespace: string) {
    let status = 0,
      message = '';
    if (!namespace) {
      status = 2;
      message = '命名空间不能为空';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  /**
   * pre:日志类型为 指定容器日志文件
   * 校验工作负载类型workloadType
   */
  validateContainerFileWorkloadType() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const WorkloadType = getState().logStashEdit.containerFileWorkloadType;
      let result = validatorActions._validateContainerFileNamespace(WorkloadType);
      dispatch({
        type: ActionType.V_ContainerFileWorkloadType,
        payload: result
      });
    };
  },

  _validateContainerFileWorkloadType(workloadType: string) {
    let status = 0,
      message = '';
    if (!workloadType) {
      status = 2;
      message = '工作负载类型不能为空';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  /**
   * pre:日志类型为 指定容器日志文件
   * 校验工作负载workload
   */
  validateContainerFileWorkload() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const Workload = getState().logStashEdit.containerFileWorkload;
      let result = validatorActions._validateContainerFileWorkload(Workload);
      dispatch({
        type: ActionType.V_ContainerFileWorkload,
        payload: result
      });
    };
  },

  _validateContainerFileWorkload(workload: string) {
    let status = 0,
      message = '';
    if (!workload) {
      status = 2;
      message = '工作负载不能为空';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  /**
   * pre：日志类型为 指定容器文件路径
   * 校验容器名
   */
  validateContainerFileContainerName(index: number) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerFilePathArr: ContainerFilePathItem[] = cloneDeep(getState().logStashEdit.containerFilePaths);
      let { containerName } = containerFilePathArr[index];
      let result = validatorActions._validateContainerFileContainerName(containerName);

      containerFilePathArr[index].v_containerName = result;
      dispatch({
        type: ActionType.UpdateContainerFilePaths,
        payload: containerFilePathArr
      });
    };
  },
  _validateContainerFileContainerName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    //验证容器名称
    if (!name) {
      status = 2;
      message = '容器名不能为空';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  /**
   * pre：日志类型为 指定容器文件路径
   * 校验容器文件路径
   */

  validateContainerFileContainerFilePath(index: number) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerFilePathArr: ContainerFilePathItem[] = cloneDeep(getState().logStashEdit.containerFilePaths);
      let result = validatorActions._validateContainerFileContainerFilePath(
        containerFilePathArr[index].containerName,
        containerFilePathArr[index].containerFilePath,
        containerFilePathArr,
        index
      );
      containerFilePathArr[index].v_containerFilePath = result;
      dispatch({
        type: ActionType.UpdateContainerFilePaths,
        payload: containerFilePathArr
      });
    };
  },

  _validateContainerFileContainerFilePath(name: string, path: string, paths: ContainerFilePathItem[], index: number) {
    let firstStrReg = /^\/$/g,
      status = 0,
      message = '';

    if (path.length === 0) {
      status = 2;
      message = '文件路径不能为空';
    } else if (path.length > 63) {
      status = 2;
      message = '文件路径最多支持63个字符';
    } else if (!firstStrReg.test(path.split('')[0])) {
      if (path !== 'stdout') {
        status = 2;
        message = '文件路径不正确';
      } else {
        status = 1;
        message = '';
      }
    } else if (
      paths.findIndex(
        (item, indexPath) => indexPath !== index && item.containerFilePath === path && item.containerName === name
      ) !== -1
    ) {
      status = 2;
      message = '文件路径不能重复';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },
  /**校验所有的容器文件路径 + 容器名 */
  validateAllContainerFilePath() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let containerFilePathArr: ContainerFilePathItem[] = cloneDeep(getState().logStashEdit.containerFilePaths);
      containerFilePathArr.forEach((item, index) => {
        let nameResult = validatorActions._validateContainerFileContainerName(item.containerName);
        let filePathResult = validatorActions._validateContainerFileContainerFilePath(
          item.containerName,
          item.containerFilePath,
          containerFilePathArr,
          index
        );
        containerFilePathArr[index].v_containerName = nameResult;
        containerFilePathArr[index].v_containerFilePath = filePathResult;
      });
      dispatch({
        type: ActionType.UpdateContainerFilePaths,
        payload: containerFilePathArr
      });
    };
  },

  _validateAllContainerFilePath(contaierFilePaths: ContainerFilePathItem[]) {
    let result = true;
    contaierFilePaths.forEach((item, index) => {
      let nameResult = validatorActions._validateContainerFileContainerName(item.containerName);
      let filePathResult = validatorActions._validateContainerFileContainerFilePath(
        item.containerName,
        item.containerFilePath,
        contaierFilePaths,
        index
      );
      result = result && nameResult.status === 1 && filePathResult.status === 1;
    });

    return result;
  }
};
