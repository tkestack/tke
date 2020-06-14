import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Button, Justify, SearchBox, Table, Text } from '@tencent/tea-component';

import { dateFormatter, downloadCsv } from '../../../../helpers';
import { SelectList } from '../../../../src/modules/common';
import { Cluster } from '../../common/models';
import { includes } from '../../common/utils';
import { allActions } from '../actions';
import {
    canCreateLogStash, canCreateLogStashInLogDaemonset, canFetchLogList, collectorStatus, logModeMap
} from '../constants/Config';
import { Log } from '../models';
import { LogDaemonSetStatus } from '../models/LogDaemonset';
import { router } from '../router';
import { RootProps } from './LogStashApp';

/** 新建日志收集规则按钮的提示内容 */
export const canCreateTip = {
  empty: t('请先创建一个集群'),
  canNotCreate: t('当前集群状态下无法新建日志采集规则'),
  max: t('当前集群最多可创建 100 个日志采集规则'),
  canNotCreateInLogDaemonset: phase => t(`当前日志采集器（${phase}）状态下无法创建日志采集规则`)
};

/**是否能搜索的提示内容 */
export const canSearchTip = phase => {
  return phase === '404' ? t('请先开通日志采集规则') : t(`当前日志采集器（${phase}）不健康`);
};

/** 判断当前集群能够进行新建日志采集规则的操作 */
export const isCanCreateLogStash = (clusterInfo: Cluster, logList: Log[], isDaemonsetNormal: LogDaemonSetStatus) => {
  let canCreate = false,
    tip = '',
    ifLogDaemonset = false;
  if (clusterInfo) {
    canCreate = includes(canCreateLogStash, clusterInfo.status.phase);
    tip = !canCreate ? canCreateTip.canNotCreate : '';
    if (canCreate) {
      // 兼容新旧日志组件并存
      canCreate = clusterInfo.spec.logAgentName || includes(canCreateLogStashInLogDaemonset, isDaemonsetNormal.phase);
      !canCreate && (tip = canCreateTip.canNotCreateInLogDaemonset(isDaemonsetNormal.phase));
      ifLogDaemonset = !canCreate; //标记是否是因为logDaemonset的状态不是runnig所以才不能创建的
    }
  } else {
    canCreate = false;
    tip = canCreateTip.empty;
  }

  // 目前限制一个集群下 日志采集规则为100个
  if (logList.length >= 100) {
    canCreate = false;
    tip = canCreateTip.max;
  }

  return { canCreate, tip, ifLogDaemonset };
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class LogStashActionPanel extends React.Component<RootProps, any> {
  render() {
    let {
      actions,
      clusterSelection,
      logQuery,
      logList,
      isDaemonsetNormal,
      route,
      namespaceList,
      namespaceSelection
    } = this.props;

    let logAgentName = '';
    if (clusterSelection && clusterSelection[0]) {
      logAgentName = clusterSelection[0].spec.logAgentName;
    }
    // 判断当前是否能够新建日志收集规则
    let { canCreate, tip } = isCanCreateLogStash(clusterSelection[0], logList.data.records, isDaemonsetNormal);
    let ifFetchLogList = includes(canFetchLogList, isDaemonsetNormal.phase);
    let handleNamespaceSwitched = namespaceSelection => {
      let namespaceFound = namespaceList.data.records.find(item => item.metadata && item.metadata.name === namespaceSelection);
      console.log('namespaceFound = ', namespaceFound);
    };
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

                <SelectList
                  value={namespaceSelection}
                  recordData={namespaceList}
                  valueField="namespace"
                  textField="namespace"
                  className="tc-15-select m"
                  onSelect={value => {
                    actions.namespace.selectNamespace(value);
                    handleNamespaceSwitched(value);
                    if (ifFetchLogList) {
                      actions.log.applyFilter({
                        clusterId: route.queries['clusterId'],
                        logAgentName,
                        namespace: value
                      });
                    } else {
                      actions.log.fetch({
                        noCache: true
                      });
                    }
                  }}
                  name="Namespace"
                  tipPosition="right"
                  style={{
                    display: 'inline-block',
                    padding: '0 10px'
                  }}
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

  /** 处理新建按钮的button */
  private _handleCreate() {
    let { actions, isOpenLogStash, route, clusterSelection } = this.props,
      urlParams = router.resolve(route);
    if (clusterSelection && clusterSelection[0] && clusterSelection[0].spec.logAgentName || isOpenLogStash) {
      router.navigate(Object.assign({}, urlParams, { mode: 'create' }), route.queries);
    } else {
      actions.workflow.authorizeOpenLog.start();
    }
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
