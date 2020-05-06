import * as React from 'react';
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { Layout, Card } from '@tea/component';
import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { AuditPanel } from './AuditPanel';

const { useState, useEffect } = React;
const { Body, Content } = Layout;
const store = configStore();

export class AuditAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <AuditApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

// @connect(state => state, mapDispatchToProps)
// @((router.serve as any)())
// class AuditApp extends React.Component<RootProps, {}> {
//   render() {
//     return (
//       <Layout>
//         <Body>
//           <Content>
//             <Content.Header title="审计记录"></Content.Header>
//             <Content.Body>
//               <AuditPanel />
//             </Content.Body>
//           </Content>
//         </Body>
//       </Layout>
//     );
//   }
// }

const AuditApp = (props: RootProps) => {
  useEffect(() => {
    return () => {
      console.log('AuditApp render unmount.... ');
      (router.serve as any)();
    };
  });
  console.log('AuditApp render ....');
  return (
    <Layout>
      <Body>
        <Content>
          <Content.Header title="审计记录"></Content.Header>
          <Content.Body>
            <AuditPanel />
          </Content.Body>
        </Content>
      </Body>
    </Layout>
  );
};

