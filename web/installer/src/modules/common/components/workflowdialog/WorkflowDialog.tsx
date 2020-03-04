import * as React from 'react';
import { Modal, Button } from '@tencent/tea-component';
import { OperationState, WorkflowState, WorkflowActionCreator, isSuccessWorkflow } from '@tencent/ff-redux';
import { BaseReactProps } from '@tencent/qcloud-lib';
import { getWorkflowError } from '../../utils';
import { WorkflowErrorTip } from '../error';

export interface WorkflowDialogProps extends BaseReactProps {
  /**提示框标题 */
  caption?: string;

  /**显示宽度 */
  width?: number;

  /**操作对象 */
  targets?: Array<any>;

  /**操作参数 */
  params?: Object;

  /**操作流 */
  workflow?: WorkflowState<any, any>;

  /**操作 */
  action?: WorkflowActionCreator<any, any>;

  /**是否开启二次确认 */
  confirmMode?: ConfirmMode;

  /**前置操作 */
  preAction?: () => void;

  /**校验操作 */
  validateAction?: () => boolean;

  /**后置操作 */
  postAction?: (obj?: any) => void;
}

interface ConfirmMode {
  /**确认label */
  label: string;

  /**二次确认的值 */
  value: string;
}

interface WorkflowDialogState {
  /**二次确认输入值 */
  confirmValue?: string;

  /**是否通过二次确认 */
  isPassed?: boolean;
}

export class WorkflowDialog extends React.Component<WorkflowDialogProps, WorkflowDialogState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      confirmValue: '',
      isPassed: false
    };
  }

  initState() {
    this.setState({
      confirmValue: '',
      isPassed: false
    });
  }

  handleConfirmValue(value) {
    this.setState({ confirmValue: value, isPassed: value === this.props.confirmMode.value });
  }

  render() {
    let {
        caption,
        width,
        targets,
        params,
        workflow,
        action,
        confirmMode,
        children,
        preAction,
        validateAction,
        postAction
      } = this.props,
      { confirmValue, isPassed } = this.state;

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

      postAction && postAction();
      this.initState();
    };

    const perform = () => {
      preAction && preAction();

      if (!validateAction || (validateAction && validateAction())) {
        action.start(targets, params);
        action.perform();
      }
      this.initState();
    };

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <Modal visible={true} caption={caption || '提示'} onClose={cancel} disableEscape>
        <Modal.Body>{children}</Modal.Body>
        <Modal.Footer>
          <WorkflowErrorTip
            isShow={failed}
            className="error"
            error={getWorkflowError(workflow) || {}}
            style={{ marginTop: '10px' }}
          />
          <Button
            className="m"
            type="primary"
            disabled={workflow.operationState === OperationState.Performing || (confirmMode && !isPassed)}
            onClick={perform}
          >
            {failed ? '重试' : '提交'}
          </Button>
          <Button className="weak m" onClick={cancel}>
            取消
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }
}
