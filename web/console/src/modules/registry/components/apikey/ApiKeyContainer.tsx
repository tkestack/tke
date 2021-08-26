/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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
