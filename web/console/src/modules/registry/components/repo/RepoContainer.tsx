/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import * as React from 'react';

import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { router } from '../../router';
import { RootProps } from '../RegistryApp';
import { CreateImagePanel } from './CreateImagePanel';
import { CreateRepoPanel } from './CreateRepoPanel';
import { RepoDetailPanel } from './RepoDetailPanel';
import { RepoTablePanel } from './RepoTablePanel';

export class RepoContainer extends React.Component<RootProps, {}> {
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
