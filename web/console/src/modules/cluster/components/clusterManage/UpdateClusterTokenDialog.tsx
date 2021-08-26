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
import { CreateResource, Validation } from 'src/modules/common';

import { Bubble, Button, Modal, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../config';
import { TipInfo } from '../../../common/components';
import { getWorkflowError } from '../../../common/utils';
import { allActions } from '../../actions';
import { RootProps } from '../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

interface UpdateClusterTokenDialogState {
  token?: string;
  v_token?: Validation;
}

@connect(state => state, mapDispatchToProps)
export class UpdateClusterTokenDialog extends React.Component<RootProps, UpdateClusterTokenDialogState> {
  state = {
    token: '',
    v_token: {
      status: 0,
      message: ''
    }
  };
  render() {
    let { actions, subRoot, route, cluster, clustercredential, updateClusterToken } = this.props;
    let action = actions.workflow.updateClusterToken;
    let workflow = updateClusterToken;

    let { token, v_token } = this.state;

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
      actions.cluster.clearClustercredential();
    };

    const perform = () => {
      if (token !== '') {
        let clustercredentialInfo = resourceConfig(this.props.clusterVersion).clustercredential;
        let data = {
          token
        };

        let createClusterData: CreateResource[] = [
          {
            id: uuid(),
            resourceInfo: clustercredentialInfo,
            mode: 'update',
            isStrategic: false,
            resourceIns: this.props.clustercredential.name,
            jsonData: JSON.stringify(data)
          }
        ];
        action.start(createClusterData, {
          clusterId: route.queries['clusterId']
        });
        action.perform();
      }
    };

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <Modal visible={true} caption={t('修改集群凭证')} onClose={cancel} size={500} disableEscape={true}>
        <Modal.Body>
          <FormPanel isNeedCard={false}>
            <FormPanel.Item label="原token" text>
              <FormPanel.InlineText>{clustercredential.token}</FormPanel.InlineText>
            </FormPanel.Item>
            <FormPanel.Item label="新token">
              <Bubble content={v_token.status === 2 ? v_token.message : null}>
                <div className={v_token.status === 2 ? 'is-error' : ''} style={{ display: 'inline-block' }}>
                  <FormPanel.Input
                    value={token}
                    onChange={value => {
                      this.setState({ token: value });
                    }}
                    onBlur={e => {
                      if (e.target.value === '') {
                        this.setState({
                          v_token: {
                            status: 2,
                            message: 'token不能为空'
                          }
                        });
                      }
                    }}
                  />
                </div>
              </Bubble>
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
}
