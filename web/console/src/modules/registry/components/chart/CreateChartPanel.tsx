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

import { OperationState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Card, ContentView, Icon, Input, Justify, Segment, Text } from '@tencent/tea-component';

import { FormPanel } from '@tencent/ff-component';
import { router } from '../../router';
import { RootProps } from '../RegistryApp';

export class CreateChartPanel extends React.Component<RootProps, any> {
  render() {
    return (
      <ContentView>
        <ContentView.Header>
          <Justify
            left={
              <React.Fragment>
                <a
                  href="javascript:;"
                  className="back-link"
                  onClick={() => {
                    this.goBack();
                  }}
                >
                  <Icon type="btnback" />
                  {t('返回')}
                </a>
                <h2>{t('新建 ChartGroup')}</h2>
              </React.Fragment>
            }
          />
          ;
        </ContentView.Header>
        <ContentView.Body>
          <Card>
            <Card.Body>
              <FormPanel isNeedCard={false}>
                <FormPanel.Item
                  label={t('名称')}
                  validator={this.props.chartCreation.v_name}
                  input={{
                    size: 'm',
                    placeholder: t('请输入 ChartGroup，不超过 63 个字符'),
                    value: this.props.chartCreation.name,
                    onChange: value => this.props.actions.charts.inputChartName(value)
                  }}
                />
                <FormPanel.Item label={t('描述')}>
                  <Input
                    multiline
                    placeholder="请输入该 ChartGroup 的描述信息"
                    value={this.props.chartCreation.displayName}
                    onChange={value => this.props.actions.charts.inputChartDesc(value)}
                  />
                </FormPanel.Item>
                <FormPanel.Item label={t('权限类型')}>
                  <Segment
                    options={[
                      { value: 'Public', text: t('公有') },
                      { value: 'Private', text: t('私有') }
                    ]}
                    value={this.props.chartCreation.visibility}
                    onChange={value => this.props.actions.charts.selectChartVisibility(value)}
                  />
                </FormPanel.Item>
              </FormPanel>
              <FormPanel.Action>
                <Button
                  type="primary"
                  disabled={this.props.chartCreation.v_name.status !== 1}
                  onClick={() => {
                    if (this.props.chartCreation.v_name.status === 1) {
                      this.props.actions.charts.createChart.start([this.props.chartCreation]);
                      this.props.actions.charts.createChart.perform();
                    }
                  }}
                >
                  <Trans>确认</Trans>
                </Button>
                <Button
                  onClick={() => {
                    this.goBack();
                  }}
                >
                  <Trans>取消</Trans>
                </Button>
                {this.renderError()}
              </FormPanel.Action>
            </Card.Body>
          </Card>
        </ContentView.Body>
      </ContentView>
    );
  }

  private goBack() {
    let urlParams = router.resolve(this.props.route);
    router.navigate(Object.assign({}, urlParams, { sub: 'chart', mode: 'list' }), {});
    this.props.actions.charts.clearEdition();
  }

  private renderError() {
    return (
      this.props.createChart.operationState === OperationState.Done &&
      this.props.createChart.results[0].error && (
        <p>
          <Text theme="danger" style={{ verticalAlign: 'middle', fontSize: '12px' }}>
            {this.props.createChart.results[0].error.message}
          </Text>
        </p>
      )
    );
  }
}
