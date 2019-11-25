import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';

interface FixedFormLayoutProps extends BaseReactProps {
  /** 是否去掉 ul 的 margin-top */
  isRemoveUlMarginTop?: boolean;

  /** style */
  style?: any;
}

export class FixedFormLayout extends React.Component<FixedFormLayoutProps, {}> {
  render() {
    let { isRemoveUlMarginTop, style } = this.props;

    return (
      <div className="run-docker-box" style={style}>
        <div className="edit-param-list">
          <div className="param-box">
            <div className="param-bd">
              <ul className="form-list fixed-layout" style={isRemoveUlMarginTop ? { marginTop: '0' } : {}}>
                {this.props.children}
              </ul>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
