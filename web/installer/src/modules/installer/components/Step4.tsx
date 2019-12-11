import * as React from 'react';
import { RootProps } from './InstallerApp';
import { Button, Input, Form, Segment } from '@tencent/tea-component';
import { getValidateStatus } from '../../common/utils';
import { validateActions } from '../actions/validateActions';

export class Step4 extends React.Component<RootProps> {
  render() {
    const { actions, editState, step } = this.props;
    return step === 'step4' ? (
      <section>
        <Form>
          <Form.Item label="镜像仓库类型">
            <Segment
              value={editState.repoType}
              options={[
                { text: 'TKE提供', value: 'tke' },
                { text: '第三方仓库', value: 'thirdParty' }
              ]}
              onChange={value => actions.installer.updateEdit({ repoType: value })}
            />
            <div className="tea-form__help-text">
              {editState.repoType === 'tke'
                ? 'TKE将根据设置的信息安装镜像仓库服务'
                : 'TKE将不会安装镜像仓库服务，使用您提供的镜像仓库作为默认镜像仓库服务'}
            </div>
            {editState.repoType === 'tke' ? (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="域名后缀"
                    required
                    status={getValidateStatus(editState.v_repoSuffix)}
                    message={editState.v_repoSuffix.message}
                  >
                    <Input
                      value={editState.repoSuffix}
                      onChange={value => actions.installer.updateEdit({ repoSuffix: value })}
                    />
                  </Form.Item>
                </Form>
              </div>
            ) : editState.repoType === 'thirdParty' ? (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="仓库地址"
                    required
                    status={getValidateStatus(editState.v_repoAddress)}
                    message={editState.v_repoAddress.message}
                  >
                    <Input
                      value={editState.repoAddress}
                      onChange={value => actions.installer.updateEdit({ repoAddress: value })}
                    />
                  </Form.Item>
                  <Form.Item
                    label="命名空间"
                    required
                    status={getValidateStatus(editState.v_repoNamespace)}
                    message={editState.v_repoNamespace.message}
                  >
                    <Input
                      value={editState.repoNamespace}
                      onChange={value => actions.installer.updateEdit({ repoNamespace: value })}
                    />
                  </Form.Item>
                  <Form.Item
                    label="用户名"
                    required
                    status={getValidateStatus(editState.v_repoUser)}
                    message={editState.v_repoUser.message}
                  >
                    <Input
                      value={editState.repoUser}
                      onChange={value => actions.installer.updateEdit({ repoUser: value })}
                    />
                  </Form.Item>
                  <Form.Item
                    label="密码"
                    required
                    status={getValidateStatus(editState.v_repoPassword)}
                    message={editState.v_repoPassword.message}
                  >
                    <Input
                      type="password"
                      value={editState.repoPassword}
                      onChange={value => actions.installer.updateEdit({ repoPassword: value })}
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
