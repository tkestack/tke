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

import { isSuccessWorkflow, OperationState } from '@tencent/ff-redux';
import { Alert, Button, Form } from '@tencent/tea-component';

import { getWorkflowError } from '../../common/utils';
import { RootProps } from './InstallerApp';

export class Step9 extends React.Component<RootProps> {
  render() {
    const { actions, editState, step, createCluster } = this.props;
    let failed = createCluster.operationState === OperationState.Done && !isSuccessWorkflow(createCluster);

    return step === 'step9' ? (
      <section>
        <Form.Title>基本配置</Form.Title>
        <Form>
          <Form.Item label="用户名">
            <Form.Text>{editState.username}</Form.Text>
          </Form.Item>
          <Form.Item label="密码">
            <Form.Text>{editState.password}</Form.Text>
          </Form.Item>
        </Form>
        <hr />
        <Form.Title>高可用设置</Form.Title>
        <Form>
          <Form.Item label="高可用类型">
            <Form.Text>
              {editState.haType === 'tke' ? 'TKE提供' : editState.haType === 'thirdParty' ? '使用已有' : '不设置'}
            </Form.Text>
          </Form.Item>
          {editState.haType === 'tke' ? (
            <Form.Item label="VIP地址">
              <Form.Text>{editState.haTkeVip}</Form.Text>
            </Form.Item>
          ) : editState.haType === 'thirdParty' ? (
            <Form.Item label="VIP地址">
              <Form.Text>{editState.haThirdVip}</Form.Text>
            </Form.Item>
          ) : (
            ''
          )}
        </Form>
        <hr />
        <Form.Title>集群设置</Form.Title>
        <Form>
          <Form.Item label="网卡名称">
            <Form.Text>{editState.networkDevice}</Form.Text>
          </Form.Item>
          <Form.Item label="GPU类型">
            <Form.Text>{editState.gpuType === 'none' ? '不使用' : editState.gpuType}</Form.Text>
          </Form.Item>
          <Form.Item label="容器网络">
            <Form.Text>{editState.cidr}</Form.Text>
          </Form.Item>
          <Form.Item label="master节点">
            {editState.machines.map((m, index) => (
              <Form.Text key={index}>{m.host}</Form.Text>
            ))}
          </Form.Item>
        </Form>
        <hr />
        <Form.Title>认证模块设置</Form.Title>
        <Form>
          <Form.Item label="认证方式">
            <Form.Text>{editState.authType}</Form.Text>
          </Form.Item>
          {editState.authType === 'oidc' ? (
            <>
              <Form.Item label="IssuerUrl">
                <Form.Text>{editState.issuerURL}</Form.Text>
              </Form.Item>
              <Form.Item label="ClientID">
                <Form.Text>{editState.clientID}</Form.Text>
              </Form.Item>
              <Form.Item label="CA证书">
                <Form.Text>{editState.caCert}</Form.Text>
              </Form.Item>
            </>
          ) : (
            <noscript />
          )}
        </Form>
        <hr />
        <Form.Title>镜像仓库设置</Form.Title>
        <Form>
          <Form.Item label="镜像仓库类型">
            <Form.Text>{editState.repoType}</Form.Text>
          </Form.Item>
          {editState.repoType === 'tke' ? (
            <>
              <Form.Item label="域名后缀">
                <Form.Text>{editState.repoSuffix}</Form.Text>
              </Form.Item>

              <Form.Item label="是否安装应用商店">
                <Form.Text>{editState.application ? '是' : '否'}</Form.Text>
              </Form.Item>
            </>
          ) : editState.repoType === 'thirdParty' ? (
            <React.Fragment>
              <Form.Item label="仓库地址">
                <Form.Text>{editState.repoAddress}</Form.Text>
              </Form.Item>
              <Form.Item label="命名空间">
                <Form.Text>{editState.repoNamespace}</Form.Text>
              </Form.Item>
              <Form.Item label="用户名">
                <Form.Text>{editState.repoUser}</Form.Text>
              </Form.Item>
              <Form.Item label="密码">
                <Form.Text>{editState.repoPassword}</Form.Text>
              </Form.Item>
            </React.Fragment>
          ) : (
            <noscript />
          )}
        </Form>
        <hr />
        <Form.Title>业务模块设置</Form.Title>
        <Form>
          <Form.Item label="是否开启">
            <Form.Text>{editState.openBusiness ? '是' : '否'}</Form.Text>
          </Form.Item>
        </Form>
        <hr />
        <Form.Title>审计模块设置</Form.Title>
        <Form>
          <Form.Item label="是否开启">
            <Form.Text>{editState.openAudit ? '是' : '否'}</Form.Text>
          </Form.Item>
          {editState.openAudit && (
            <Form.Item label="ES地址">
              <Form.Text>{editState.auditEsUrl}</Form.Text>
            </Form.Item>
          )}

          {(editState.auditEsUsername || editState.auditEsPassword) && (
            <>
              <Form.Item label="用户名">
                <Form.Text>{editState.auditEsUsername}</Form.Text>
              </Form.Item>
              <Form.Item label="密码">
                <Form.Text>{editState.auditEsPassword}</Form.Text>
              </Form.Item>
            </>
          )}
        </Form>
        <hr />
        <Form.Title>监控模块设置</Form.Title>
        <Form>
          <Form.Item label="监控存储类型">
            <Form.Text>{editState.monitorType}</Form.Text>
          </Form.Item>
          {editState.monitorType === 'es' ? (
            <>
              <Form.Item label="ES地址">
                <Form.Text>{editState.esUrl}</Form.Text>
              </Form.Item>
              {(editState.esUsername || editState.esPassword) && (
                <>
                  <Form.Item label="用户名">
                    <Form.Text>{editState.esUsername}</Form.Text>
                  </Form.Item>
                  <Form.Item label="密码">
                    <Form.Text>{editState.esPassword}</Form.Text>
                  </Form.Item>
                </>
              )}
            </>
          ) : editState.monitorType === 'external-inflexdb' ? (
            <>
              <Form.Item label="InfluxDB地址">
                <Form.Text>{editState.influxDBUrl}</Form.Text>
              </Form.Item>
              <Form.Item label="用户名">
                <Form.Text>{editState.influxDBUsername}</Form.Text>
              </Form.Item>
              <Form.Item label="密码">
                <Form.Text>{editState.influxDBPassword}</Form.Text>
              </Form.Item>
            </>
          ) : (
            <noscript />
          )}
        </Form>
        <hr />
        <Form.Title>控制台设置</Form.Title>
        <Form>
          <Form.Item label="是否开启">
            <Form.Text>{editState.openConsole ? '是' : '否'}</Form.Text>
          </Form.Item>
          <Form.Item label="控制台域名">
            <Form.Text>{editState.consoleDomain || '无'}</Form.Text>
          </Form.Item>
          <Form.Item label="证书类型">
            <Form.Text>{editState.certType}</Form.Text>
          </Form.Item>
        </Form>
        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button style={{ marginRight: '10px' }} type="weak" onClick={() => actions.installer.stepNext('step8')}>
            上一步
          </Button>
          <Button
            type="primary"
            onClick={() => {
              actions.installer.createCluster.start([editState]);
              actions.installer.createCluster.perform();
            }}
          >
            安装
          </Button>
          {failed && (
            <Alert
              type="error"
              style={{
                display: 'inline-block',
                marginTop: '10px',
                marginBottom: '0px'
              }}
            >
              {getWorkflowError(createCluster)}
            </Alert>
          )}
        </Form.Action>
      </section>
    ) : (
      <noscript></noscript>
    );
  }
}
