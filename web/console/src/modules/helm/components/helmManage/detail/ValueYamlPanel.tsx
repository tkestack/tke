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
