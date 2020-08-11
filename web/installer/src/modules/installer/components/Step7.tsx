import * as React from 'react';

import { Button, Form, Input, Segment } from '@tencent/tea-component';

import { getValidateStatus } from '../../common/utils';
import { validateActions } from '../actions/validateActions';
import { RootProps } from './InstallerApp';

export class Step7 extends React.Component<RootProps> {
  render() {
    const { actions, editState, step } = this.props;
    return step === 'step7' ? (
      <section>
        <Form>
          <Form.Item label="监控存储类型">
            <Segment
              value={editState.monitorType}
              options={[
                { text: 'TKE提供', value: 'tke-influxdb' },
                { text: '外部InfluxDB', value: 'external-influxdb' },
                { text: '外部ES', value: 'es' },
                { text: '不使用', value: 'none' }
              ]}
              onChange={value => actions.installer.updateEdit({ monitorType: value })}
            />
            <div className="tea-form__help-text">
              {editState.monitorType === 'tke-influxdb'
                ? 'TKE默认将安装InfluxDB作为监控数据存储'
                : editState.monitorType === 'external-influxdb'
                ? '使用您提供的InfluxDB作为监控数据存储，TKE将不再安装监控存储组件'
                : editState.monitorType === 'es'
                ? '使用您提供的Elasticsearch作为监控数据的存储，TKE将不再安装监控存储组件'
                : '不安装监控存储组件，将导致平台不提供监控服务，请谨慎选择'}
            </div>
            {editState.monitorType === 'es' ? (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="ES地址"
                    required
                    status={getValidateStatus(editState.v_esUrl)}
                    message={editState.v_esUrl.message}
                  >
                    <Input
                      value={editState.esUrl}
                      onChange={value => actions.installer.updateEdit({ esUrl: value })}
                      placeholder={'http://10.0.0.1:9200'}
                    />
                  </Form.Item>
                  <Form.Item
                    label="用户名"
                    status={getValidateStatus(editState.v_esUsername)}
                    message={editState.v_esUsername.message}
                  >
                    <Input
                      value={editState.esUsername}
                      onChange={value => actions.installer.updateEdit({ esUsername: value })}
                    />
                  </Form.Item>
                  <Form.Item
                    label="密码"
                    status={getValidateStatus(editState.v_esPassword)}
                    message={editState.v_esPassword.message}
                  >
                    <Input
                      type="password"
                      value={editState.esPassword}
                      onChange={value => actions.installer.updateEdit({ esPassword: value })}
                    />
                  </Form.Item>
                </Form>
              </div>
            ) : editState.monitorType === 'external-influxdb' ? (
              <div className="run-docker-box" style={{ marginTop: '10px', width: '100%' }}>
                <Form>
                  <Form.Item
                    label="InfluxDB地址"
                    required
                    status={getValidateStatus(editState.v_influxDBUrl)}
                    message={editState.v_influxDBUrl.message}
                  >
                    <Input
                      value={editState.influxDBUrl}
                      onChange={value => actions.installer.updateEdit({ influxDBUrl: value })}
                    />
                  </Form.Item>
                  <Form.Item
                    label="用户名"
                    required
                    status={getValidateStatus(editState.v_influxDBUsername)}
                    message={editState.v_influxDBUsername.message}
                  >
                    <Input
                      value={editState.influxDBUsername}
                      onChange={value => actions.installer.updateEdit({ influxDBUsername: value })}
                    />
                  </Form.Item>
                  <Form.Item
                    label="密码"
                    required
                    status={getValidateStatus(editState.v_influxDBPassword)}
                    message={editState.v_influxDBPassword.message}
                  >
                    <Input
                      type="password"
                      value={editState.influxDBPassword}
                      onChange={value => actions.installer.updateEdit({ influxDBPassword: value })}
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
