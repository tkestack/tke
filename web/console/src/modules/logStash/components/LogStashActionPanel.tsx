import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Button, Justify, SearchBox, Select, Table, Text } from '@tencent/tea-component';

import { dateFormatter, downloadCsv } from '../../../../helpers';
import { SelectList } from '../../../../src/modules/common';
import { Cluster } from '../../common/models';
import { includes } from '../../common/utils';
import { allActions } from '../actions';
import {
  canCreateLogStash,
  canCreateLogStashInLogDaemonset,
  canFetchLogList,
  collectorStatus,
  logModeMap
} from '../constants/Config';
import { Log } from '../models';
import { LogDaemonSetStatus } from '../models/LogDaemonset';
import { router } from '../router';
import { RootProps } from './LogStashApp';

/** 新建日志收集规则按钮的提示内容 */
export const canCreateTip = {
  empty: t('请先创建或者选择一个集群'),
  canNotCreate: t('当前集群状态下无法新建日志采集规则'),
  max: t('当前集群最多可创建 100 个日志采集规则'),
  canNotCreateInLogDaemonset: phase => t(`当前日志采集器（${phase}）状态下无法创建日志采集规则`)
};

/**是否能搜索的提示内容 */
export const canSearchTip = phase => {
  return phase === '404' ? t('请先开通日志采集规则') : t(`当前日志采集器（${phase}）不健康`);
};

