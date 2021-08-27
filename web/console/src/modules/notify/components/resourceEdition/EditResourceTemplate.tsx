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
import { t } from '@tencent/tea-app/lib/i18n';
import { Form } from '@tea/component/form';
import { Input } from '@tea/component/input';
import { onChange } from '../../schema/schemaUtil';
import { Select } from '@tencent/tea-component';
import { EditResource } from './EditResource';
import { templateSchema } from '../../schema/templateSchema';
import { router } from '../../router';

export class EditResourceTemplate extends EditResource {
  componentDidMount() {
    this.props.actions.resource.channel.fetch({});
  }
  getSchema() {
    return templateSchema;
  }

  renderForm() {
    let resource = templateSchema;
    resource = this.state.resource;
    const namespaceOptions = this.props.channel.list.data.records.map(c => ({
      value: c.metadata.name,
      text: `${c.spec.displayName}(${c.metadata.name})`
    }));

    if (resource.properties.metadata.properties.namespace.value) {
      const channel = this.props.channel.list.data.records.find(
        c => c.metadata.name === resource.properties.metadata.properties.namespace.value
      );
      let type = 'text';
      if (channel) {
        if (channel.spec.smtp) {
          type = 'text';
        }
        if (channel.spec.tencentCloudSMS) {
          type = 'tencentCloudSMS';
        }
        if (channel.spec.wechat) {
          type = 'wechat';
        }
        if (channel.spec.webhook) {
          type = 'webhook';
        }
      }
      resource.properties.spec.pick = type;
    }

    // 更新模式下disbale渠道
    const { route } = this.props;
    const { mode } = router.resolve(route);

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

        <Form.Item label={t('渠道')} required>
          <Select
            disabled={mode === 'update'}
            size="l"
            placeholder={t('请选择渠道')}
            options={namespaceOptions}
            value={resource.properties.metadata.properties.namespace.value}
            onChange={onChange(resource.properties.metadata.properties.namespace)}
          />
        </Form.Item>
        {this.renderFields(resource.properties.spec.properties[resource.properties.spec['pick']])}
      </Form>
    );
  }
}
