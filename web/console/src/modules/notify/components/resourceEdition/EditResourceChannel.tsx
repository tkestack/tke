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
import * as React from 'react';
import { t } from '@tencent/tea-app/lib/i18n';
import { Form } from '@tea/component/form';
import { Input } from '@tea/component/input';
import { channelSchema } from '../../schema/channelSchema';
import { onChange } from '../../schema/schemaUtil';
import { Radio } from '@tencent/tea-component';
import { EditResource } from './EditResource';
import { PermissionProvider } from '@common';

export class EditResourceChannel extends EditResource {
  getSchema() {
    return channelSchema;
  }

  renderForm() {
    let resource = channelSchema;
    resource = this.state.resource;
    return (
      <Form>
        <Form.Item label={t('名称')} required>
          <Input
            size="l"
            placeholder={t('请填写名称')}
            value={resource.properties.spec.properties.displayName.value}
            onChange={onChange(resource.properties.spec.properties.displayName)}
          />
        </Form.Item>

        <Form.Item label={t('渠道')}>
          <Radio.Group value={resource.properties.spec.pick} onChange={onChange(resource.properties.spec)}>
            <Radio name="smtp">{t('邮件')}</Radio>
            <PermissionProvider value="platform.notify.sms_wechat_webhook">
              <Radio name="tencentCloudSMS">{t('短信')}</Radio>
              <Radio name="wechat">{t('微信公众号')}</Radio>
              <Radio name="webhook">{t('webhook')}</Radio>
            </PermissionProvider>
          </Radio.Group>
        </Form.Item>
        {this.renderFields(resource.properties.spec.properties[resource.properties.spec['pick']])}
      </Form>
    );
  }
}
