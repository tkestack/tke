import * as React from 'react';
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { Layout, Card } from '@tea/component';
import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { AlarmRecordPanel } from './AlarmRecordPanel';

const { useState, useEffect } = React;
const { Body, Content } = Layout;
const store = configStore();

export class AlarmRecordContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <AlarmRecordApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
@((router.serve as any)())
class AlarmRecordApp extends React.Component<RootProps, {}> {
  render() {
    return (
      <Layout>
        <Body>
          <Content>
            <Content.Header title="历史告警记录"></Content.Header>
            <Content.Body>
              <AlarmRecordPanel />
            </Content.Body>
          </Content>
        </Body>
      </Layout>
    );
  }
}
