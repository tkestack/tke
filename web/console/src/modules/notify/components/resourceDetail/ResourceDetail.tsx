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
import { RootProps } from '../NotifyApp';
import { LinkButton } from '../../../common/components';
import { FormPanel } from '@tencent/ff-component';
import { router } from '../../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon } from '@tencent/tea-component';

export class ResourceDetail extends React.Component<RootProps, {}> {
  render() {
    let { actions, route } = this.props;
    let id = route.queries.resourceIns;
    let urlParams = router.resolve(route);
    let resource = this.props[urlParams.resourceName] || this.props.channel;
    let ins = resource.list.data.records.find(ins => ins.metadata.name === id);
    return ins ? this.renderIns(ins) : <Icon type="loading" />;
  }

  renderIns(ins) {
    let { actions, route } = this.props;
    let urlParams = router.resolve(route);
    return (
      <FormPanel
        title={t('基本信息')}
        operation={
          <LinkButton
            onClick={() => {
              router.navigate({ ...urlParams, mode: 'update' }, route.queries);
            }}
          >
            {t('编辑')}
          </LinkButton>
        }
      >
        <FormPanel.Item text label={t('名称')}>
          {ins.spec.displayName}
        </FormPanel.Item>
        {this.renderMore(ins)}
      </FormPanel>
    );
  }

  renderMore(ins) {
    return <React.Fragment />;
  }
}
