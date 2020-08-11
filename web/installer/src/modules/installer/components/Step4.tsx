import * as React from 'react';

import { Button, Form, Input, Segment } from '@tencent/tea-component';

import { getValidateStatus } from '../../common/utils';
import { validateActions } from '../actions/validateActions';
import { RootProps } from './InstallerApp';

export class Step4 extends React.Component<RootProps> {
  render() {
    const { actions, editState, step } = this.props;
    return step === 'step4' ? (
      <section>
        <Form>
          <Form.Item label="认证方式">
            <Segment
              value={editState.authType}
              options={[
                { text: 'TKE提供', value: 'tke' },
                { text: 'OIDC认证', value: 'oidc' }
              ]}
              onChange={value => actions.installer.updateEdit({ authType: value })}
            />
            <div className="tea-form__help-text">
              {editState.authType === 'tke' ? '使用TKE提供的用户认证功能' : '接入已有的OIDC认证'}
            </div>
            {editState.authType === 'oidc' ? (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="IssuerUrl"
                    required
                    status={getValidateStatus(editState.v_issuerURL)}
                    message={editState.v_issuerURL.message}
                  >
                    <Input
                      value={editState.issuerURL}
                      onChange={value => actions.installer.updateEdit({ issuerURL: value })}
                    />
                  </Form.Item>
                  <Form.Item
                    label="ClientID"
                    required
                    status={getValidateStatus(editState.v_clientID)}
                    message={editState.v_clientID.message}
                  >
                    <Input
                      value={editState.clientID}
                      onChange={value => actions.installer.updateEdit({ clientID: value })}
                    />
                  </Form.Item>
                  <Form.Item
                    label="CA证书"
                    required
                    status={getValidateStatus(editState.v_caCert)}
                    message={editState.v_caCert.message}
                  >
                    <Input
                      multiline
                      style={{ width: '400px' }}
                      value={editState.caCert}
                      onChange={value => actions.installer.updateEdit({ caCert: value })}
                    />
                  </Form.Item>
                </Form>
              </div>
            ) : (
              <noscript />
            )}
          </Form.Item>
        </Form>
        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button style={{ marginRight: '10px' }} type="weak" onClick={() => actions.installer.stepNext('step3')}>
            上一步
          </Button>
          <Button
            type="primary"
            onClick={() => {
              actions.validate.validateStep4(editState);
              if (validateActions._validateStep4(editState)) {
                actions.installer.stepNext('step5');
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
