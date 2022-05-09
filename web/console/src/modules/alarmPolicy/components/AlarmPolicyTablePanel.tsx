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

import { Bubble, Button, TableColumn } from 'tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { router as addonRouter } from '../../addon/router';
import { LinkButton, TipInfo } from '../../common/components';
import { MetricNameMap } from '../constants/Config';
import { AlarmPolicy, MetricsObject } from '../models/AlarmPolicy';
import { router } from '../router';
import { RootProps } from './AlarmPolicyApp';
import { TablePanelColumnProps, TablePanel } from '@tencent/ff-component';

export class AlarmPolicyTablePanel extends React.Component<RootProps, {}> {
  render() {
    return this._renderTablePanel();
  }

  getColumns() {
    const { actions, route } = this.props;
    const columns: TableColumn<AlarmPolicy>[] = [
      {
        key: 'alarmPolicyName',
        header: t('告警策略名称'),
        width: '10%',
        render: x => {
          return (
            <LinkButton
              title={x.alarmPolicyName}
              onClick={() => {
                actions.alarmPolicy.fetchAlarmPolicyDetail(x);
                router.navigate({ sub: 'detail' }, Object.assign({}, route.queries, { alarmPolicyId: x.id }));
              }}
              overflow
            >
              {x.alarmPolicyName}
            </LinkButton>
          );
        }
      },
      {
        key: 'PolicyType',
        header: t('策略类型'),
        width: '8%',
        render: x => {
          const map = {
            cluster: t('集群'),
            node: t('节点'),
            pod: t('Pod'),
            virtualMachine: t('虚拟机')
          };

          return map?.[x?.alarmPolicyType] ?? '-';
        }
      },
      {
        key: 'PolicyRule',
        header: t('触发条件'),
        width: '15%',
        render: x => {
          const { content, hoverContent } = this._getAlarmMetricsContent(x.alarmMetrics, x);
          return (
            <Bubble placement="right" content={hoverContent || null}>
              {content}
            </Bubble>
          );
        }
      },
      {
        key: 'PolicyNotify',
        header: t('告警渠道'),
        width: '16%',
        render: x => {
          return this._rendAlarmNotifyWay(x);
        }
      }
    ];
    return columns;
  }

  private _renderTablePanel() {
    const { actions, alarmPolicy, cluster } = this.props;

    const columns = this.getColumns();
    const emptyTips: JSX.Element =
      cluster.selection && cluster.selection.spec.hasPrometheus ? (
        <div className="text-center">
          <Trans>
            您选择的集群的告警设置列表为空，您可以
            <Button
              type="link"
              onClick={() => {
                this._handleCreate();
              }}
            >
              新建告警设置
            </Button>
            ，或切换到其他集群
          </Trans>
        </div>
      ) : (
        <div className="text-center">
          <Trans>您选择的集群的告警设置列表为空</Trans>
        </div>
      );

    return (
      <React.Fragment>
        {this.renderPromTip()}
        <TablePanel
          left={
            <React.Fragment>
              <Button
                type="primary"
                onClick={() => this.handleCreate()}
                disabled={!(cluster.selection && cluster.selection.spec.hasPrometheus)}
              >
                {/* <b className="icon-add" /> */}
                {t('新建')}
              </Button>
              <Button
                disabled={alarmPolicy.selections.length === 0}
                onClick={() => actions.workflow.deleteAlarmPolicy.start(alarmPolicy.selections)}
              >
                {t('删除')}
              </Button>
            </React.Fragment>
          }
          columns={columns}
          emptyTips={emptyTips}
          model={alarmPolicy}
          action={actions.alarmPolicy}
          getOperations={record => this.getOperations(record)}
          selectable={{
            value: alarmPolicy.selections.map(item => item.id as string),
            onChange: keys => {
              actions.alarmPolicy.selects(
                alarmPolicy.list.data.records.filter(item => keys.indexOf(item.id as string) !== -1)
              );
            }
          }}
          isNeedPagination={true}
        />
      </React.Fragment>
    );
  }

