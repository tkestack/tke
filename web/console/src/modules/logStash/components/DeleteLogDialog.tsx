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

import { Alert, Button, Modal } from '@tea/component';
import { bindActionCreators, isSuccessWorkflow, OperationState } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { allActions } from '../actions';
import { RootProps } from './LogStashApp';

const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class DeleteLogDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, inlineDeleteLog } = this.props;
    const workflow = inlineDeleteLog;
    const action = actions.workflow.inlineDeleteLog;
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
      action.start(workflow.targets, workflow.params);
      action.perform();
    };

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    if (workflow.operationState === OperationState.Done && isSuccessWorkflow(workflow)) {
      tips.success(t('操作成功'), 1000);
    }

    const targetName = workflow.targets[0].resourceIns;
    return (
      <Modal visible={true} caption={t('删除日志收集规则')} onClose={cancel} disableEscape={true} size={575}>
        <Modal.Body>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                <strong>{t('您确定要删除日志收集规则"{{targetName}}"吗？', { targetName })}</strong>
              </p>
              <div className="block-help-text">
                {t('删除日志收集规则后将不再继续按照规则收集日志，但已收集日志仍存于消费端不受影响。')}
              </div>
            </div>

            {failed && <Alert type="error">{workflow.results[0].error.message}</Alert>}
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" disabled={workflow.operationState === OperationState.Performing} onClick={perform}>
            {workflow.operationState === OperationState.Performing ? <i className="n-loading-icon" /> : ''}
            {failed ? t('重试') : t('确定')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}
