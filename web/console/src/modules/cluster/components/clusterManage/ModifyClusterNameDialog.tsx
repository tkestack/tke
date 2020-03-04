import * as React from 'react';
import classNames from 'classnames';
import { OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { Modal, Button } from '@tea/component';
import { CreateResource, Validation, initValidator } from '../../../common/models';
import { RootProps } from '../ClusterApp';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { getWorkflowError, TipInfo } from '../../../../modules/common';
import { resourceConfig } from '../../../../../config';
import { uuid, bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../../actions';
import { FormPanel } from '@tencent/ff-component';

interface ModifyClusterNameState {
  name: string;
  v_name: Validation;
}
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ModifyClusterNameDialog extends React.Component<RootProps, ModifyClusterNameState> {
  constructor(props, content) {
    super(props, content);
    this.state = {
      name: '',
      v_name: initValidator
    };
  }
  _validateName(name) {
    let status = 0,
      message = '';

    //验证集群名称
    if (!name) {
      status = 2;
      message = t('集群名称不能为空');
    } else if (name.length > 60) {
      status = 2;
      message = t('集群名称不能超过60个字符');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  }

  validateName() {
    let { name } = this.state;
    let result = this._validateName(name);
    this.setState({ v_name: result });
  }
  render() {
    let { name, v_name } = this.state;
    let { modifyClusterName, actions, cluster } = this.props;
    let clusterInfo = resourceConfig().cluster;
    const workflow = modifyClusterName;
    const action = actions.workflow.modifyClusterName;
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
      this.validateName();
      if (this._validateName(this.state.name).status !== 2) {
        const createResource: CreateResource[] = [
          {
            id: uuid(),
            resourceInfo: clusterInfo,
            mode: 'update',
            clusterId: cluster.selection && cluster.selection.metadata.name,
            jsonData: JSON.stringify({
              spec: {
                displayName: name
              }
            })
          }
        ];

        action.start(createResource, 1);
        action.perform();
        this.setState({ name: '', v_name: initValidator });
      }
    };

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);
    let oldName = cluster.selection && cluster.selection.spec.displayName;
    return (
      <Modal visible={true} caption={t('编辑集群名称')} onClose={cancel} disableEscape={true}>
        <Modal.Body>
          <FormPanel isNeedCard={false}>
            <FormPanel.Item text label={t('原名称')}>
              {oldName}
            </FormPanel.Item>
            <FormPanel.Item
              label={t('新名称')}
              validator={v_name}
              message="最长60个字符"
              errorTipsStyle={'Icon'}
              input={{
                value: name,
                onBlur: () => {
                  this.validateName();
                },
                onChange: value => {
                  this.setState({ name: value });
                }
              }}
            />
          </FormPanel>
          {failed && <TipInfo className="error">{getWorkflowError(workflow)}</TipInfo>}
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" disabled={workflow.operationState === OperationState.Performing} onClick={perform}>
            {failed ? t('重试') : t('提交')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}
