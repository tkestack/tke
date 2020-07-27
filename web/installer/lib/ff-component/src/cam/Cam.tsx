import * as React from 'react';

import { BaseReactProps } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble } from '@tencent/tea-component';

export interface CamProps extends BaseReactProps {
  message?: string | React.ReactNode;
  error?: any;
  position?: 'top' | 'right' | 'bottom' | 'left';
  align?: 'start' | 'end';
  style?: any;
}

const interfaceReg = /\(\w+:\w+\)/g;
const camReg = /\(([\w|\/]*:){5}.*?\)/g;

function transMsg(message: string) {
  let msg = message
    .replace(/^\(.+?\)/, '')
    .trim()
    .replace(interfaceReg, '<span class="text-success">$&</span>')
    .replace(camReg, '<span class="text-warning">$&</span>')
    .replace(/^(.+)$/gm, '<p class="rich-text">$1</p>');
  return (
    <div className="authority-wrap" style={{ textAlign: 'left', lineHeight: '1.2em' }}>
      <p className="authority-inf text-weak">
        <Trans>
          <span style={{ verticalAlign: '-1px' }}>该操作需要授权，请联系您的开发商为您添加权限。</span>
          <a target="_blank" href="//www.qcloud.com/document/product/378/4509">
            查看授权操作指南 <i className="external-link-icon" />
          </a>
        </Trans>
      </p>
      {message && (
        <div>
          <p className="authority-tit text-weak">{t('失败信息描述：')}</p>
          <div className="rich-textarea">
            <div dangerouslySetInnerHTML={{ __html: msg }} className="rich-content" />
          </div>
        </div>
      )}
    </div>
  );
}

export class CamBox extends React.Component<CamProps, {}> {
  render() {
    let { message = '' } = this.props;

    return transMsg(message + '');
  }
}

export class CamTips extends React.Component<CamProps, {}> {
  render() {
    let { message = '', position, align, style } = this.props;
    return (
      <Bubble
        // align={align}
        placement={position || 'bottom'}
        // style={style ? style : { width: '520px' }}
        content={transMsg(message + '') || null}
      >
        <span className="text" style={{ fontSize: '14px' }}>
          {' '}
          ******
        </span>
      </Bubble>
    );
  }
}

export class CamPanel extends React.Component<CamProps, {}> {
  render() {
    let { message = '', position, align } = this.props;
    return (
      <div style={{ lineHeight: '1.92em' }}>
        <p>
          <strong>{t('您没有权限访问此数据')}</strong>
        </p>
        <div>
          <span style={{ display: 'table-cell' }}>{t('该操作需要授权，请联系您的开发商为您添加权限。')}</span>
          <Bubble
            // align={align}
            placement={position || 'bottom'}
            // style={{ width: '520px' }}
            content={transMsg(message + '') || null}
          >
            <a href="javascript:;" style={{ fontSize: '12px' }}>
              {t('了解原因')}
            </a>
          </Bubble>
        </div>
      </div>
    );
  }
}

export function isCamRefused(e) {
  return (
    e &&
    (e.code + '' === '4102' ||
      e.code + '' === '42' ||
      (e.code + '').indexOf('UnauthorizedOperation') !== -1 ||
      (e.code + '').indexOf('CamNoAuth') !== -1)
  );
}
