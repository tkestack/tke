import * as React from 'react';
import { connect } from 'react-redux';
import { CreateResource } from 'src/modules/cluster/models';

import { Button, ExternalLink, Segment, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import {
    bindActionCreators, FetchState, isSuccessWorkflow, OperationState, uuid
} from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { SegmentOption } from '@tencent/tea-component/lib/segment/SegmentOption';

import { resourceConfig } from '../../../../config/resourceConfig';
import { SelectList, TipInfo } from '../../common/components';
import { cloneDeep, getWorkflowError } from '../../common/utils';
import { allActions } from '../actions';
import { validatorActions } from '../actions/validatorActions';
import { logModeList } from '../constants/Config';
import {
    ContainerLogInput, ContainerLogNamespace, ElasticsearchOutput, HostLogInput, KafkaOutpot,
    LogStashEditYaml, PodLogInput
} from '../models/LogStashEdit';
import { router } from '../router';
import { EditConsumerPanel } from './EditConsumerPanel';
import { EditOriginContainerFilePanel } from './EditOriginContainerFilePanel';
import { EditOriginContainerPanel } from './EditOriginContainerPanel';
import { EditOriginNodePanel } from './EditOriginNodePanel';
import { isCanCreateLogStash } from './LogStashActionPanel';
import { RootProps } from './LogStashApp';

/** 日志采集类型的提示 */
const logModeTip = {
  container: {
    text: t('采集集群内任意服务下的容器日志，仅支持Stderr和Stdout的日志。')
  },
  node: {
    text: t('采集集群内指定节点路径的文件。')
  },
  containerFile: {
    text: t('采集集群内指定容器路径的文件。')
  }
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class EditLogStashPanel extends React.Component<RootProps, any> {
  componentWillUnmount() {
    let { actions, route } = this.props;
    // 清理logStashEdit的内容
    actions.editLogStash.clearLogStashEdit();
    // 重置workflow
    actions.workflow.modifyLogStash.reset();
  }

  render() {
    let {
        actions,
        logStashEdit,
        regionSelection,
        route,
        clusterSelection,
        logList,
        clusterList,
        isOpenLogStash,
        modifyLogStashFlow,
        isDaemonsetNormal,
        logDaemonset,
        projectList,
        projectSelection
      } = this.props,
      { logStashName, v_logStashName, v_clusterSelection, logMode } = logStashEdit,
      urlParams = router.resolve(route);
    // 当前的类型 create | update
    let { mode } = urlParams;
    let byProject = window.location.href.includes('/tkestack-project');

    let isCreateMode: boolean = mode === 'create';

    // 判断当前是否能够新建日志收集规则
    let { canCreate, tip, ifLogDaemonset } = isCanCreateLogStash(
      clusterSelection[0],
      logList.data.records,
      isDaemonsetNormal,
      isOpenLogStash
    );

    /** 渲染日志类型 */
    let selectedLogMode = Object.values(logModeList).find(item => item.value === logMode);

    let ifLogDaemonsetNeedLoading = logDaemonset.fetchState === FetchState.Fetching;

    /** 创建日志采集规则失败 */
    let failed = modifyLogStashFlow.operationState === OperationState.Done && !isSuccessWorkflow(modifyLogStashFlow);

    //渲染集群列表selectList选择项
    const selectClusterList = cloneDeep(clusterList);
    selectClusterList.data.records = clusterList.data.records.map(cluster => {
      return { clusterId: cluster.metadata.name, clusterName: cluster.spec.displayName };
    });

    const logModeSegments: SegmentOption[] = Object.keys(logModeList).map(mode => {
      return {
        value: logModeList[mode].value,
        text: logModeList[mode].name
      };
    });
    // 如果是在业务侧，去掉"节点文件路径"
    if (byProject) {
      logModeSegments.pop();
    }

    // 根据当前是平台侧/业务侧返回业务/集群选择器，或者显示业务/集群信息（修改）
    let getSelector = () => {
      if (byProject) {
        if (!isCreateMode) {
          return (
            <FormPanel.Item label={t('所属业务')} text>
              {projectSelection}
            </FormPanel.Item>
          );
        }
        let projectListOptions = projectList.map((p, index) => ({
          text: p.displayName,
          value: p.name
        }));

        return (<>
          <FormPanel.Item label={t('业务：')}>
            <FormPanel.Select
              options={projectListOptions}
              value={projectSelection}
              onChange={value => {
                actions.cluster.selectProject(value);
              }}
            ></FormPanel.Select>
          </FormPanel.Item>
        </>);
      }
      if (!isCreateMode) {
        return (
          <FormPanel.Item label={t('所属集群')} text>
            {clusterSelection[0] &&
            clusterSelection[0].metadata.name + '(' + clusterSelection[0].spec.displayName + ')'}
          </FormPanel.Item>
        );
      }

      return (
        <FormPanel.Item
          label={t('所属集群')}
          message={
            <React.Fragment>
              <Text parent="p">
                <Trans>
                  如现有的集群不合适，您可以去控制台
                  <ExternalLink href={`/tke/cluster/create?rid=${route.queries['rid']}`} target="_self">
                    导入集群
                  </ExternalLink>
                  或者
                  <ExternalLink href={`/tke/cluster/createIC?rid=${route.queries['rid']}`} target="_self">
                    新建一个独立集群
                  </ExternalLink>
                </Trans>
              </Text>
              {!(clusterSelection && clusterSelection[0] && clusterSelection[0].spec.logAgentName || isOpenLogStash) && (
                <Text theme="danger">
                  <Trans>
                    该集群未开启日志收集功能，
                    <Button type="link" onClick={() => actions.workflow.authorizeOpenLog.start()}>
                      立即开启
                    </Button>
                  </Trans>
                </Text>
              )}
            </React.Fragment>
          }
        >
          <SelectList
            value={clusterSelection[0] ? clusterSelection[0].metadata.name : ''}
            recordData={selectClusterList}
            valueField="clusterId"
            textField="clusterName"
            textFields={['clusterId', 'clusterName']}
            textFormat={`\${clusterId} (\${clusterName})`}
            className="tc-15-select m"
            style={{ marginRight: '5px' }}
            onSelect={value => actions.cluster.selectCluster(value)}
            name={t('集群')}
            emptyTip=""
            tipPosition="left"
            align="start"
            validator={v_clusterSelection}
            isUnshiftDefaultItem={false}
          />
        </FormPanel.Item>
      );
    };

    return (
      <FormPanel>
        {!isCreateMode ? (
          <FormPanel.Item label={t('收集规则名称')} text>
            {route.queries['stashName']}
          </FormPanel.Item>
        ) : (
          <FormPanel.Item
            label={t('收集规则名称')}
            text={mode === 'update'}
            validator={v_logStashName}
            message={t('最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')}
            input={{
              placeholder: t('请输入日志收集规则名称'),
              value: logStashName,
              onChange: value => {
                actions.editLogStash.inputStashName(value);
              },
              onBlur: actions.validate.validateStashName
            }}
          />
        )}

        {getSelector()}

        {!isCreateMode ? (
          <FormPanel.Item label={t('类型')} text>
            {selectedLogMode.name}
          </FormPanel.Item>
        ) : (
          <FormPanel.Item
            label={t('类型')}
            message={
              <React.Fragment>
                <Text>
                  <Trans>
                    {logModeTip[logMode].text}
                    <ExternalLink href={logModeTip[logMode].href}>{t('查看实例')}</ExternalLink>
                  </Trans>
                </Text>
              </React.Fragment>
            }
          >
            <Segment
              options={logModeSegments}
              value={selectedLogMode.value}
              onChange={value => actions.editLogStash.changeLogMode(value)}
            />
          </FormPanel.Item>
        )}

        <EditOriginContainerPanel isEdit={mode === 'update'} />

        <EditOriginNodePanel />

        <EditOriginContainerFilePanel isEdit={mode === 'update'} />

        <EditConsumerPanel />
        <FormPanel.Footer>
          <Button
            type="primary"
            disabled={modifyLogStashFlow.operationState === OperationState.Performing || !canCreate}
            onClick={() => {
              this._handleSubmit(mode);
            }}
            style={{
              marginRight: '20px'
            }}
          >
            {t('完成')}
          </Button>
          <Button
            type="weak"
            onClick={e => {
              let newRouteQueies = JSON.parse(
                JSON.stringify(Object.assign({}, route.queries, { stashName: undefined, namespace: undefined }))
              );
              router.navigate({}, newRouteQueies);
            }}
          >
            {t('取消')}
          </Button>
          <TipInfo
            isShow={clusterList.fetched === true && (failed || !canCreate)}
            className="error"
            style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
          >
            {clusterList && !canCreate ? tip : getWorkflowError(modifyLogStashFlow)}
            {!canCreate && ifLogDaemonset ? (
              ifLogDaemonsetNeedLoading ? (
                <Button type="icon" icon="loading" style={{ background: 'transparent' }} />
              ) : (
                <Button
                  type="icon"
                  icon="refresh"
                  style={{ background: 'transparent' }}
                  onClick={() => {
                    actions.logDaemonset.fetch();
                  }}
                >
                  （点击刷新状态）
                </Button>
              )
            ) : null}
          </TipInfo>
        </FormPanel.Footer>
      </FormPanel>
    );
  }

  private async _handleSubmit(mode) {
    let { actions, route, logStashEdit, clusterVersion, logSelection, clusterSelection } = this.props;
    // let { rid, clusterId } = route.queries;
    let { rid } = route.queries;
    let { logAgentName } = clusterSelection[0].spec;
    let { name: clusterId } = clusterSelection[0].metadata;
    let {
      logStashName,
      logMode,
      consumerMode,
      isSelectedAllNamespace,
      containerLogs,
      metadatas,
      nodeLogPath,
      addressIP,
      addressPort,
      topic,
      esAddress,
      indexName,
      containerFileNamespace,
      containerFileWorkload,
      containerFileWorkloadType,
      containerFilePaths
    } = logStashEdit;

    // 创建日志的命名空间
    let namespace =
      mode === 'update'
        ? route.queries['namespace']
        : logMode === 'containerFile'
        ? containerFileNamespace
        : 'kube-system';

    // 验证所有的编辑内容是否合法
    actions.validate.validateLogStashEdit();

    let valResult = await validatorActions._validateLogStashEdit(
      logStashEdit,
      namespace,
      clusterVersion,
      clusterId,
      mode,
      +rid
    );

    if (valResult) {
      // let { rid, clusterId } = route.queries;
      let { rid } = route.queries;

      let logResourceInfo = resourceConfig(clusterVersion)['logcs'];

      let inputType;
      let outputType;
      // 处理日志源的相关的选项
      if (logMode === 'container') {
        // 这里需要去判断 是否选择了所有的Namespace，如果是所有的namespace，则namespaces为空数组即可
        let namespaces = [];
        if (isSelectedAllNamespace === 'selectOne') {
          namespaces = containerLogs.map(
            (containerLog): ContainerLogNamespace => {
              let workloads = [];
              let workloadTypeKeys = Object.keys(containerLog.workloadSelection);
              workloadTypeKeys.forEach(item => {
                containerLog.workloadSelection[item].forEach((workloadItem: string) => {
                  workloads.push({
                    name: workloadItem,
                    type: item
                  });
                });
              });
              return {
                namespace: containerLog.namespaceSelection.replace(new RegExp(`^${clusterId}-`), ''),
                all_containers: containerLog.collectorWay === 'container',
                workloads
              };
            }
          );
          // 按照v1.3的日志采集规范，指定容器只允许设置一个ns，这里要把namespace改写成这里指定的具体的ns，而不是kube-system
          namespace = namespaces[0] && namespaces[0].namespace;
        }
        let containerLogInput: ContainerLogInput = {
          container_log_input: {
            all_namespaces: isSelectedAllNamespace === 'selectOne' ? false : true,
            namespaces
          },
          type: 'container-log'
        };
        inputType = containerLogInput;
      } else if (logMode === 'node') {
        let labels = {};
        metadatas.forEach(item => {
          labels[item.metadataKey] = item.metadataValue;
        });
        let hostLogInput: HostLogInput = {
          host_log_input: {
            labels,
            path: nodeLogPath
          },
          type: 'host-log'
        };
        inputType = hostLogInput;
      } else if (logMode === 'containerFile') {
        //聚合 将containerFilePaths中相同cantainerName的containerFilePath聚合在一块成数组
        let containerFiles = {};
        containerFilePaths.forEach(item => {
          if (containerFiles[item.containerName]) {
            containerFiles[item.containerName].push({ path: item.containerFilePath });
          } else {
            containerFiles[item.containerName] = [];
            containerFiles[item.containerName].push({ path: item.containerFilePath });
          }
        });

        let podLogInput: PodLogInput = {
          pod_log_input: {
            container_log_files: containerFiles,
            metadata: true,
            workload: {
              name: containerFileWorkload,
              type: containerFileWorkloadType
            }
          },
          type: 'pod-log'
        };

        inputType = podLogInput;
      }

      // 处理消费端的相关配置
      if (consumerMode === 'kafka') {
        let kafkaOutput: KafkaOutpot = {
          kafka_output: {
            host: addressIP,
            port: +addressPort,
            topic: topic
          },
          type: 'kafka'
        };
        outputType = kafkaOutput;
      } else if (consumerMode === 'es') {
        // 这里的和 kafka的类型一样
        let [scheme, address] = esAddress.split('://');
        let esOutput: ElasticsearchOutput = {
          elasticsearch_output: {
            hosts: [address],
            index: indexName
          },
          type: 'elasticsearch'
        };
        outputType = esOutput;
      }

      let logStashEditYaml: LogStashEditYaml = {
        kind: logResourceInfo.headTitle,
        apiVersion: (logResourceInfo.group ? logResourceInfo.group + '/' : '') + logResourceInfo.version,
        metadata: {
          name: logStashName,
          namespace: namespace.replace(new RegExp(`^${clusterId}-`), '')
        },
        spec: {
          input: inputType,
          output: outputType
        }
      };

      //更新方式为put，需要添加resoureVersion
      if (mode === 'update' && logSelection[0]) {
        logStashEditYaml.metadata.resourceVersion = logSelection[0].metadata.resourceVersion;
      }

      //容器文件需要添加label
      if (logMode === 'containerFile') {
        logStashEditYaml.metadata.labels = {
          'log.tke.cloud.tencent.com/pod-log': 'true'
        };
      }

      let jsonData = JSON.stringify(logStashEditYaml);
      let resource: CreateResource = {
        id: uuid(),
        resourceInfo: logResourceInfo,
        mode: mode === 'update' ? 'modify' : mode, // 更新方式为put，不是patch，update对应的为patch，modify对应为put
        namespace: namespace.replace(new RegExp(`^${clusterId}-`), ''),
        clusterId,
        logAgentName,
        jsonData,
        isStrategic: false,
        resourceIns: mode === 'update' ? logStashName : '' // 更新的需要需要带上具体的name
      };
      actions.workflow.modifyLogStash.start([resource], +rid);
      actions.workflow.modifyLogStash.perform();
    }
  }
}
