import * as React from 'react';

import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { router } from '../../router';
import { RootProps } from '../RegistryApp';
import { CreateImagePanel } from './CreateImagePanel';
import { CreateRepoPanel } from './CreateRepoPanel';
import { RepoDetailPanel } from './RepoDetailPanel';
import { RepoTablePanel } from './RepoTablePanel';

export class RepoContainer extends React.Component<RootProps, {}> {
  componentDidMount() {
    this.props.actions.image.fetchDockerRegUrl.fetch();
  }

  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);

    if (urlParam['sub'] === 'repo') {
      if (urlParam['mode'] === 'list') {
        return <RepoTablePanel {...this.props} />;
      } else if (urlParam['mode'] === 'create') {
        return <CreateRepoPanel {...this.props} />;
      } else if (urlParam['mode'] === 'icreate') {
        return <CreateImagePanel {...this.props} />;
      } else if (urlParam['mode'] === 'detail' && urlParam['tab'] === 'images') {
        return <RepoDetailPanel {...this.props} />;
      } else {
        return <RepoTablePanel {...this.props} />;
      }
    } else {
      return <RepoTablePanel {...this.props} />;
    }
  }
}
