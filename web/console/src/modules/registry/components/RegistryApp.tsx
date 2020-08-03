import * as React from 'react';
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { ApiKeyContainer } from './apikey/ApiKeyContainer';
import { RepoContainer } from './repo/RepoContainer';
import { ChartApp } from './chart/ChartApp';
import { ChartGroupApp } from './chartgroup/ChartGroupApp';
import { AppCenter } from './AppCenter';

const store = configStore();

export class RegistryAppContainer extends React.Component<any, any> {
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <RegistryApp />
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
class RegistryApp extends React.Component<RootProps, {}> {
  componentDidMount() {
    this.props.actions.image.fetchDockerRegUrl.fetch();
  }

  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);

    if (urlParam['sub'] === 'apikey') {
      return <ApiKeyContainer {...this.props} />;
    } else if (urlParam['sub'] === 'repo') {
      return <RepoContainer {...this.props} />;
    } else if (urlParam['sub'] === 'chart') {
      return <AppCenter {...this.props} />;
    } else if (urlParam['sub'] === 'chartgroup') {
      // return <ChartGroupApp {...this.props} />;
      return <AppCenter {...this.props} />;
    } else {
      return <RepoContainer {...this.props} />;
    }
  }
}
