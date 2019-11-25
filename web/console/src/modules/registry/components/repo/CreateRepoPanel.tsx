import * as React from 'react';

import { OperationState } from '@tencent/qcloud-redux-workflow';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Card, ContentView, Icon, Input, Justify, Segment, Text } from '@tencent/tea-component';

import { FormPanel } from '../../../common/components';
import { router } from '../../router';
import { RootProps } from '../RegistryApp';

export class CreateRepoPanel extends React.Component<RootProps, any> {
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
                <h2>{t('新建命名空间')}</h2>
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
                  validator={this.props.repoCreation.v_name}
                  input={{
                    size: 'm',
                    placeholder: t('请输入命名空间，不超过 63 个字符'),
                    value: this.props.repoCreation.name,
                    onChange: value => this.props.actions.repo.inputRepoName(value)
                  }}
                />
                <FormPanel.Item label={t('描述')}>
                  <Input
                    multiline
                    placeholder="请输入该命名空间的描述信息"
                    value={this.props.repoCreation.displayName}
                    onChange={value => this.props.actions.repo.inputRepoDesc(value)}
                  />
                </FormPanel.Item>
                <FormPanel.Item label={t('权限类型')}>
                  <Segment
                    options={[
                      { value: 'Public', text: t('公有') },
                      { value: 'Private', text: t('私有') }
                    ]}
                    value={this.props.repoCreation.visibility}
                    onChange={value => this.props.actions.repo.selectRepoVisibility(value)}
                  />
                </FormPanel.Item>
              </FormPanel>
              <FormPanel.Action>
                <Button
                  type="primary"
                  disabled={this.props.repoCreation.v_name.status !== 1}
                  onClick={() => {
                    if (this.props.repoCreation.v_name.status === 1) {
                      this.props.actions.repo.createRepo.start([this.props.repoCreation]);
                      this.props.actions.repo.createRepo.perform();
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
    router.navigate(Object.assign({}, urlParams, { sub: 'repo', mode: 'list' }), {});
    this.props.actions.repo.clearEdition();
  }

  private renderError() {
    return (
      this.props.createRepo.operationState === OperationState.Done &&
      this.props.createRepo.results[0].error && (
        <p>
          <Text theme="danger" style={{ verticalAlign: 'middle', fontSize: '12px' }}>
            {this.props.createRepo.results[0].error.message}
          </Text>
        </p>
      )
    );
  }
}
