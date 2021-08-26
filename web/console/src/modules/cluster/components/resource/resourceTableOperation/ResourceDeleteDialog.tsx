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

import { WorkflowDialog } from '../../../../common/components';
import { allActions } from '../../../actions';
import { CreateResource } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';

const deleteTipsMap = {
  resource: t('该Workload下所有Pods将一并销毁，销毁后不可恢复，请谨慎操作。'),
  namespace: t('删除Namespace将销毁Namespace下的所有资源，销毁后不可恢复，请谨慎操作。'),
  service: {
    svc: t('该Service下的负载均衡将一并销毁，销毁后不可恢复，请谨慎操作。'),
    ingress: t('该Ingress下的所有规则将一并删除，销毁后不可恢复，请谨慎操作。')
  }
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceDeleteDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, route, subRoot, namespaceSelection, region } = this.props,
      urlParams = router.resolve(route),
      { deleteResourceFlow, resourceOption, resourceInfo, mode } = subRoot,
      { resourceDeleteSelection } = resourceOption;

    let resourceIns = resourceDeleteSelection[0] ? resourceDeleteSelection[0].metadata.name : '';

    // 需要提交的数据
    let resource: CreateResource = {
      id: uuid(),
      resourceInfo,
      namespace: namespaceSelection,
      clusterId: route.queries['clusterId'],
      resourceIns
    };

    let deleteTips: string | object = deleteTipsMap[urlParams['type']];

    // 这里主要是考虑在更新实例数量的时候，会调用删除接口删除hpa，不应该展示出dialog
    return mode === 'update' ? (
      <noscript />
    ) : (
      <WorkflowDialog
        caption={t('删除资源')}
        workflow={deleteResourceFlow}
        action={actions.workflow.deleteResource}
        params={region.selection ? region.selection.value : ''}
        targets={[resource]}
        isDisabledConfirm={resourceIns ? false : true}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <p style={{ wordWrap: 'break-word' }}>
            <strong>
              {t('您确定要删除{{headTitle}}：{{resourceIns}}吗？', {
                headTitle: resourceInfo.headTitle,
                resourceIns
              })}
            </strong>
          </p>
          {deleteTips && (
            <div className="block-help-text text-danger">
              {typeof deleteTips === 'string' ? deleteTips : deleteTips[urlParams['resourceName']]}
            </div>
          )}
        </div>
      </WorkflowDialog>
    );
  }

  componentDidUpdate() {
    setTimeout(() => {
      let { subRoot, actions } = this.props,
        { deleteResourceFlow } = subRoot;

      if (deleteResourceFlow.operationState === OperationState.Done && isSuccessWorkflow(deleteResourceFlow)) {
        actions.workflow.deleteResource.reset();
      }
    });
  }
}
