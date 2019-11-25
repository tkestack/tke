import * as React from 'react';
import { ResourceDetail } from './ResourceDetail';
import { FormPanel } from '../../../common/components';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export class ResourceDetailChannel extends ResourceDetail {
  renderMore(ins) {
    return (
      <React.Fragment>
        <FormPanel.Item text label={t('类型')}>
          {this.renderType(ins)}
        </FormPanel.Item>
        {['smtp', 'tencentCloudSMS', 'wechat']
          .filter(key => ins.spec[key])
          .map(key => {
            return Object.keys(ins.spec[key])
              .filter(property => property !== 'password')
              .map(property => (
                <FormPanel.Item text key={property} label={property}>
                  {ins.spec[key][property] || '-'}
                </FormPanel.Item>
              ));
          })}
      </React.Fragment>
    );
  }

  renderType(x) {
    if (x.spec.smtp) {
      return t('邮件');
    }
    if (x.spec.tencentCloudSMS) {
      return t('短信');
    }
    if (x.spec.wechat) {
      return t('微信公众号');
    }

    return '-';
  }
}
