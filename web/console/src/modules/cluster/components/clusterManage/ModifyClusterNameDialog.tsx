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
import classNames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Modal } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../config';
import { getWorkflowError, TipInfo } from '../../../../modules/common';
import { CreateResource, initValidator, Validation } from '../../../common/models';
import { allActions } from '../../actions';
import { RootProps } from '../ClusterApp';

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
