import * as React from 'react';
import { MediaObject, Icon, ExternalLink } from '@tencent/tea-component';

interface AuthTipProps {
  name?: string;

  isAuthorized?: boolean;

  message?: string;

  redirect?: string;
}

export class AuthTip extends React.Component<AuthTipProps> {
  render() {
    const { name, redirect } = this.props;
    return (
      <div style={{ position: 'relative' }}>
        <div className="dialog-panel" style={{ position: 'relative', padding: '25px' }}>
          <div className="tc-15-rich-dialog m" style={{ margin: '100px auto' }}>
            <div className="tc-15-rich-dialog-bd">
              <MediaObject media={<Icon type="error" size="l" />}>
                <h3 className="tc-dialog-title">系统提示</h3>
                <p>
                  对不起，你没有《{name}》模块的访问权限
                  {redirect ? (
                    <React.Fragment>
                      ，<ExternalLink href={redirect}>申请权限</ExternalLink>
                    </React.Fragment>
                  ) : (
                    <noscript />
                  )}
                </p>
              </MediaObject>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
