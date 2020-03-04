import * as React from 'react';

import { FormPanel } from '@tencent/ff-component';
import { OperationState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import {
    Button, Card, ContentView, Icon, Input, InputNumber, Justify, Segment, Text
} from '@tencent/tea-component';

import { router } from '../../router';
import { RootProps } from '../RegistryApp';

export class CreateApiKeyPanel extends React.Component<RootProps, any> {
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
                <h2>{t('创建访问凭证')}</h2>
              </React.Fragment>
            }
          />
          ;
        </ContentView.Header>
        <ContentView.Body>
          <Card>
            <Card.Body>
              <FormPanel isNeedCard={false}>
                <FormPanel.Item label={t('凭证描述')}>
                  <Input
                    multiline
                    placeholder="请输入凭证的描述信息"
                    value={this.props.apiKeyCreation.description}
                    onChange={value => this.props.actions.apiKey.inputApiKeyDesc(value)}
                  />
                </FormPanel.Item>
                <FormPanel.Item label={t('过期时间')}>
                  <InputNumber
                    value={this.props.apiKeyCreation.expire}
                    min={1}
                    onChange={value => this.props.actions.apiKey.inputApiKeyExpire(+value)}
                  ></InputNumber>
                  <Segment
                    style={{ marginLeft: '8px' }}
                    rimless={true}
                    options={[
                      { value: 'h', text: t('小时') },
                      { value: 'm', text: t('分钟') }
                    ]}
                    value={this.props.apiKeyCreation.unit}
                    onChange={value => this.props.actions.apiKey.selectApiKeyUnit(value)}
                  />
                </FormPanel.Item>
              </FormPanel>
              <FormPanel.Action>
                <Button
                  type="primary"
                  onClick={() => {
                    this.props.actions.apiKey.createApiKey.start([this.props.apiKeyCreation]);
                    this.props.actions.apiKey.createApiKey.perform();
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
    router.navigate(Object.assign({}, urlParams, { sub: 'apikey', mode: 'list' }), {});
    this.props.actions.apiKey.clearEdition();
  }

  private renderError() {
    return (
      this.props.createApiKey.operationState === OperationState.Done &&
      this.props.createApiKey.results[0].error && (
        <p>
          <Text theme="danger" style={{ verticalAlign: 'middle', fontSize: '12px' }}>
            {this.props.createApiKey.results[0].error.message}
          </Text>
        </p>
      )
    );
  }
}
