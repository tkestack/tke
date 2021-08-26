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

import { Bubble, Button, Modal, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { TipInfo } from '../../../../common/components';
import { getWorkflowError, isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validatorActions } from '../../../actions/validatorActions';
import { RootProps } from '../../ClusterApp';

/**
 * 目前支持的 Taint effect 类型：
NoSchedule：新的Pod不调度到该Node上，不影响正在运行的Pod
PreferNoSchedule：soft版的NoSchedule，尽量不调度到该Node上
NoExecute：新的Pod不调度到该Node上，并且删除（evict）已在运行的Pod。Pod可以增加一个时间
 */
const taintEffectOptions = [
  { text: 'NoSchedule', value: 'NoSchedule' },
  {
    text: 'PreferNoSchedule',
    value: 'PreferNoSchedule'
  },
  {
    text: 'NoExecute',
    value: 'NoExecute'
  }
];
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class UpdateNodeTaintDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, subRoot, route } = this.props,
      {
        taintEdition: { taints },
        updateNodeTaint
      } = subRoot.computerState;
    let action = actions.workflow.updateNodeTaint;
    let workflow = updateNodeTaint;

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
      let { taintEdition } = this.props.subRoot.computerState;
      actions.validate.validateAllComputerTaint();
      if (validatorActions._validateAllComputerTaint(taints)) {
        action.start([taintEdition], {
          clusterId: route.queries['clusterId']
        });
        action.perform();
      }
    };
    let canAdd = isEmpty(taints.filter(x => !x.key));

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <Modal visible={true} caption={t('编辑节点Taint')} onClose={cancel} size={800} disableEscape={true}>
        <Modal.Body>
          <FormPanel isNeedCard={false}>
            <FormPanel.Item
              label="Taint"
              tips={t('设置节点的Taint')}
              message={t(
                '长度不超过63个字符，只能包含字母、数字及"-./"，必须以字母或者数字开头结尾，且不能包含"kubernetes"保留字'
              )}
            >
              {this._renderTaintList()}
              <Button
                type="link"
                style={{ display: 'block' }}
                disabled={!canAdd}
                tooltip={!canAdd ? t('请先完成待编辑项') : null}
                onClick={() => {
                  actions.computer.addTaint();
                }}
              >
                {t('新增Taint')}
              </Button>
            </FormPanel.Item>
          </FormPanel>
          {failed && <TipInfo type="error">{getWorkflowError(workflow)}</TipInfo>}
        </Modal.Body>
        <Modal.Body>
          <Modal.Footer>
            <Button type="primary" disabled={workflow.operationState === OperationState.Performing} onClick={perform}>
              {failed ? t('重试') : t('确定')}
            </Button>
            <Button onClick={cancel}>{t('取消')}</Button>
          </Modal.Footer>
        </Modal.Body>
      </Modal>
    );
  }
  private _renderTaintList() {
    let { actions, subRoot, route } = this.props,
      { taints } = subRoot.computerState.taintEdition;
    return taints.map((taint, index) => {
      return (
        <div key={index} style={{ marginBottom: '10px' }}>
          <Bubble
            content={taint.disabled ? t('默认标签不可以编辑') : taint.v_key.status === 2 ? taint.v_key.message : null}
          >
            <div className={taint.v_key.status === 2 ? 'is-error' : ''} style={{ display: 'inline-block' }}>
              <FormPanel.Input
                placeholder={t('taint名称')}
                value={taint.key}
                disabled={taint.disabled}
                onChange={value => actions.computer.updateTaint({ key: value }, taint.id + '')}
                onBlur={value => actions.validate.validateComputerTaintKey(taint.id)}
              />
            </div>
          </Bubble>
          <FormPanel.InlineText parent="span">=</FormPanel.InlineText>
          <Bubble
            content={
              taint.disabled ? t('默认标签不可以编辑') : taint.v_value.status === 2 ? taint.v_value.message : null
            }
          >
            <div
              className={taint.v_value.status === 2 ? 'is-error' : ''}
              style={{ display: 'inline-block', margin: '0 5px' }}
            >
              <FormPanel.Input
                placeholder={t('taint值')}
                value={taint.value}
                size={'s'}
                disabled={taint.disabled}
                onChange={value => actions.computer.updateTaint({ value: value }, taint.id + '')}
                onBlur={value => actions.validate.validateComputerTaintValue(taint.id)}
              />
            </div>
          </Bubble>
          <Bubble content={taint.disabled ? t('默认标签不可以编辑') : null} style={{ display: 'inline-block' }}>
            <FormPanel.Select
              value={taint.effect}
              options={taintEffectOptions}
              disabled={taint.disabled}
              onChange={value => actions.computer.updateTaint({ effect: value }, taint.id + '')}
            />
          </Bubble>
          {!(route.queries['clusterId'] === 'global' && taint.disabled) && (
            <Button onClick={() => actions.computer.deleteTaint(taint.id + '')} icon="close" />
          )}
        </div>
      );
    });
  }
}