export const isCanCreateLogStash = (
  clusterInfo: Cluster,
  logList: Log[],
  isDaemonsetNormal: LogDaemonSetStatus,
  isOpenLogStash
) => {
  let canCreate = true,
    tip = '',
    ifLogDaemonset = false;
  // 没有选择集群
  if (!clusterInfo) {
    canCreate = false;
    tip = canCreateTip.empty;
  } else if (!includes(canCreateLogStash, clusterInfo.status.phase)) {
    // 集群没有运行
    canCreate = false;
    tip = canCreateTip.canNotCreate;
  } else if (!clusterInfo.spec.logAgentName && !isOpenLogStash) {
    // 日志组件是否安装
    canCreate = false;
    tip = '日志组件尚未安装，请先安装LogAgent组件';
  } else if (
    !(
      (clusterInfo.spec.logAgentName && clusterInfo.spec.logAgentStatus === 'Running') ||
      (isOpenLogStash && includes(canCreateLogStashInLogDaemonset, isDaemonsetNormal.phase))
    )
  ) {
    // 日志组件的状态是否正常。安装了logAgent并且状态是运行的，或者安装了logCollector但是状态是运行的都可以创建
    canCreate = false;
    if (clusterInfo.spec.logAgentName) {
      tip = canCreateTip.canNotCreateInLogDaemonset(clusterInfo.spec.logAgentStatus);
    } else {
      tip = canCreateTip.canNotCreateInLogDaemonset(isDaemonsetNormal.phase);
    }
    ifLogDaemonset = true; //标记是否是因为日志组件的状态不是runnig所以才不能创建的
  } else if (logList.length >= 100) {
    // 目前限制一个集群下 日志采集规则为100个
    canCreate = false;
    tip = canCreateTip.max;
  }

  return { canCreate, tip, ifLogDaemonset };
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class LogStashActionPanel extends React.Component<RootProps, any> {
  state = {
    isBusiness: window.location.href.includes('tkestack-project'),
    namespace: this.props.namespaceSelection
  };

  render() {
    let {
      actions,
      clusterSelection,
      logQuery,
      logList,
      isOpenLogStash,
      isDaemonsetNormal,
      route,
      namespaceList: namespaceRecords,
      namespaceSelection
    } = this.props;
    let { isBusiness } = this.state;

    let logAgentName = '';
    if (clusterSelection && clusterSelection[0]) {
      logAgentName = clusterSelection[0].spec.logAgentName;
    }
    // 判断当前是否能够新建日志收集规则
    let { canCreate, tip } = isCanCreateLogStash(
      clusterSelection[0],
      logList.data.records,
      isDaemonsetNormal,
      isOpenLogStash
    );
    let ifFetchLogList = logAgentName || includes(canFetchLogList, isDaemonsetNormal.phase);
    let handleNamespaceSwitched = namespaceSelection => {
      let namespaceFound = namespaceRecords.data.records.find(item => item.namespace === namespaceSelection);
      actions.cluster.selectClusterFromNamespace(namespaceFound.cluster);
    };
    let namespaces = [];
    let namespaceList = [];
    let groups = {};
    if (!isBusiness) {
      namespaces = namespaceRecords.data.records.map(item => ({
        name: item.name,
      }));
      namespaceList = namespaces.map(({ name }) => ({
        // value: fullName,
        value: name,
        text: name
      }));
    } else {
      namespaces = namespaceRecords.data.records.map(item => ({
        name: item.namespace,
        fullName: item.namespaceValue,
        clusterName: item.cluster.metadata.name,
        clusterDisplayName: item.cluster.spec.displayName
      }));
      namespaceList = namespaces.map(({ name, clusterName, fullName }) => ({
        value: fullName,
        groupKey: clusterName,
        text: name
      }));
      groups = namespaces.reduce((accu, item, index, arr) => {
        let { clusterName, clusterDisplayName } = item;
        if (!accu[clusterName]) {
          accu[clusterName] = clusterDisplayName;
        }
        return accu;
      }, {});
    }

    return (
      <Table.ActionPanel>
        <Justify
          left={
            <React.Fragment>
              <Bubble content={!canCreate ? tip : null}>
                <Button type="primary" disabled={!canCreate} onClick={this._handleCreate.bind(this)}>
                  {t('新建')}
                </Button>
              </Bubble>
            </React.Fragment>
          }
          right={
            <React.Fragment>
              <div style={{ display: 'inline-block', fontSize: '12px' }}>
                <Text theme="label" style={{ verticalAlign: '-4px' }}>
                  {t('命名空间')}
                </Text>

                <Select
                  searchable
                  boxSizeSync
                  // groups={isPlatform ? undefined : groups}
                  size="m"
                  type="simulate"
                  appearence="button"
                  options={namespaceList}
                  groups={groups}
                  value={namespaceSelection}
                  onChange={this.handleNamespaceChanged}
                />
              </div>
              <Bubble
                content={
                  !namespaceSelection
                    ? '请先选择一个命名空间'
                    : !ifFetchLogList
                    ? canSearchTip(isDaemonsetNormal.phase)
                    : null
                }
              >
                <SearchBox
                  style={{ marginLeft: 8 }}
                  value={logQuery.keyword ? logQuery.keyword : ''}
                  onChange={actions.log.changeKeyword}
                  onSearch={value => {
                    if (ifFetchLogList) {
                      actions.log.performSearch(value);
                    } else {
                      actions.log.fetch({
                        noCache: true
                      });
                    }
                  }}
                  onClear={() => {
                    if (ifFetchLogList) {
                      actions.log.performSearch('');
                    } else {
                      actions.log.fetch({
                        noCache: true
                      });
                    }
                  }}
                  placeholder={t('请输入日志名称')}
                  disabled={!namespaceSelection}
                />
              </Bubble>
              <Button
                icon="download"
                title={t('下载')}
                onClick={() => this._downloadHandle(this.props.logList.data.records, isDaemonsetNormal)}
              />
            </React.Fragment>
          }
        />
      </Table.ActionPanel>
    );
  }

  getList = namespaceValue => {
    let { actions, namespaceList, route, namespaceSelection, isDaemonsetNormal, clusterSelection } = this.props;
    let { isBusiness } = this.state;
    let {
      namespace: { selectNamespace },
      cluster: { selectClusterFromNamespace },
      log: { applyFilter, fetch }
    } = actions;
    let logAgentName = '';
    if (clusterSelection && clusterSelection[0]) {
      logAgentName = clusterSelection[0].spec.logAgentName;
    }
    // 如果是平台侧的话，切换集群的时候地址栏中的clusterId参数已经体现出当前选中的集群了
    let clusterId = route.queries['clusterId'];
    // 兼容业务侧对集群的处理，从命名空间关联到集群
    let namespace = namespaceValue;
    if (isBusiness) {
      let namespaceFound = namespaceList.data.records.find(item => item.namespaceValue === namespaceValue);
      if (namespaceFound !== undefined) {
        // 业务侧下再覆盖namespace
        namespace = namespaceFound.namespace;
        // 取附加到ns上的集群信息中的集群id
        clusterId = namespaceFound.cluster.metadata.name;
        logAgentName = namespaceFound.cluster.spec.logAgentName;
        selectClusterFromNamespace(namespaceFound.cluster);
      }
    }
    let ifFetchLogList = logAgentName || includes(canFetchLogList, isDaemonsetNormal.phase);
    if (ifFetchLogList) {
      applyFilter({
        clusterId,
        logAgentName,
        namespace: namespace
      });
    } else {
      fetch({
        noCache: true
      });
    }
  };

  private handleNamespaceChanged = value => {
    let { actions, namespaceList, route, namespaceSelection, isDaemonsetNormal, clusterSelection } = this.props;
    let {
      namespace: { selectNamespace },
      cluster: { selectClusterFromNamespace },
      log: { applyFilter, fetch }
    } = actions;
    this.setState({ namespace: value });
    selectNamespace(value);

    this.getList(value);
  };

  /** 处理新建按钮的button */
  private _handleCreate() {
    let { actions, isOpenLogStash, route, clusterSelection } = this.props,
      urlParams = router.resolve(route);
    if ((clusterSelection && clusterSelection[0] && clusterSelection[0].spec.logAgentName) || isOpenLogStash) {
      router.navigate(Object.assign({}, urlParams, { mode: 'create' }), route.queries);
    } else {
      actions.workflow.authorizeOpenLog.start();
    }
  }

  private _handleInstall() {
    let { actions, isOpenLogStash, route, clusterSelection } = this.props,
      urlParams = router.resolve(route);
    actions.cluster.enableLogAgent(clusterSelection[0]);
    actions.cluster.applyFilter({});
  }

  private _handleDelete() {
    let { actions, isOpenLogStash, route, clusterSelection } = this.props,
      urlParams = router.resolve(route);
    actions.cluster.disableLogAgent(clusterSelection[0]);
    actions.cluster.applyFilter({});
  }

  /** 下载操作 */
  private _downloadHandle(logList: Log[], isDaemonsetNormal: LogDaemonSetStatus) {
    let rows = [],
      head = [t('名称'), t('状态'), t('日志类型'), t('命名空间'), t('创建时间')];

    logList.forEach((item: Log) => {
      let row = [
        item.metadata.name,
        collectorStatus[isDaemonsetNormal.phase].text,
        logModeMap[item.spec.input.type],
        item.metadata.namespace,
        dateFormatter(new Date(item.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')
      ];
      rows.push(row);
    });

    downloadCsv(rows, head, 'tkestack.csv');
  }
}
