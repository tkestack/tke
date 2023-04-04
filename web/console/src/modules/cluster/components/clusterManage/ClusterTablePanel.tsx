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
import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, Button, Drawer, Dropdown, Icon, List, Text } from '@tea/component';
import { TablePanel, TablePanelColumnProps } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { ChartPanel } from '@tencent/tchart';
import { Trans, t } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../../helpers';
import { Clip, LinkButton, TipInfo } from '../../../common/components';
import { Cluster } from '../../../common/models';
import { allActions } from '../../actions';
import { ClusterTypeMap } from '../../constants/Config';
import { DialogNameEnum } from '../../models';
import { MonitorPanelProps, getClusterTables } from '../../models/MonitorPanel';
import { router } from '../../router';
import { RootProps } from '../ClusterApp';
import { KubectlDialog } from '../KubectlDialog';
import { UpdateClusterTokenDialog } from './UpdateClusterTokenDialog';

/** 集群的状态颜色的展示 */
export const ClusterStatus = {
  Running: 'success',
  Initializing: 'label',
  Failed: 'danger',
  Terminating: 'label'
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

interface State {
  showOsTips?: boolean;
  selectCluster?: Cluster;
  monitorPanelProps?: MonitorPanelProps;
  isShowMonitor?: boolean;
}

@connect(state => state, mapDispatchToProps)
export class ClusterTablePanel extends React.Component<RootProps, State> {
  state = {
    showOsTips: false,
    selectCluster: null,
    monitorPanelProps: undefined,
    isShowMonitor: false
  };

  componentWillUnmount() {
    const { actions } = this.props;
    actions.cluster.clearPolling();
  }

  render() {
    return (
      <React.Fragment>
        {this._renderTablePanel()}
        <KubectlDialog {...this.props} />
        <UpdateClusterTokenDialog />
        {this.renderMonitor()}
      </React.Fragment>
    );
  }

  private _renderTablePanel() {
    const { actions, cluster, region } = this.props;

    const columns: TablePanelColumnProps<Cluster>[] = [
      {
        key: 'name',
        header: t('ID/名称'),
        width: '10%',
        render: x => (
          <React.Fragment>
            <Text parent="div" className="m-width" overflow>
              {x.status.phase.toLowerCase() !== 'running' ? (
                x.metadata.name
              ) : (
                <React.Fragment>
                  <a
                    id={x.metadata.name}
                    title={x.metadata.name}
                    href="javascript:;"
                    onClick={() => {
                      this._handleClickForCluster(x);
                    }}
                    className="tea-text-overflow"
                  >
                    {x.metadata.name || '-'}
                  </a>
                </React.Fragment>
              )}
            </Text>
            <Clip target={`#${x.metadata.name}`} />
            <Text parent="div">
              {x.spec.displayName || '-'}
              <Icon
                onClick={() => {
                  actions.cluster.selectCluster([x]);
                  actions.workflow.modifyClusterName.start([x], 1);
                }}
                style={{ cursor: 'pointer' }}
                type="pencil"
              />
            </Text>
          </React.Fragment>
        )
      },
      {
        key: 'monitor',
        header: t('监控'),
        width: '7%',
        render: x => {
          const { promethus } = x.spec;
          return (
            <div>
              <p className="text-overflow m-width">
                <i
                  className="dosage-icon"
                  style={{ cursor: promethus ? 'pointer' : 'not-allowed' }}
                  data-monitor
                  data-title={promethus ? t('查看监控') : '监控告警组件未开启'}
                  onClick={() => {
                    promethus && this._handleMonitor(x);
                  }}
                />
                {/* {!x.clusterBMonitor && <span className="alarm-label-tips">{t('未配告警')}</span>} */}
              </p>
            </div>
          );
        }
      },
      {
        key: 'status',
        header: t('状态'),
        width: '8%',
        render: x => (
          <React.Fragment>
            <Text theme={ClusterStatus[x.status.phase]} verticalAlign="middle">
              {x.status.phase || '-'}
            </Text>
            {x.spec.updateInfo.worker.message === '有节点正在升级中' && (
              <>
                <Icon className="tea-ml-1n" type="loading" />{' '}
                <Button
                  type="link"
                  onClick={() => {
                    router.navigate(
                      { sub: 'sub', mode: 'list', type: 'nodeManage', resourceName: 'node' },
                      {
                        rid: region?.selection?.value + '',
                        clusterId: x.metadata.name
                      }
                    );
                  }}
                >
                  有节点正在升级中
                </Button>
              </>
            )}
            {x.status.phase !== 'Running' && <Icon className="tea-ml-1n" type="loading" />}
            {x.status.phase !== 'Running' && x.status.phase !== 'Terminating' && (
              <Button
                type="link"
                onClick={() => {
                  actions.cluster.select(x);
                  actions.dialog.updateDialogState(DialogNameEnum.clusterStatusDialog);
                }}
              >
                查看详情
              </Button>
            )}
          </React.Fragment>
        )
      },
      {
        key: 'version',
        header: t('K8S版本'),
        width: '8%',
        render: x => <Text>{x.status.version || '-'}</Text>
      },
      {
        key: 'type',
        header: t('集群类型'),
        width: '8%',
        render: x => <Text>{ClusterTypeMap[(x.spec.type as string).toLowerCase()]}</Text>
      },
      {
        key: 'createTime',
        header: t('创建时间'),
        width: '15%',
        render: x => this._reduceTime(x.metadata.creationTimestamp)
      },
      { key: 'operation', header: t('操作'), width: '16%', render: x => this._renderOperationCell(x) }
    ];

    const emptyTips: JSX.Element = (
      <div className="text-center">
        <Trans>
          当前集群列表为空，您可以
          <a
            href="javascript:void(0);"
            onClick={() => router.navigate({ sub: 'createIC' }, { rid: region?.selection?.value + '' })}
          >
            [新建一个集群]
          </a>
        </Trans>
      </div>
    );

    return (
      <TablePanel
        columns={columns}
        emptyTips={emptyTips}
        action={actions.cluster}
        model={cluster}
        bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
        rowDisabled={record => record.status.phase === 'Terminating'}
      />
    );
  }

  private _handleMonitor(cluster) {
    this.setState({
      isShowMonitor: true,
      selectCluster: cluster,
      monitorPanelProps: {
        title: cluster.metadata.name,
        subTitle: cluster.spec.displayName,
        tables: getClusterTables(cluster.metadata.name),
        groupBy: []
      }
    });
  }
  private renderMonitor() {
    return (
      <Drawer
        visible={this.state.isShowMonitor}
        title={(this.state.monitorPanelProps && this.state.monitorPanelProps.title) || ''}
        subTitle={(this.state.monitorPanelProps && this.state.monitorPanelProps.subTitle) || ''}
        onClose={() => this.setState({ isShowMonitor: false })}
        outerClickClosable={true}
        placement={'right'}
        size={'l'}
        style={{ zIndex: 4 }}
      >
        {this.renderPromTip()}
        {this.state.monitorPanelProps && (
          <ChartPanel
            tables={this.state.monitorPanelProps.tables}
            groupBy={this.state.monitorPanelProps.groupBy}
            height={250}
          />
        )}
      </Drawer>
    );
  }

  private renderPromTip() {
    const { selectCluster } = this.state;
    return (
      selectCluster &&
      !selectCluster.spec.hasPrometheus && (
        <TipInfo className="warning">
          <span style={{ verticalAlign: 'middle' }}>
            <Trans>
              该集群未安装Prometheus组件, 请前往
              <a href={`/tkestack/cluster/sub/list/basic/info?clusterId=${selectCluster.selection.metadata.name}`}>
                集群基本信息
              </a>
              进行安装
            </Trans>
          </span>
        </TipInfo>
      )
    );
  }

  /** 处理创建时间 */
  private _reduceTime(time = '') {
    const showTime = dateFormatter(new Date(time), 'YYYY-MM-DD HH:mm:ss');

    return (
      <Bubble placement="left" content={showTime || null}>
        <Text>{showTime}</Text>
      </Bubble>
    );
  }

  /** 处理集群点击跳转 */
  private _handleClickForCluster(cluster: Cluster) {
    const { actions, region } = this.props;

    // 进行路由的跳转
    const routeQueries = {
      rid: region?.selection?.value + '',
      clusterId: cluster.metadata.name
    };
    router.navigate({ sub: 'sub', mode: 'list', type: 'resource', resourceName: 'deployment' }, routeQueries);
    // 进行deployment数据的拉取
    actions.resource.initResourceInfoAndFetchData(true, 'deployment');

    // 选择当前选中的集群信息, true 即需要初始化k8s的版本
    actions.cluster.selectCluster([cluster], true);
  }

  /** 渲染操作按钮 */
  private _renderOperationCell(cluster: Cluster) {
    const { actions, region } = this.props;
    const isDisabledButon = cluster.status.phase === 'Terminating';

    const routeQueries = {
      rid: region?.selection?.value + '',
      clusterId: cluster.metadata.name
    };

    const renderConditions = () => {
      return (
        <LinkButton
          disabled={isDisabledButon}
          onClick={() => {
            actions.cluster.select(cluster);
            actions.dialog.updateDialogState(DialogNameEnum.clusterStatusDialog);
          }}
        >
          查看创建详情
        </LinkButton>
      );
    };

    const renderDeleteButton = () => {
      const isCanNotDelete = cluster.metadata.name === 'global' || isDisabledButon;
      return (
        <LinkButton
          disabled={isCanNotDelete}
          errorTip={isDisabledButon ? '' : 'global集群不可删除'}
          tipDirection="left"
          onClick={() => {
            if (!isCanNotDelete) {
              actions.workflow.deleteCluster.start([cluster]);
            }
          }}
        >
          删除
        </LinkButton>
      );
    };

    const renderKuberctlButton = () => {
      return (
        <LinkButton
          disabled={isDisabledButon}
          tipDirection="left"
          onClick={() => {
            actions.cluster.select(cluster);
            actions.cluster.fetchClustercredential(cluster.metadata.name);
            actions.dialog.updateDialogState(DialogNameEnum.kuberctlDialog);
          }}
        >
          {t('查看集群凭证')}
        </LinkButton>
      );
    };

    const renderUpdateTokenButton = () => {
      return (
        <Button
          disabled={isDisabledButon}
          type="link"
          style={{ marginLeft: '5px' }}
          onClick={() => {
            actions.cluster.fetchClustercredential(cluster.metadata.name);
            actions.workflow.updateClusterToken.start([]);
          }}
        >
          {t('修改集群凭证')}
        </Button>
      );
    };

    const renderMoreButton = () => {
      return (
        <Dropdown
          trigger="hover"
          style={{ marginRight: '5px' }}
          button="更多"
          onOpen={() => console.log('open')}
          onClose={() => console.log('close')}
        >
          <List type="option">
            <List.Item>
              <LinkButton
                disabled={!cluster.spec.updateInfo.master.isNeed || cluster.status.phase !== 'Running'}
                errorTip={cluster.spec.updateInfo.master.message}
                onClick={() => {
                  router.navigate(
                    { sub: 'cluster-update' },
                    { ...routeQueries, clusterVersion: cluster?.status?.version }
                  );
                }}
              >
                升级Master
              </LinkButton>
            </List.Item>

            <List.Item>
              <LinkButton
                disabled={!cluster.spec.updateInfo.worker.isNeed || cluster.status.phase !== 'Running'}
                errorTip={cluster.spec.updateInfo.worker.message}
                onClick={() => {
                  router.navigate(
                    { sub: 'worker-update' },
                    { ...routeQueries, clusterVersion: cluster?.status?.version }
                  );
                }}
              >
                升级Worker
              </LinkButton>
            </List.Item>
          </List>
        </Dropdown>
      );
    };

    return (
      <React.Fragment>
        {/* {cluster.status.phase !== 'Running' && cluster.status.phase !== 'Terminating' && renderConditions()} */}
        {renderDeleteButton()}
        {renderKuberctlButton()}
        {/* {cluster.spec.type === 'Imported' && renderUpdateTokenButton()} */}
        {/* {cluster.spec.type === 'Imported' && renderUpdateTokenButton()} */}
        {cluster.spec.type !== 'Imported' && renderMoreButton()}
      </React.Fragment>
    );
  }
}
