import * as React from 'react';

import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { router } from '../../router';
import { CreateApiKeyPanel } from './CreateApiKeyPanel';
import { ApiKeyTablePanel } from './ApiKeyTablePanel';
import { RootProps } from '../RegistryApp';

export class ApiKeyContainer extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);

    if (urlParam['sub'] === 'apikey') {
      if (urlParam['mode'] === 'create') {
        return <CreateApiKeyPanel {...this.props} />;
      } else if (urlParam['mode'] === 'create') {
        return <ApiKeyTablePanel {...this.props} />;
      } else {
        return <ApiKeyTablePanel {...this.props} />;
      }
    } else {
      return null;
    }
  }
}
