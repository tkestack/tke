import * as React from 'react';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { RootState } from '../models';
import { allActions } from '../actions';
import { MainBodyLayout } from '../../common/layouts';
import { configStore } from '../stores/RootStore';
import { router } from '../router';
import { ResetStoreAction } from '../../../../helpers';
import { ProjectHeadPanel } from './ProjectHeadPanel';
import { ProjectActionPanel } from './ProjectActionPanel';
import { ProjectTablePanel } from './ProjectTablePanel';
import { ProjectDetail } from './ProjectDetail';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ContentView } from '@tencent/tea-component';
import { CreateProjectPanel } from './CreateProjectPanel';
import { CreateNamespacePanel } from './CreateNamespacePanel';

const store = configStore();

export class ProjectAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <ProjectApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
@((router.serve as any)())
class ProjectApp extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);

    if (!urlParam['sub']) {
      return (
        <ContentView>
          <ContentView.Header>
            <ProjectHeadPanel title={t('项目管理')} />
          </ContentView.Header>
          <ContentView.Body>
            <ProjectActionPanel />
            <ProjectTablePanel />
          </ContentView.Body>
        </ContentView>
      );
    } else if (urlParam['sub'] === 'detail') {
      return <ProjectDetail {...this.props} />;
    } else if (urlParam['sub'] === 'create') {
      return (
        <ContentView>
          <ContentView.Header>
            <ProjectHeadPanel isNeedBack={true} title={t('新建业务')} />
          </ContentView.Header>
          <ContentView.Body>
            <CreateProjectPanel />
          </ContentView.Body>
        </ContentView>
      );
    } else if (urlParam['sub'] === 'createNS') {
      return (
        <ContentView>
          <ContentView.Header>
            <ProjectHeadPanel isNeedBack={true} title={t('新建Namespace')} />
          </ContentView.Header>
          <ContentView.Body>
            <CreateNamespacePanel />
          </ContentView.Body>
        </ContentView>
      );
    }
  }
}
