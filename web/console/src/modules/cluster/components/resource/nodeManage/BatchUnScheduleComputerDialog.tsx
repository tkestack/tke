import * as React from 'react';
import * as classnames from 'classnames';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { RootProps } from '../../ClusterApp';
import { OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { TipInfo } from '../../../../common/components';
import { getWorkflowError } from '../../../../common/utils';
import { Computer, CreateResource } from '../../../models';
import { allActions } from '../../../actions';
import { resourceConfig } from '../../../../../../config';
import { Button, Modal, Text } from '@tencent/tea-component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });
@connect(state => state, mapDispatchToProps)
export class BatchUnScheduleComputerDialog extends React.Component<RootProps, {}> {
  render() {
    let { subRoot, actions, route, cluster, clusterVersion } = this.props,
      { computer, batchUnScheduleComputer } = subRoot.computerState;
    let action = actions.workflow.batchUnScheduleComputer;
    let workflow = batchUnScheduleComputer;
    let resourceIns = computer.selections[0] && computer.selections[0].metadata.name;
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
          unschedulable: true
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
        resourceIns: computer.selections[0].metadata.name,
        jsonData: JSON.stringify(jsonData)
      };
      action.start([resource], +route.queries['rid']);
      action.perform();
    };

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <Modal
        visible={workflow.operationState !== OperationState.Pending}
        caption={`您确定封锁节点${resourceIns}么？`}
        onClose={cancel}
        size={485}
        disableEscape={true}
      >
        <Modal.Body>
          <Text>封锁节点后，将不接受新的Pod调度到该节点，需要手动取消封锁的节点。</Text>
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
