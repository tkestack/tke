/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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
import { RootProps } from '../NotifyApp';
import { insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Table, Card, Text } from '@tencent/tea-component';
insertCSS(
  'notifyPartBodyIntro',
  `.h3 { margin: 8px 0; }
`
);
export const BodyIntro = props => {
  const bodyTipList = [
    {
      field: 'startsAt',
      explanation: '开始时间'
    },
    {
      field: 'alarmPolicyType',
      explanation: '告警类型 (pod/node/cluster)'
    },
    {
      field: 'alarmPolicyName',
      explanation: '告警策略名称'
    },
    {
      field: 'alertName',
      explanation: '告警的metric (k8s_pod_rate_cpu_core_used_limit)'
    },
    {
      field: 'value',
      explanation: '告警值'
    },
    {
      field: 'workloadKind',
      explanation: 'workload类型'
    },
    {
      field: 'workloadName',
      explanation: 'workload名称'
    },
    {
      field: 'clusterID',
      explanation: '集群ID'
    },
    {
      field: 'namespace',
      explanation: '命名空间'
    },
    {
      field: 'podName',
      explanation: 'pod名称'
    },
    {
      field: 'nodeName',
      explanation: '节点名称'
    },
    {
      field: 'nodeRole',
      explanation: '节点类型 (node/master)'
    },
    {
      field: 'metricDisplayName',
      explanation: '告警指标展示名（前端展示的中文名字，如CPU利用率）'
    },
    {
      field: 'evaluateType',
      explanation: '告警阈值条件（>, <, ==）'
    },
    {
      field: 'evaluateValue',
      explanation: '告警阈值'
    },
    {
      field: 'unit',
      explanation: '告警值的单位'
    }
  ];

  return (
    <Card>
      <Card.Body>
        <h3 className="h3">
          <Trans>body模板说明</Trans>
        </h3>
        <p>{'body支持go template 格式，{{}}内表示引用的消息字段，{{}}外可以使用任何字符, 当前支持以下字段：'}</p>
        <Table
          records={bodyTipList}
          recordKey="field"
          columns={[
            {
              key: 'field',
              header: t('字段')
            },
            {
              key: 'explanation',
              header: t('说明')
            }
          ]}
        />
        <p>
          <Text>
            {
              'summary [TKEStack alarm] {{.startsAt}} {{.alarmPolicyType}} {{.metricDisplayName}} {{.value}} {{.unit}} {{.evaluateType}} {{.evaluateValue}}, 告警策略名:{{.alarmPolicyName}}, 指标名:{{.alertName}}, 集群ID:{{.clusterID}}, 工作负载类型:{{.workloadKind}}, 命名空间:{{.namespace}}, POD名称:{{.podName}}, 节点名称:{{.nodeName}}, 节点类型:{{.nodeRole}}'
            }
          </Text>
        </p>
        <h3 className="h3">
          <Trans>推荐示例</Trans>
        </h3>
        <p>
          <Trans>平台提供了summary字段将上述字段进行了组合，推荐使用。</Trans>
        </p>
        <p>
          <strong>使用示例：</strong>
          {'{{.summary}}'}
        </p>
        <h3 className="h3">
          <Trans>自定义示例</Trans>
        </h3>
        <p>
          <Trans>想要显示 在某个时间，哪个集群/命名空间/负载，触发了哪个告警策略，产生了什么告警内容格式如下</Trans>
        </p>
        <p>
          <strong>body格式：</strong>
          {
            '{{.startsAt}} 集群负载{{.clusterID}}/{{.namespace}}/{{.workloadKind}}/{{.workloadName}}/{{.podName}} 触发了告警策略 {{.alarmPolicyName}} 告警内容 {{.metricDisplayName}} {{.value}} {{.evaluateType}} {{.evaluateValue}}'
          }
        </p>
        <p>
          <Trans>
            <strong>报警内容：</strong>
            {
              '报警内容：2020-04-20T03:16:36Z 集群负载cls-k8kxb7hf/kube-system/Deployment/coredns/coredns-8bf485d54-s2lnc 触发了告警策略 kube-system-deployment 告警内容 CPU利用率 91% > 90%'
            }
          </Trans>
        </p>
      </Card.Body>
    </Card>
  );
};
