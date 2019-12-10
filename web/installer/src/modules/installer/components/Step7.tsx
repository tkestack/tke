import * as React from 'react';
import { RootProps } from './InstallerApp';
import { Button, Input, Form, Switch, Segment } from '@tencent/tea-component';
import { getValidateStatus } from '../../common/utils';
import { validateActions } from '../actions/validateActions';

export class Step7 extends React.Component<RootProps> {
  render() {
    const { actions, editState, step } = this.props;
    return step === 'step7' ? (
      <section>
        <Form>
          <Form.Item label="是否开启" message="建议默认开启，否则平台将不会安装控制台，仅可使用命令行工具或API管理集群">
            <Switch
              value={editState.openConsole}
              onChange={value => actions.installer.updateEdit({ openConsole: value })}
            />
          </Form.Item>
          {editState.openConsole ? (
            [
              <Form.Item label="控制台域名">
                <Input
                  value={editState.consoleDomain}
                  onChange={value => actions.installer.updateEdit({ consoleDomain: value })}
                />
              </Form.Item>,
              <Form.Item label="证书类型">
                <Segment
                  value={editState.certType}
                  options={[
                    { text: '自签名证书', value: 'selfSigned' },
                    { text: '指定服务端证书', value: 'thirdParty' }
                  ]}
                  onChange={value => actions.installer.updateEdit({ certType: value })}
                />
                {editState.certType === 'thirdParty' ? (
                  <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                    <Form>
                      <Form.Item
                        label="证书"
                        required
                        status={getValidateStatus(editState.v_certificate)}
                        message={editState.v_certificate.message}
                      >
                        <Input
                          style={{ width: '400px' }}
                          multiline
                          value={editState.certificate}
                          onChange={value => actions.installer.updateEdit({ certificate: value })}
                        />
                      </Form.Item>
                      <Form.Item
                        label="私钥"
                        required
                        status={getValidateStatus(editState.v_privateKey)}
                        message={editState.v_privateKey.message}
                      >
                        <Input
                          style={{ width: '400px' }}
                          multiline
                          value={editState.privateKey}
                          onChange={value => actions.installer.updateEdit({ privateKey: value })}
                        />
                      </Form.Item>
                    </Form>
                  </div>
                ) : (
                  <noscript />
                )}
              </Form.Item>
            ]
          ) : (
            <noscript />
          )}
        </Form>
        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button style={{ marginRight: '10px' }} type="weak" onClick={() => actions.installer.stepNext('step6')}>
            上一步
          </Button>
          <Button
            type="primary"
            onClick={() => {
              actions.validate.validateStep7(editState);
              if (validateActions._validateStep7(editState)) {
                actions.installer.stepNext('step8');
              }
            }}
          >
            下一步
          </Button>
        </Form.Action>
      </section>
    ) : (
      <noscript></noscript>
    );
  }
}
