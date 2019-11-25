import * as React from 'react';
import * as classnames from 'classnames';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { RootProps } from '../../ClusterApp';
import { OperationState, isSuccessWorkflow } from '@tencent/qcloud-redux-workflow';
import { TipInfo } from '../../../../common/components';
import { getWorkflowError } from '../../../../common/utils';
import { Computer, CreateResource } from '../../../models';
import { allActions } from '../../../actions';
import { resourceConfig } from '../../../../../../config';
import { Modal } from '@tencent/tea-component/lib/modal/ModalMain';
import { Button, Text } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

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
