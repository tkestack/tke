/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import * as React from 'react';
import { ResourceDetail } from './ResourceDetail';
import { FormPanel } from '@tencent/ff-component';
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
