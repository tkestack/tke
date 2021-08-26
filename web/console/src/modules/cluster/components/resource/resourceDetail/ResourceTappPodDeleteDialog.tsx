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

import { Text } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { WorkflowDialog } from '../../../../common/components';
import { allActions } from '../../../actions';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceTappPodDeleteDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, subRoot, region } = this.props,
      { resourceDetailState } = subRoot,
      { podSelection, removeTappPodFlow } = resourceDetailState;

    return (
      <WorkflowDialog
        caption={t('删除pod')}
        workflow={removeTappPodFlow}
        action={actions.workflow.removeTappPod}
        params={region.selection ? region.selection.value : ''}
        targets={removeTappPodFlow.targets}
        // preAction={() => {
        //   actions.resourceDetail.pod.podSelect([]);
        // }}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <div className="docker-dialog jiqun">
            <p>
              <strong>{t('您确认删除以下pod吗？')}</strong>
            </p>
            <p style={{ maxWidth: '550px' }}>
              {podSelection.map((item, index) => (
                <Text key={index} style={{ marginRight: '10px', wordBreak: 'break-all' }}>
                  {item.metadata.name}
                </Text>
              ))}
            </p>
            <Text theme="danger">{t('pod删除后将不可恢复，请提前备份好数据。')}</Text>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
