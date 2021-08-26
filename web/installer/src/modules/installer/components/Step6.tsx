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

import { Button, Form, Switch, Input, InputNumber, InputAdorment } from '@tencent/tea-component';

import { RootProps } from './InstallerApp';
import { getValidateStatus } from '../../common/utils';
import { validateActions } from '../actions/validateActions';

export class Step6 extends React.Component<RootProps> {
  render() {
    const { actions, editState, step } = this.props;
    return step === 'step6' ? (
      <section>
        <Form>
          <Form.Item label="是否开启业务模块" message="关闭业务模块，平台将不安装业务管理相关的功能，建议默认开启">
            <Switch
              value={editState.openBusiness}
              onChange={value => actions.installer.updateEdit({ openBusiness: value })}
            />
          </Form.Item>
          <Form.Item
            label="是否开启审计功能"
            message="审计模块为平台提供了操作记录,用户可以在平台管理进行查询，需用用户提供ES资源"
          >
            <Switch
              value={editState.openAudit}
              onChange={value => actions.installer.updateEdit({ openAudit: value })}
            />
            {editState.openAudit && (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="ES地址"
                    required
                    status={getValidateStatus(editState.v_auditEsUrl)}
                    message={editState.v_auditEsUrl.message}
                  >
                    <Input
                      value={editState.auditEsUrl}
                      placeholder={'http://10.0.0.1:9200'}
                      onChange={value => actions.installer.updateEdit({ auditEsUrl: value })}
                    />
                  </Form.Item>
                  <Form.Item label="保留数据时间" required align={'middle'}>
                    <InputAdorment after={'天'} appearence={'pure'}>
                      <InputNumber
                        min={1}
                        value={editState.auditEsReserveDays}
                        onChange={value => actions.installer.updateEdit({ auditEsReserveDays: value })}
                      />
                    </InputAdorment>
                  </Form.Item>
                  <Form.Item
                    label="用户名"
                    status={getValidateStatus(editState.v_auditEsUsername)}
                    message={editState.v_auditEsUsername.message}
                  >
                    <Input
                      value={editState.auditEsUsername}
                      onChange={value => actions.installer.updateEdit({ auditEsUsername: value })}
                    />
                  </Form.Item>
                  <Form.Item
                    label="密码"
                    status={getValidateStatus(editState.v_auditEsPassword)}
                    message={editState.v_auditEsPassword.message}
                  >
                    <Input
                      type="password"
                      value={editState.auditEsPassword}
                      onChange={value => actions.installer.updateEdit({ auditEsPassword: value })}
                    />
                  </Form.Item>
                </Form>
              </div>
            )}
          </Form.Item>
        </Form>
        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button style={{ marginRight: '10px' }} type="weak" onClick={() => actions.installer.stepNext('step5')}>
            上一步
          </Button>
          <Button
            type="primary"
            onClick={() => {
              actions.validate.validateStep6(editState);
              if (validateActions._validateStep6(editState)) {
                actions.installer.stepNext('step7');
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
