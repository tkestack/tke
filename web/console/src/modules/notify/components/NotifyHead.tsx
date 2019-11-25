import * as React from 'react';
import { RootProps } from './NotifyApp';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component';
export class NotifyHead extends React.Component<RootProps, {}> {
  render() {
    return (
      <Justify
        left={
          <div style={{ lineHeight: '28px' }}>
            <h2 style={{ float: 'left' }}>{t('通知设置')}</h2>
          </div>
        }
      />
    );
  }
}
