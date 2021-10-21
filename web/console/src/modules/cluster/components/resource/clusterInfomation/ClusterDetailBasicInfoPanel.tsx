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

import { FormPanel } from '@tencent/ff-component';
import { FetchState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Icon, Switch, Text } from '@tencent/tea-component';

import { dateFormatter } from '../../../../../../helpers';
import { Cluster } from '../../../../common';
import { Clip } from '../../../../common/components';
import { DialogNameEnum } from '../../../models';
import { RootProps } from '../../ClusterApp';
import { ClusterStatus } from '../../clusterManage/ClusterTablePanel';
import { ContainerRuntimeEnum } from '@src/modules/cluster/constants/Config';

export class ClusterDetailBasicInfoPanel extends React.Component<RootProps, {}> {
  render() {
    const { cluster } = this.props;

    return (
      <FormPanel title={t('基本信息')}>
        {cluster.list.fetched !== true || cluster.list.fetchState === FetchState.Fetching ? (
          <Icon type="loading" />
        ) : (
          this._renderClusterInfo()
        )}
      </FormPanel>
    );
  }
  _renderNodeMax() {
    const { cluster } = this.props;
    const clusterInfo: Cluster = cluster.selection;
    if (clusterInfo && clusterInfo.spec.clusterCIDR) {
      const b = clusterInfo.spec.clusterCIDR.split('/')[1];
      const { maxNodePodNum, maxClusterServiceNum } = clusterInfo.spec.properties;
      return Math.pow(2, 32 - parseInt(b)) / maxNodePodNum - Math.ceil(maxClusterServiceNum / maxNodePodNum);
    } else {
      return '';
    }
  }
  /** 处理开关日志采集组件的的操作 */
  private _handleSwitch(cluster: Cluster) {
    const { actions, route } = this.props;
    const enableLogAgent = !cluster.spec.logAgentName;
    if (enableLogAgent) {
      actions.cluster.enableLogAgent(cluster);
    } else {
      actions.cluster.disableLogAgent(cluster);
    }

    actions.cluster.applyFilter({});

    return;
  }
  /** 展示集群的基本信息 */
  private _renderClusterInfo() {
    const { actions, cluster } = this.props;
    const clusterInfo: Cluster = cluster.selection;
    const nodeMax = this._renderNodeMax();
    return cluster.selection ? (
      <React.Fragment>
        <FormPanel.Item label={t('集群名称')} text>
          <Text id="detailClusterName">{clusterInfo.spec.displayName}</Text>
          <Clip target={`#detailClusterName`} />
        </FormPanel.Item>
        <FormPanel.Item label={t('集群ID')} text>
          <Text id="detailClusterId">{clusterInfo.metadata.name}</Text>
          <Clip target={`#detailClusterId`} />
        </FormPanel.Item>
        <FormPanel.Item label={t('状态')} text>
          <Text theme={ClusterStatus[clusterInfo.status.phase]}>{clusterInfo.status.phase || '-'}</Text>
        </FormPanel.Item>
        <FormPanel.Item label={t('Kubernetes版本')} text>
          {clusterInfo.status.version}
        </FormPanel.Item>
        <FormPanel.Item label={t('运行时组件')} text>
          {clusterInfo?.spec?.features?.containerRuntime ?? ContainerRuntimeEnum.DOCKER}
        </FormPanel.Item>
        {clusterInfo.spec.networkDevice && (
          <FormPanel.Item label={t('网卡名称')} text>
            {clusterInfo.spec.networkDevice}
          </FormPanel.Item>
        )}
        {clusterInfo.spec.clusterCIDR && (
          <FormPanel.Item text label={t('容器网络')}>
            <p>{clusterInfo.spec.clusterCIDR}</p>
            {clusterInfo.spec.properties && (
              <p>
                {t('{{ maxClusterServiceNum }}个Service/集群，{{ maxNodePodNum }}个Pod/节点,{{ nodeMax }}个节点/集群', {
                  maxClusterServiceNum: clusterInfo.spec.properties.maxClusterServiceNum,
                  maxNodePodNum: clusterInfo.spec.properties.maxNodePodNum,
                  nodeMax: nodeMax
                })}
              </p>
            )}
          </FormPanel.Item>
        )}
        <FormPanel.Item label={t('集群凭证')} text>
          <Button
            type="link"
            onClick={() => {
              actions.cluster.fetchClustercredential(clusterInfo.metadata.name);
              actions.dialog.updateDialogState(DialogNameEnum.kuberctlDialog);
            }}
          >
            {t('查看集群凭证')}
          </Button>
        </FormPanel.Item>
        <FormPanel.Item label={t('超售比')} text>
          {clusterInfo.spec.properties && clusterInfo.spec.properties.oversoldRatio ? (
            <React.Fragment>
              <Text>
                {clusterInfo.spec.properties.oversoldRatio.cpu
                  ? `CPU:${clusterInfo.spec.properties.oversoldRatio.cpu} `
                  : ''}
              </Text>
              <Text>
                {clusterInfo.spec.properties.oversoldRatio.memory
                  ? `Memory:${clusterInfo.spec.properties.oversoldRatio.memory}`
                  : ''}
              </Text>
            </React.Fragment>
          ) : (
            <Text>{t('暂无设置超售比')}</Text>
          )}
          <Button
            icon="pencil"
            onClick={() => {
              actions.cluster.initClusterAllocationRatio(clusterInfo);
              actions.workflow.updateClusterAllocationRatio.start([]);
            }}
          />
        </FormPanel.Item>
        <FormPanel.Item label={t('创建时间')} text>
          {dateFormatter(new Date(clusterInfo.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}
        </FormPanel.Item>
        {/* <FormPanel.Item label={t('日志采集')} text>
          <Switch
            value={Boolean(clusterInfo.spec.logAgentName)}
            onChange={value => {
              this._handleSwitch(clusterInfo);
            }}
          />
        </FormPanel.Item> */}
      </React.Fragment>
    ) : (
      <noscript />
    );
  }
}
