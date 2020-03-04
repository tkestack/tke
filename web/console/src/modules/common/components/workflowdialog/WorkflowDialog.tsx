import * as React from 'react';
import { Modal, Button } from '@tea/component';
import { OperationState, WorkflowState, WorkflowActionCreator, isSuccessWorkflow } from '@tencent/ff-redux';
import { BaseReactProps, Identifiable } from '@tencent/qcloud-lib';
import { getWorkflowError } from '../../utils';
import { TipInfo } from '../tipinfo';
import { InputField } from '../inputfield';
import { FormItem } from '../formitem';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

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
  postAction?: () => void;

  /** 是否禁用提交的按钮 */
  isDisabledConfirm?: boolean;
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
        postAction,
        isDisabledConfirm = false
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
      <Modal visible={true} caption={caption || t('提示')} onClose={cancel} size={width || 495} disableEscape={true}>
        <Modal.Body>
          {children}
          {confirmMode ? (
            <FormItem label={confirmMode.label}>
              <InputField
                type="text"
                value={confirmValue}
                onChange={this.handleConfirmValue.bind(this)}
                placeholder=""
                tip={t('请输入{{label}}进行确认', {
                  label: confirmMode.label
                })}
              />
            </FormItem>
          ) : (
            <noscript />
          )}

          <TipInfo type="error" isShow={failed}>
            {getWorkflowError(workflow)}
          </TipInfo>
        </Modal.Body>
        <Modal.Footer>
          <Button
            type="primary"
            disabled={
              workflow.operationState === OperationState.Performing || (confirmMode && !isPassed) || isDisabledConfirm
            }
            onClick={perform}
          >
            {failed ? t('重试') : t('确定')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}
