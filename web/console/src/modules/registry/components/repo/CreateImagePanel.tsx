import * as React from 'react';

import { OperationState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Card, ContentView, Icon, Input, Justify, Segment, Text } from '@tencent/tea-component';

import { FormPanel } from '@tencent/ff-component';
import { router } from '../../router';
import { RootProps } from '../RegistryApp';

export class CreateImagePanel extends React.Component<RootProps, any> {
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
                <h2>{t('新建镜像')}</h2>
              </React.Fragment>
            }
          />
        </ContentView.Header>
        <ContentView.Body>
          <Card>
            <Card.Body>
              <FormPanel isNeedCard={false}>
                <FormPanel.Item
                  label={t('名称')}
                  validator={this.props.imageCreation.v_name}
                  input={{
                    size: 'm',
                    placeholder: t('请输入镜像名称，不超过 63 个字符'),
                    value: this.props.imageCreation.name,
                    onChange: value => this.props.actions.image.inputImageName(value)
                  }}
                />
                <FormPanel.Item label={t('命名空间')} text>
                  {this.props.route.queries['nsName']}
                </FormPanel.Item>
                <FormPanel.Item label={t('描述')}>
                  <Input
                    multiline
                    placeholder="请输入镜像的描述信息"
                    value={this.props.imageCreation.displayName}
                    onChange={value => this.props.actions.image.inputImageDesc(value)}
                  />
                </FormPanel.Item>
                <FormPanel.Item label={t('权限类型')}>
                  <Segment
                    options={[
                      { value: 'Public', text: t('公有') },
                      { value: 'Private', text: t('私有') }
                    ]}
                    value={this.props.imageCreation.visibility}
                    onChange={value => this.props.actions.image.selectImageVisibility(value)}
                  />
                </FormPanel.Item>
              </FormPanel>
              <FormPanel.Action>
                <Button
                  type="primary"
                  disabled={this.props.imageCreation.v_name.status !== 1}
                  onClick={() => {
                    if (this.props.imageCreation.v_name.status === 1) {
                      this.props.actions.image.createImage.start([
                        Object.assign({}, this.props.imageCreation, {
                          namespace: this.props.route.queries['ns'],
                          namespaceName: this.props.route.queries['nsName']
                        })
                      ]);
                      this.props.actions.image.createImage.perform();
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
    router.navigate(Object.assign({}, urlParams, { sub: 'repo', mode: 'detail' }), this.props.route.queries);
    this.props.actions.image.clearEdition();
  }

  private renderError() {
    return (
      this.props.createImage.operationState === OperationState.Done &&
      this.props.createImage.results[0].error && (
        <p>
          <Text theme="danger" style={{ verticalAlign: 'middle', fontSize: '12px' }}>
            {this.props.createImage.results[0].error.message}
          </Text>
        </p>
      )
    );
  }
}
