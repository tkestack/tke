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
