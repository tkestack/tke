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
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { Card, Col, Layout, NavMenu, Row, Stepper } from '@tencent/tea-component';

import { ResetStoreAction } from '../../../../helpers';
import { actions } from '../actions';
import { RootState } from '../models';
import { configStore } from '../stores/RootStore';
import { Step10 } from './Step10';
import { Step2 } from './Step2';
import { Step3 } from './Step3';
import { Step4 } from './Step4';
import { Step5 } from './Step5';
import { Step6 } from './Step6';
import { Step7 } from './Step7';
import { Step8 } from './Step8';
import { Step9 } from './Step9';

const { Header, Content, Body } = Layout;

const store = configStore();

export class InstallerAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <InstallerApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof actions;
}

const mapDispatchToProps = dispatch => Object.assign({}, bindActionCreators({ actions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
class InstallerApp extends React.Component<RootProps> {
  componentDidMount() {
    const { actions } = this.props;
    actions.installer.cluster.fetch();
  }

  render() {
    const { step } = this.props;
    const steps = [
      // { id: 'step1', label: '准备工作' },
      { id: 'step2', label: '基本设置' },
      { id: 'step3', label: '集群设置' },
      { id: 'step4', label: '认证设置' },
      { id: 'step5', label: '镜像仓库设置' },
      { id: 'step6', label: '业务设置' },
      { id: 'step7', label: '监控设置' },
      { id: 'step8', label: '控制台设置' },
      { id: 'step9', label: '配置预览' },
      { id: 'step10', label: '安装' }
    ];

    const stepItem = steps.find(s => s.id === step);
    return (
      <Layout>
        <Header>
          <NavMenu
            left={
              <NavMenu.Item>
                <img src="./static/icon/logo.svg" alt="logo" style={{ height: '32px' }} />
              </NavMenu.Item>
            }
          />
        </Header>
        <Body>
          <Content>
            <Content.Body style={{ overflowY: 'auto' }}>
              <div
                style={{
                  maxWidth: '1200px',
                  minHeight: '600px',
                  margin: '0 auto'
                }}
              >
                <h2 style={{ margin: '40px 0px', fontWeight: 600 }}>TKE Stack</h2>
                <Row>
                  <Col span={4}>
                    <Stepper type="process-vertical" current={step} steps={steps} />
                  </Col>
                  <Col span={20}>
                    <Card style={{ height: '100%', position: 'relative' }} className="affix-target">
                      <Card.Body>
                        <h2>{stepItem.label}</h2>
                        <div
                          style={{
                            padding: '20px 60px 120px',
                            fontSize: '14px',
                            backgroundColor: '#fff'
                          }}
                        >
                          {/* <Step1 {...this.props} /> */}
                          <Step2 {...this.props} />
                          <Step3 {...this.props} />
                          <Step4 {...this.props} />
                          <Step5 {...this.props} />
                          <Step6 {...this.props} />
                          <Step7 {...this.props} />
                          <Step8 {...this.props} />
                          <Step9 {...this.props} />
                          <Step10 {...this.props} />
                        </div>
                      </Card.Body>
                    </Card>
                  </Col>
                </Row>
              </div>
            </Content.Body>
          </Content>
        </Body>
      </Layout>
    );
  }
}
