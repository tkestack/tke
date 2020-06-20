import * as React from 'react';

import { BaseReactProps } from '@tencent/ff-redux';

interface WorkflowErrorTipProps extends BaseReactProps {
  isShow?: boolean;

  error: AuthObject;
}

interface AuthObject {
  code?: number;

  name?: string;

  isAuthorized?: boolean;

  isLoginedSec?: boolean;

  message?: string;

  redirect?: string;
}

export class WorkflowErrorTip extends React.Component<WorkflowErrorTipProps> {
  render() {
    let { className, error, style, isShow } = this.props;
    let { isAuthorized, isLoginedSec, message, redirect, name } = error;
    return isShow ? (
      <div className={'tc-15-msg ' + className} style={style}>
        {isLoginedSec === undefined ? (
          <span className="tip-info">操作失败，{message}</span>
        ) : isLoginedSec ? (
          isAuthorized === false ? (
            <p>
              对不起，你没有《{name}》模块的访问权限
              {redirect ? (
                <span>
                  ，
                  <a href={redirect} target="_blank">
                    申请权限
                  </a>
                </span>
              ) : (
                <noscript />
              )}
            </p>
          ) : (
            <span className="tip-info">操作失败，{message}</span>
          )
        ) : (
          <p>
            对不起，《{name}》模块的权限状态已失效，请
            {redirect ? (
              <span>
                <a href={redirect} target="_blank">
                  重新登录
                </a>
              </span>
            ) : (
              <noscript />
            )}
          </p>
        )}
      </div>
    ) : (
      <noscript />
    );
  }
}
