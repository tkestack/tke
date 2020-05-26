import * as React from 'react';

import { Button, Form, Input, Segment, Text, Bubble, Icon } from '@tencent/tea-component';

import { getValidateStatus } from '../../common/utils';
import { validateActions } from '../actions/validateActions';
import { RootProps } from './InstallerApp';

export class Step2 extends React.Component<RootProps> {
  render() {
    const { actions, editState, step } = this.props;
    return step === 'step2' ? (
      <section>
        <Form.Title>
          账户设置
          <Text theme="label" style={{ fontSize: '12px', fontWeight: 400, marginLeft: '10px' }}>
            设置控制台管理员账号信息
          </Text>
        </Form.Title>

        <Form layout="fixed">
          <Form.Item
            label="用户名"
            required
            status={getValidateStatus(editState.v_username)}
            message={editState.v_username.message}
          >
            <Input value={editState.username} onChange={value => actions.installer.updateEdit({ username: value })} />
          </Form.Item>
          <Form.Item
            label="密码"
            required
            status={getValidateStatus(editState.v_password)}
            message={editState.v_password.message}
          >
            <Input
              type="password"
              value={editState.password}
              onChange={value => actions.installer.updateEdit({ password: value })}
            />
          </Form.Item>
          <Form.Item
            label="确认密码"
            required
            status={getValidateStatus(editState.v_confirmPassword)}
            message={editState.v_confirmPassword.message}
          >
            <Input
              type="password"
              value={editState.confirmPassword}
              onChange={value => actions.installer.updateEdit({ confirmPassword: value })}
            />
          </Form.Item>
        </Form>
        <hr />
        <Form.Title>
          高可用设置
          <Text theme="label" style={{ fontSize: '12px', fontWeight: 400, marginLeft: '10px' }}>
            设置控制台、APIServer的高可用VIP
          </Text>
        </Form.Title>
        <Form>
          <Form.Item label="高可用类型">
            <Segment
              value={editState.haType}
              options={[
                { value: 'tke', text: 'TKE提供' },
                { value: 'thirdParty', text: '使用已有' },
                { value: 'none', text: '不设置' }
              ]}
              onChange={value => actions.installer.updateEdit({ haType: value })}
            />
            {editState.haType === 'tke' ? (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="VIP地址"
                    required
                    status={getValidateStatus(editState.v_haTkeVip)}
                    message={editState.v_haTkeVip.message}
                  >
                    <Input
                      value={editState.haTkeVip}
                      onChange={value => actions.installer.updateEdit({ haTkeVip: value })}
                    />
                    <Text theme={'weak'}>
                      {
                        '用户需要提供一个可用的IP地址，保证该IP和各master节点可以正常联通，TKE会为集群部署keepalived并配置该IP为VIP'
                      }
                    </Text>
                  </Form.Item>
                </Form>
              </div>
            ) : editState.haType === 'thirdParty' ? (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="VIP地址"
                    required
                    status={getValidateStatus(editState.v_haThirdVip)}
                    message={editState.v_haThirdVip.message}
                  >
                    <Input
                      value={editState.haThirdVip}
                      onChange={value => actions.installer.updateEdit({ haThirdVip: value })}
                    />
                    <React.Fragment>
                      <Input
                        size={'s'}
                        value={editState.haThirdVipPort}
                        onChange={value => actions.installer.updateEdit({ haThirdVipPort: value })}
                      />
                      <Bubble content={'后端6443端口的映射端口'}>
                        <Icon type="info" />
                      </Bubble>
                    </React.Fragment>
                    <Text theme={'weak'}>
                      <p>
                        在用户自定义VIP情况下，VIP后端需要绑定6443（kube-apiserver端口）端口，同时请确保该VIP有至少两个LB后端（master),
                      </p>
                      <p>由于LB自身路由问题，单LB后端情况下存在集群不可用风险。</p>
                    </Text>
                  </Form.Item>
                </Form>
              </div>
            ) : (
              <noscript />
            )}
          </Form.Item>
        </Form>
        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button style={{ marginRight: '10px' }} type="weak" onClick={() => actions.installer.stepNext('step1')}>
            上一步
          </Button>
          <Button
            type="primary"
            onClick={() => {
              actions.validate.validateStep2(editState);
              if (validateActions._validateStep2(editState)) {
                actions.installer.stepNext('step3');
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
