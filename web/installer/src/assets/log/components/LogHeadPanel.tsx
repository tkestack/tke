import * as React from 'react';
import { RootProps } from '../components/LogApp';

export class LogHeadPanel extends React.Component<RootProps, void> {
  render() {
    return (
      <div className="manage-area-title">
        <h2 style={{ float: 'left' }}>控制台日志</h2>
      </div>
    );
  }
}
