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
import { RootProps } from '../../HelmApp';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import { Card } from '@tencent/tea-component';

export class ValueYamlPanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, detailState, route } = this.props,
      { helm } = detailState;
    if (!helm) {
      return <noscript />;
    }
    return (
      <Card>
        <Card.Body>
          <CodeMirror
            className={'codeMirrorHeight'}
            value={this.props.detailState.helm.valueYaml}
            options={{
              lineNumbers: true,
              mode: 'yaml',
              theme: 'monokai',
              readOnly: true,
              lineWrapping: true, // 自动换行
              styleActiveLine: true // 当前行背景高亮
            }}
          />
        </Card.Body>
      </Card>
    );
  }
}
