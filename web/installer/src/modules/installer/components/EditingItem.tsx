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
import * as React from 'react';

import { Button, Form, Input, Justify, Segment, Text } from '@tencent/tea-component';

import { getValidateStatus } from '../../common/utils/getValidateStatus';
import { validateActions } from '../actions/validateActions';
import { RootProps } from './InstallerApp';

interface EditingItemProps extends RootProps {
  id?: string | number;
}

const list = [
  { name: '密码认证', value: 'password' },
  { name: '密钥认证', value: 'cert' }
];

export class EditingItem extends React.Component<EditingItemProps> {
  render() {
    const { id, editState, actions } = this.props,
      { machines } = editState;
    const machine = machines.find(m => m.id === id);
    return (
      <div className="run-docker-box" style={{ width: '100%', backgroundColor: '#f2f2f2' }}>
        <Justify
          right={
            <section>
              <Button
                tooltip="保存"
                type="link"
                onClick={() => {
                  const canSave = validateActions._validateMachine(machine);
                  actions.validate.validateMachine(machine);
                  if (canSave) {
                    actions.installer.updateMachine({ status: 'edited' }, id);
                  }
                }}
              >
                <i className="icon-submit-gray" />
              </Button>
              <Button
                disabled={machines.length === 1}
                tooltip={machines.length === 1 ? '不可删除，至少指定一台机器' : '删除'}
                type="link"
                onClick={() => actions.installer.removeMachine(id)}
              >
                <i className="icon-cancel-icon" />
              </Button>
            </section>
          }
        />
        <Form>
          <Form.Item
            label="访问地址"
            required
            status={getValidateStatus(machine.v_host)}
            message={machine.v_host.message}
          >
            <Input
              placeholder="请输入访问地址,eg:ip1;ip2;ip3"
              value={machine.host}
              onChange={value => actions.installer.updateMachine({ host: value }, id)}
            />
            <Text theme="text">注意：要求当前运行安装器所在设备网络可达目标机器</Text>
          </Form.Item>
          <Form.Item
            label="SSH端口"
            required
            status={getValidateStatus(machine.v_port)}
            message={machine.v_port.message}
          >
            <Input
              style={{ width: '80px' }}
              value={machine.port}
              onChange={port => actions.installer.updateMachine({ port }, id)}
            />
          </Form.Item>
          <Form.Item>
            <Segment
              options={[
                { text: '密码认证', value: 'password' },
                { text: '密钥认证', value: 'cert' }
              ]}
              value={machine.authWay}
              onChange={value => actions.installer.updateMachine({ authWay: value }, id)}
            />
          </Form.Item>
          <Form.Item
            label="用户名"
            required
            status={getValidateStatus(machine.v_user)}
            message={machine.v_user.message}
          >
            <Input
              // disabled={true}
              placeholder="请输入特权用户名"
              value={machine.user}
              onChange={user => actions.installer.updateMachine({ user }, id)}
            />
          </Form.Item>
          <Form.Item
            label="密码"
            required
            status={getValidateStatus(machine.v_password)}
            message={machine.v_password.message}
          >
            <Input
              type="password"
              value={machine.password}
              onChange={password => actions.installer.updateMachine({ password }, id)}
            />
          </Form.Item>
          <Form.Item
            label="证书"
            required
            status={getValidateStatus(machine.v_cert)}
            message={machine.v_cert.message}
            style={{
              display: machine.authWay === 'cert' ? 'table-row' : 'none'
            }}
          >
            <Input
              value={machine.cert}
              multiline
              style={{ width: '400px' }}
              onChange={cert => actions.installer.updateMachine({ cert }, id)}
            />
          </Form.Item>
        </Form>
      </div>
    );
  }
}
