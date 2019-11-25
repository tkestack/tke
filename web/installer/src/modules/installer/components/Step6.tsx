import * as React from 'react';
import { RootProps } from './InstallerApp';
import { Button, Form, Switch } from '@tencent/tea-component';

export class Step6 extends React.Component<RootProps> {
  render() {
    const { actions, editState, step } = this.props;
    return step === 'step6' ? (
      <section>
        <Form>
          <Form.Item label="是否开启" message="关闭业务模块，平台将不安装业务管理相关的功能，建议默认开启">
            <Switch
              value={editState.openBusiness}
              onChange={value => actions.installer.updateEdit({ openBusiness: value })}
            />
          </Form.Item>
        </Form>
        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button style={{ marginRight: '10px' }} type="weak" onClick={() => actions.installer.stepNext('step5')}>
            上一步
          </Button>
          <Button type="primary" onClick={() => actions.installer.stepNext('step7')}>
            下一步
          </Button>
        </Form.Action>
      </section>
    ) : (
      <noscript></noscript>
    );
  }
}
