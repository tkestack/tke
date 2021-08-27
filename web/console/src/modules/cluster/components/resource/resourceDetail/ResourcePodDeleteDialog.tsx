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

import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../../config';
import { WorkflowDialog } from '../../../../common/components';
import { allActions } from '../../../actions';
import { CreateResource } from '../../../models';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourcePodDeleteDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, route, subRoot, region, clusterVersion } = this.props,
      { resourceDetailState } = subRoot,
      { podSelection, deletePodFlow } = resourceDetailState;

    let podResourceInfo = resourceConfig(clusterVersion)['pods'];
    let deleteResourceIns = podSelection[0] ? podSelection[0].metadata.name : '';
    let namespace = route.queries['np'];

    // 需要提交的数据
    let resource: CreateResource = {
      id: uuid(),
      resourceInfo: podResourceInfo,
      namespace,
      clusterId: route.queries['clusterId'],
      resourceIns: deleteResourceIns
    };

    return (
      <WorkflowDialog
        caption={t('销毁实例')}
        workflow={deletePodFlow}
        action={actions.workflow.deletePod}
        params={region.selection ? region.selection.value : ''}
        targets={[resource]}
        preAction={() => {
          actions.resourceDetail.pod.podSelect([]);
        }}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                <strong>{t('您确定要销毁实例{{ deleteResourceIns }}吗？', { deleteResourceIns })}</strong>
              </p>
              <div className="block-help-text">{t('实例销毁重建后将不可恢复，请提前备份好数据。')}</div>
            </div>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
