import * as React from 'react';

import { RootProps } from './InstallerApp';

export class InstallerHeadPanel extends React.Component<RootProps, void> {
  render() {
    return (
      <div className="qc-header-nav" id="topnav" data-reactid=".1.0">
        <div className="qc-header-inner" data-reactid=".1.0.0">
          <div
            className="qc-header-unit qc-header-logo"
            data-reactid=".1.0.0.0"
          >
            <div className="qc-nav-logo" data-reactid=".1.0.0.0.0">
              <a
                className="qc-logo-inner"
                href="javascript:;"
                title="容器服务"
                data-reactid=".1.0.0.0.0.0"
              >
                <span className="qc-logo-icon" data-reactid=".1.0.0.0.0.0.0">
                  腾讯云
                </span>
              </a>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
