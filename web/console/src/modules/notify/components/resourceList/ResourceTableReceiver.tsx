/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { ResourceTable } from './ResourceTable';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { TablePanelColumnProps } from '@tencent/ff-component';
import { Resource } from '../../../common';

export class ResourceTableReceiver extends ResourceTable {
  getColumns(): TablePanelColumnProps<Resource>[] {
    return [
      {
        key: 'username',
        header: t('username'),
        render: x => {
          return x.spec.username || '-';
        }
      },
      {
        key: 'mobile',
        header: t('移动号码'),
        render: x => {
          return (x.spec.identities && x.spec.identities.mobile) || '-';
        }
      },
      {
        key: 'email',
        header: t('电子邮件'),
        render: x => {
          return (x.spec.identities && x.spec.identities.email) || '-';
        }
      },
      {
        key: 'wechat_openid',
        header: t('微信OpenID'),
        render: x => {
          return (x.spec.identities && x.spec.identities.wechat_openid) || '-';
        }
      }
    ];
  }
}
