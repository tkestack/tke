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
import { RootProps } from './AlarmPolicyApp';
import { LinkButton } from '../../common/components';
import { MetricsObject, AlarmPolicy } from '../models/AlarmPolicy';
import { router } from '../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { MetricNameMap } from '../constants/Config';
import { router as notifyRouter } from '../../notify/router';
import { FormPanel } from '@tencent/ff-component';

export class AlarmPolicyDetailPanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, route, alarmPolicyDetail, channel, template, receiverGroup } = this.props;
    return (
      <FormPanel
        title={t('基本信息')}
        operation={
          <LinkButton
            onClick={() => {
              router.navigate(
                { sub: 'update' },
                Object.assign({}, route.queries, { alarmPolicyId: alarmPolicyDetail.id })
              );
            }}
          >
            {t('编辑')}
          </LinkButton>
        }
      >
        <FormPanel.Item text label={t('告警策略名称')}>
          {alarmPolicyDetail.alarmPolicyName}
        </FormPanel.Item>
        {/* <FormPanel.Item text label={t('备注')}>
          {alarmPolicyDetail.alarmPolicyDescription || '-'}
        </FormPanel.Item> */}
        <FormPanel.Item text label={t('策略类型')}>
          {alarmPolicyDetail.alarmPolicyType === 'cluster'
            ? t('集群')
            : alarmPolicyDetail.alarmPolicyType === 'node'
            ? t('节点')
            : 'Pod'}
        </FormPanel.Item>
        <FormPanel.Item text label={t('对象类型')}>
          {alarmPolicyDetail.alarmObjetcsType === 'all'
            ? t('全部选择')
            : alarmPolicyDetail.alarmObjetcsType === 'part'
            ? t('按工作负载选择')
            : t('按k8sLabel选择')}
        </FormPanel.Item>
        <FormPanel.Item text label={t('触发条件')}>
          {alarmPolicyDetail.alarmMetrics && this._rendAlarmMetrics(alarmPolicyDetail.alarmMetrics, alarmPolicyDetail)}
        </FormPanel.Item>
        <FormPanel.Item text label={t('接收组')}>
          {alarmPolicyDetail.receiverGroups &&
            alarmPolicyDetail.receiverGroups.map((gid, index) => {
              let group = receiverGroup.list.data.records.find(g => g.metadata.name === gid);
              return (
                <React.Fragment key={index}>
                  {index > 0 && ', '}
                  <LinkButton
                    disabled={!group}
                    onClick={() => {
                      notifyRouter.navigate(
                        {
                          mode: 'detail',
                          resourceName: 'receiverGroup'
                        },
                        { resourceIns: gid }
                      );
                    }}
                    className="tea-text-overflow"
                  >
                    {gid}
                  </LinkButton>
                  {group && `(${group.spec.displayName})`}
                </React.Fragment>
              );
            })}
        </FormPanel.Item>
        <FormPanel.Item text label={t('接收方式')}>
          {alarmPolicyDetail.notifyWays &&
            alarmPolicyDetail.notifyWays.map(notifyWay => {
              let c = channel.list.data.records.find(c => c.metadata.name === notifyWay.channel);
              let tp = template.list.data.records.find(c => c.metadata.name === notifyWay.template);
              return (
                <p key={notifyWay.id}>
                  {t('渠道')}:
                  <LinkButton
                    disabled={!c}
                    title={notifyWay.channel}
                    onClick={() => {
                      notifyRouter.navigate(
                        {
                          mode: 'detail',
                          resourceName: 'channel'
                        },
                        { resourceIns: notifyWay.channel }
                      );
                    }}
                    className="tea-text-overflow"
                  >
                    {notifyWay.channel}
                  </LinkButton>
                  {c && `(${c.spec.displayName})`} {t('模版')}:
                  <LinkButton
                    title={notifyWay.channel}
                    disabled={!tp}
                    onClick={() => {
                      notifyRouter.navigate(
                        {
                          mode: 'detail',
                          resourceName: 'template'
                        },
                        { resourceIns: notifyWay.template }
                      );
                    }}
                    className="tea-text-overflow"
                  >
                    {notifyWay.template}
                  </LinkButton>
                  {tp && `(${tp.spec.displayName})`}
                </p>
              );
            })}
        </FormPanel.Item>
      </FormPanel>
    );
  }
  private _rendAlarmMetrics(alarmMetrics: MetricsObject[], alarmPolicy: AlarmPolicy) {
    let len = alarmMetrics.length;
    let content: JSX.Element[] = [];
    for (let i = 0; i < len; ++i) {
      let evaluator = alarmMetrics[i].type === 'boolean' ? '=' : alarmMetrics[i].evaluatorType === 'gt' ? '>' : '<';
      content.push(
        <p>
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
    }
    return <div>{content}</div>;
  }
}
