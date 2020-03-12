import * as React from 'react';

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ContentView, Icon, Justify, TabPanel, Tabs, Text } from '@tencent/tea-component';

import { router } from '../../router';
import { RootProps } from '../RegistryApp';
import { ChartTablePanel } from './ChartTablePanel';

export class ChartDetailPanel extends React.Component<RootProps, any> {
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
                  {t('Chart包列表')} <span>({this.props.route.queries['cgName'] || '-'})</span>
                </h2>
              </React.Fragment>
            }
          />
          ;
        </ContentView.Header>
        <ContentView.Body>
          {/* <Tabs ceiling animated={false} tabs={[{ id: 'charts', label: t('Chart列表') }]} placement="top">
            <TabPanel id="charts"> */}
          <ChartTablePanel {...this.props} />
          {/* </TabPanel>
          </Tabs> */}
        </ContentView.Body>
      </ContentView>
    );
  }

  private goBack() {
    let urlParams = router.resolve(this.props.route);
    router.navigate(Object.assign({}, urlParams, { sub: 'chart', mode: 'list' }), this.props.route.queries);
  }
}
