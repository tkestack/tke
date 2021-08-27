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
import * as classnames from 'classnames';
import * as React from 'react';

import { Button, Justify, Text } from '@tencent/tea-component';

import { RootProps } from './InstallerApp';

interface ListItemProps extends RootProps {
  id?: string | number;
}

export class ListItem extends React.Component<ListItemProps> {
  render() {
    const { id, editState, actions } = this.props,
      { machines } = editState;

    let machine = machines.find(m => m.id === id);

    return (
      <div className="run-docker-box" style={{ width: '100%', height: '40px', lineHeight: "40px", marginBottom: '10px', backgroundColor: '#f2f2f2' }}>
        <Justify
          left={<Text style={{ fontSize: '14px', marginLeft: '10px' }}>{machine.host}</Text>}
          right={
            <section>
              <Button
                disabled={false}
                type="link"
                tooltip="编辑"
                onClick={() => actions.installer.updateMachine({ status: 'editing' }, id)}
              >
                <i className="icon-edit-gray" />
              </Button>
              <Button
                disabled={editState.machines.length === 1}
                type="link"
                tooltip={editState.machines.length === 1 ? '不可删除，至少指定一台机器' : '删除'}
                onClick={() => actions.installer.removeMachine(id)}
              >
                <i className="icon-cancel-icon" />
              </Button>
            </section>
          }
        />
      </div>
    );
  }
}
