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
