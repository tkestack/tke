import * as React from 'react';
import { ResourceDetail } from './ResourceDetail';
import { FormPanel } from '@tencent/ff-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export class ResourceDetailReceiver extends ResourceDetail {
  renderMore(ins) {
    return (
      <React.Fragment>
        <FormPanel.Item text label={t('username')}>
          {ins.spec.username || '-'}
        </FormPanel.Item>

        <FormPanel.Item text label={t('手机号码')}>
          {(ins.spec.identities && ins.spec.identities.mobile) || '-'}
        </FormPanel.Item>

        <FormPanel.Item text label={t('电子邮件')}>
          {(ins.spec.identities && ins.spec.identities.email) || '-'}
        </FormPanel.Item>

        <FormPanel.Item text label={t('微信OpenId')}>
          {(ins.spec.identities && ins.spec.identities.wechat_openid) || '-'}
        </FormPanel.Item>
      </React.Fragment>
    );
  }
}
