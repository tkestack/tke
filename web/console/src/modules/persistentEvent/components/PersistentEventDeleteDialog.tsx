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

import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { WorkflowDialog } from '../../common/components';
import { allActions } from '../actions';
import { CreateResource } from '../models';
import { router } from '../router';
import { RootProps } from './PersistentEventApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class PersistentEventDeleteDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, route, peSelection, resourceInfo, deletePeFlow } = this.props;

    let resourceIns: string = peSelection[0] ? peSelection[0].metadata.name : '';
    let clusterId = route.queries['clusterId'];

    // 需要提交的数据
    let resource: CreateResource = {
      id: uuid(),
      resourceInfo,
      mode: 'delete',
      clusterId,
      resourceIns
    };

    return (
      <WorkflowDialog
        caption={t('删除资源')}
        workflow={deletePeFlow}
        action={actions.workflow.deletePeFlow}
        params={route.queries['rid']}
        targets={[resource]}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                <strong>
                  {t('您确定要删除当前集群 {{clusterId}} 的 {{headTitle}} 资源吗？', {
                    clusterId,
                    headTitle: resourceInfo.headTitle
                  })}
                </strong>
              </p>
              <div className="block-help-text text-danger">{t('该资源下所有Pods将一并销毁，请提前备份好数据。')}</div>
            </div>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