  private renderPromTip() {
    const { cluster } = this.props;
    const showTip = cluster.selection && !cluster.selection.spec.hasPrometheus;
    return (
      showTip && (
        <TipInfo className="warning">
          <span style={{ verticalAlign: 'middle' }}>
            {
              /// #if tke
              <Trans>
                该集群未安装Prometheus组件, 请前往
                <a href={`/tkestack/cluster/sub/list/basic/info?clusterId=${cluster.selection.metadata.name}`}>
                  集群基本信息
                </a>
                进行安装
              </Trans>
              /// #endif
            }
            {
              /// #if project
              <Trans>该集群未安装Prometheus组件, 请通知平台管理员进行安装</Trans>
              /// #endif
            }
          </span>
        </TipInfo>
      )
    );
  }

  private handleCreate() {
    const { route, regionSelection, cluster } = this.props;
    //actions.mode.changeMode("expand");
    const clusterId = route.queries['clusterId'] || (cluster.selection ? cluster.selection.metadata.name : '');
    router.navigate({ sub: 'create' }, Object.assign({}, route.queries, { clusterId }));
  }
  private _getAlarmMetricsContent(alarmMetrics: MetricsObject[], alarmPolicy: AlarmPolicy) {
    const len = alarmMetrics.length;
    const hoverContent: JSX.Element[] = [];
    const content: JSX.Element[] = [];
    for (let i = 0; i < len; ++i) {
      const evaluator = alarmMetrics[i].type === 'boolean' ? '=' : alarmMetrics[i].evaluatorType === 'gt' ? '>' : '<';
      const temp = (
        <p key={i}>
          <span className="text-overflow">
            {`${MetricNameMap[alarmMetrics[i].metricName] || alarmMetrics[i].metricName}${evaluator}${
              alarmMetrics[i].type === 'boolean'
                ? +alarmMetrics[i].evaluatorValue
                  ? 'False'
                  : 'True'
                : alarmMetrics[i].evaluatorValue
            }${alarmMetrics[i].unit},` +
              t('持续{{count}}分钟告警', {
                count: (alarmMetrics[i].continuePeriod * alarmPolicy.statisticsPeriod) / 60
              })}
          </span>
        </p>
      );
      hoverContent.push(temp);
      if (i < 3) {
        content.push(temp);
      }
    }
    return { hoverContent: <div>{hoverContent}</div>, content: <div>{content}</div> };
  }

  private _rendAlarmNotifyWay(alarmPolicy: AlarmPolicy) {
    const { notifyWays, receiverGroups } = alarmPolicy;

    return (
      <div>
        <p>
          <span className="text-overflow">
            {t('接收组:{{count}}个', {
              count: receiverGroups.length
            })}
          </span>
        </p>
        <p>
          <span className="text-overflow">
            {t('接收方式:{{count}}个', {
              count: notifyWays.length
            })}
          </span>
        </p>
      </div>
    );
  }
  private _handleDeleteAlarmPolicy(alarmPolicy: AlarmPolicy) {
    const { actions } = this.props;
    actions.alarmPolicy.selects([alarmPolicy]);
    actions.workflow.deleteAlarmPolicy.start([alarmPolicy]);
  }

  private _handleCopyAlarmPolicy(alarmPolicy: AlarmPolicy) {
    const { actions, route } = this.props;
    router.navigate({ sub: 'copy' }, Object.assign({}, route.queries, { alarmPolicyId: alarmPolicy.id }));
    actions.alarmPolicy.selects([alarmPolicy]);
  }
  private _handleCreate() {
    const { route, regionSelection, cluster } = this.props;
    //actions.mode.changeMode("expand");
    const rid = route.queries['rid'] || regionSelection.value + '',
      clusterId = route.queries['clusterId'] || (cluster.selection ? cluster.selection.metadata.name : '');
    router.navigate({ sub: 'create' }, { rid, clusterId });
  }
  private getOperations(alarmPolicy: AlarmPolicy) {
    const { cluster, route } = this.props;

    const clusterId = cluster.selection ? cluster.selection.metadata.name : route.queries['clusterId'] || '';

    const renderDeleteButton = () => {
      return (
        <LinkButton key={0} onClick={() => this._handleDeleteAlarmPolicy(alarmPolicy)}>
          {t('删除')}
        </LinkButton>
      );
    };

    const renderCopyButton = () => {
      return (
        <LinkButton key={1} onClick={() => this._handleCopyAlarmPolicy(alarmPolicy)}>
          {t('复制')}
        </LinkButton>
      );
    };
    return [renderDeleteButton(), renderCopyButton()];
  }
}
