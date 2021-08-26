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
import { Form } from '@tea/component/form';
import { Input } from '@tea/component/input';
import { onChange } from '../../schema/schemaUtil';
import { EditResource } from './EditResource';
import { receiverSchema } from '../../schema/receiverSchema';

export class EditResourceReceiver extends EditResource {
  getSchema() {
    return receiverSchema;
  }

  renderForm() {
    let resource = receiverSchema;
    resource = this.state.resource;
    return (
      <Form>
        {Object.keys(resource.properties.spec.properties)
          .filter(key => !resource.properties.spec.properties[key].properties)
          .map(key => (
            <Form.Item key={key} label={resource.properties.spec.properties[key].name || key} required>
              <Input
                placeholder={resource.properties.spec.properties[key].name || key}
                value={resource.properties.spec.properties[key].value}
                onChange={onChange(resource.properties.spec.properties[key])}
              />
            </Form.Item>
          ))}

        {Object.keys(resource.properties.spec.properties.identities.properties)
          .filter(key => !resource.properties.spec.properties.identities.properties[key].properties)
          .map(key => (
            <Form.Item key={key} label={resource.properties.spec.properties.identities.properties[key].name || key}>
              <Input
                placeholder={resource.properties.spec.properties.identities.properties[key].name || key}
                value={resource.properties.spec.properties.identities.properties[key].value}
                onChange={onChange(resource.properties.spec.properties.identities.properties[key])}
              />
            </Form.Item>
          ))}
      </Form>
    );
  }
}
