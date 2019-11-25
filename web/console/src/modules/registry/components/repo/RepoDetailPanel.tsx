import * as React from 'react';

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ContentView, Icon, Justify, TabPanel, Tabs, Text } from '@tencent/tea-component';

import { router } from '../../router';
import { RootProps } from '../RegistryApp';
import { ImageTablePanel } from './ImageTablePanel';

export class RepoDetailPanel extends React.Component<RootProps, any> {
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
                <h2>
                  {t('命名空间')} <span>({this.props.route.queries['nsName'] || '-'})</span>
                </h2>
              </React.Fragment>
            }
          />
          ;
        </ContentView.Header>
        <ContentView.Body>
          <Tabs ceiling animated={false} tabs={[{ id: 'images', label: t('镜像列表') }]} placement="top">
            <TabPanel id="images">
              <ImageTablePanel {...this.props} />
            </TabPanel>
          </Tabs>
        </ContentView.Body>
      </ContentView>
    );
  }

  private goBack() {
    let urlParams = router.resolve(this.props.route);
    router.navigate(Object.assign({}, urlParams, { sub: 'repo', mode: 'list' }), this.props.route.queries);
  }
}
