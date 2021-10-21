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
import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Text } from '@tencent/tea-component';
import { Modal } from '@tencent/tea-component/lib/modal/ModalMain';

import { resourceConfig } from '../../../../../../config';
import { TipInfo } from '../../../../common/components';
import { getWorkflowError } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { Computer, CreateResource } from '../../../models';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });
@connect(state => state, mapDispatchToProps)
export class BatchTurnOnScheduleComputerDialog extends React.Component<RootProps, {}> {
  render() {
    let { subRoot, actions, route, cluster, clusterVersion } = this.props,
      { computer, batchTurnOnSchedulingComputer } = subRoot.computerState;
    let action = actions.workflow.batchTurnOnScheduleComputer;
    let workflow = batchTurnOnSchedulingComputer;
    let resourceIns = computer.selections[0] && computer.selections[0].metadata.name;
    if (workflow.operationState === OperationState.Pending) {
      return <noscript />;
    }

    const cancel = () => {
      if (workflow.operationState === OperationState.Done) {
        action.reset();
      }
      if (workflow.operationState === OperationState.Started) {
        action.cancel();
      }
    };

    const perform = () => {
      let { subRoot, route } = this.props,
        { mode } = subRoot;
      let jsonData = {
        spec: {
          unschedulable: false
        }
      };
      // 去除当中不需要的数据
      jsonData = JSON.parse(JSON.stringify(jsonData));
      let resource: CreateResource = {
        id: uuid(),
        resourceInfo: resourceConfig(clusterVersion)['node'],
        mode,
        namespace: route.queries['np'],
        clusterId: route.queries['clusterId'],
        resourceIns: resourceIns,
        jsonData: JSON.stringify(jsonData)
      };
      action.start([resource], +route.queries['rid']);
      action.perform();
    };

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <Modal
        visible={workflow.operationState !== OperationState.Pending}
        caption={`您确定取消封锁节点${resourceIns}么？`}
        onClose={cancel}
        size={485}
        disableEscape={true}
      >
        <Modal.Body>
          <Trans>
            <Text>节点取消封锁后将允许新的Pod调度到该节点。</Text>
          </Trans>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" disabled={workflow.operationState === OperationState.Performing} onClick={perform}>
            {failed ? '重试' : '确定'}
          </Button>
          <Button onClick={cancel}>取消</Button>
          {failed && <TipInfo type="error">{getWorkflowError(workflow)}</TipInfo>}
        </Modal.Footer>
      </Modal>
    );
  }
}
