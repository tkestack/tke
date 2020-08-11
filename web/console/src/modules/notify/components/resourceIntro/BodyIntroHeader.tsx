import * as React from 'react';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component';
export const BodyIntroHeader = () => {
  return (
    <Justify
      left={
        <div style={{ lineHeight: '28px' }}>
          <h2 style={{ float: 'left' }}>{t('模板说明')}</h2>
        </div>
      }
    />
  );
};
