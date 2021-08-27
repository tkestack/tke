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

import { Button, ExternalLink, Form } from '@tencent/tea-component';

import { RootProps } from './InstallerApp';

export class Step1 extends React.Component<RootProps> {
  render() {
    const { isVerified, actions, step } = this.props;
    return step === 'step1' ? (
      <section>
        <Form>
          <Form.Item>
            <Form.Text>
              感谢您选择TKEStack，TKEStack遵循Apache LICENSE 2.0（许可），您可以在此查看许可原文：
              <ExternalLink href="http://www.apache.org/licenses/LICENSE-2.0">
                http://www.apache.org/licenses/LICENSE-2.0
              </ExternalLink>
              除非适用法律要求或以书面形式同意，否则根据“许可”分发的软件将按“原样”分发，不附带任何明示或暗示的保证或条件。
              请参阅许可，以了解许可下的权限和限制。
            </Form.Text>
          </Form.Item>
        </Form>

        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button
            type="primary"
            onClick={() => {
              actions.installer.stepNext('step2');
            }}
          >
            开始
          </Button>
        </Form.Action>
      </section>
    ) : (
      <noscript></noscript>
    );
  }
}
