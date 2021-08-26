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
import { connect } from 'react-redux';
import { CreateResource } from 'src/modules/cluster/models';

import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../config';
import { WorkflowDialog } from '../../common/components';
import { allActions } from '../actions';
import { RootProps } from './LogStashApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class OpenLogStashDialog extends React.Component<RootProps, any> {
  render() {
    let { actions, route, authorizeOpenLogFlow, clusterSelection, clusterVersion } = this.props;

    let logcollectorResourceInfo = resourceConfig(clusterVersion)['addon_logcollector'];
    let clusterId = route.queries['clusterId'];
    let jsonData = JSON.stringify({
      kind: logcollectorResourceInfo.headTitle,
      apiVersion:
        (logcollectorResourceInfo.group ? logcollectorResourceInfo.group + '/' : '') + logcollectorResourceInfo.version,
      metadata: {
        name: clusterId
      },
      spec: {
        clusterName: clusterId
      }
    });
    let resource: CreateResource = {
      id: uuid(),
      clusterId,
      resourceInfo: logcollectorResourceInfo,
      mode: 'create',
      jsonData
    };
    return (
      <WorkflowDialog
        caption={t('集群日志采集功能')}
        workflow={authorizeOpenLogFlow}
        action={actions.workflow.authorizeOpenLog}
        params={+route.queries['rid']}
        targets={[resource]}
        isDisabledConfirm={clusterSelection[0] ? false : true}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <p>{t('新建日志收集规则，需要先开通日志收集功能。当前您所选的集群尚未开通')}。</p>
          <p>{t('开通日志收集功能：')}</p>
          <ul style={{ marginLeft: '30px' }}>
            <li>{t('1. 将在集群内所有节点（包括后续新增节点）创建日志采集服务。')}</li>
            <li>
              {t('  2. 请为每个节点预留 ')}
              <em className="text-warning">{t('0.3核 250M')}</em> {t('以上可用资源。')}
            </li>
          </ul>
        </div>
      </WorkflowDialog>
    );
  }
}
