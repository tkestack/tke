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

import { Button, Form, Input, Segment, Text, Bubble, Icon } from '@tencent/tea-component';

import { getValidateStatus } from '../../common/utils';
import { validateActions } from '../actions/validateActions';
import { RootProps } from './InstallerApp';

export class Step2 extends React.Component<RootProps> {
  render() {
    const { actions, editState, step } = this.props;
    return step === 'step2' ? (
      <section>
        <Form.Title>
          账户设置
          <Text theme="label" style={{ fontSize: '12px', fontWeight: 400, marginLeft: '10px' }}>
            设置控制台管理员账号信息
          </Text>
        </Form.Title>

        <Form layout="fixed">
          <Form.Item
            label="用户名"
            required
            status={getValidateStatus(editState.v_username)}
            message={editState.v_username.message}
          >
            <Input value={editState.username} onChange={value => actions.installer.updateEdit({ username: value })} />
          </Form.Item>
          <Form.Item
            label="密码"
            required
            status={getValidateStatus(editState.v_password)}
            message={editState.v_password.message}
          >
            <Input
              type="password"
              value={editState.password}
              onChange={value => actions.installer.updateEdit({ password: value })}
            />
          </Form.Item>
          <Form.Item
            label="确认密码"
            required
            status={getValidateStatus(editState.v_confirmPassword)}
            message={editState.v_confirmPassword.message}
          >
            <Input
              type="password"
              value={editState.confirmPassword}
              onChange={value => actions.installer.updateEdit({ confirmPassword: value })}
            />
          </Form.Item>
        </Form>
        <hr />
        <Form.Title>
          高可用设置
          <Text theme="label" style={{ fontSize: '12px', fontWeight: 400, marginLeft: '10px' }}>
            设置控制台、APIServer的高可用VIP
          </Text>
        </Form.Title>
        <Form>
          <Form.Item label="高可用类型">
            <Segment
              value={editState.haType}
              options={[
                { value: 'none', text: '不设置' },
                { value: 'tke', text: 'TKE提供' },
                { value: 'thirdParty', text: '使用已有' }
              ]}
              onChange={value => actions.installer.updateEdit({ haType: value })}
            />
            {editState.haType === 'tke' ? (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="VIP地址"
                    required
                    status={getValidateStatus(editState.v_haTkeVip)}
                    message={
                      <>
                        {editState.v_haTkeVip.message}{' '}
                        <p>
                          <Text theme={'label'} reset>
                            {'用户提供可用的IP地址，TKE部署Keepalive，配置该IP为Master集群的VIP'}
                          </Text>
                        </p>
                      </>
                    }
                  >
                    <Input
                      value={editState.haTkeVip}
                      onChange={value => actions.installer.updateEdit({ haTkeVip: value })}
                    />
                  </Form.Item>
                </Form>
              </div>
            ) : editState.haType === 'thirdParty' ? (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="VIP地址"
                    required
                    status={getValidateStatus(editState.v_haThirdVip)}
                    message={
                      <>
                        {editState.v_haThirdVip.message}{' '}
                        <p>
                          <Text theme={'label'} reset>
                            <p>
                              VIP绑定Master集群的80（tke控制台）、443（tke控制台）、6443（kube-apiserver端口）、31138（tke-auth-api端口）端口，
                            </p>
                            <p>同时确保该VIP有至少两个LB后端（Master节点），以避免LB单后端不可用风险</p>
                          </Text>
                        </p>
                      </>
                    }
                  >
                    <Input
                      value={editState.haThirdVip}
                      placeholder={'请输入ip地址'}
                      onChange={value => actions.installer.updateEdit({ haThirdVip: value })}
                    />
                    <React.Fragment>
                      <Input
                        disabled
                        size={'s'}
                        placeholder={'请输入端口'}
                        value={editState.haThirdVipPort}
                        onChange={value => actions.installer.updateEdit({ haThirdVipPort: value })}
                      />
                      <Bubble content={'后端6443端口的映射端口'}>
                        <Icon type="info" />
                      </Bubble>
                    </React.Fragment>
                  </Form.Item>
                </Form>
              </div>
            ) : (
              <noscript />
            )}
          </Form.Item>
        </Form>
        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button
            type="primary"
            onClick={() => {
              actions.validate.validateStep2(editState);
              if (validateActions._validateStep2(editState)) {
                actions.installer.stepNext('step3');
              }
            }}
          >
            下一步
          </Button>
        </Form.Action>
      </section>
    ) : (
      <noscript></noscript>
    );
  }
}
