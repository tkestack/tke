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
import { ResourceDetail } from './ResourceDetail';
import { LinkButton } from '../../../common/components';
import { FormPanel } from '@tencent/ff-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Text } from '@tencent/tea-component';
import { router } from '../../router';

export class ResourceDetailTempalte extends ResourceDetail {
  componentDidMount() {
    this.props.actions.resource.channel.fetch({});
  }

  renderMore(ins) {
    return (
      <React.Fragment>
        <FormPanel.Item text label={t('渠道')}>
          {this.renderChannel(ins)}
        </FormPanel.Item>
        {['text', 'tencentCloudSMS', 'wechat', 'webhook']
          .filter(key => ins.spec[key])
          .map(key => {
            return Object.keys(ins.spec[key]).map(property => (
              <FormPanel.Item text key={property} label={property}>
                {ins.spec[key][property] || '-'}
              </FormPanel.Item>
            ));
          })}
      </React.Fragment>
    );
  }

  renderChannel(x) {
    let channelId = x.metadata.namespace;
    let channel = this.props.channel.list.data.records.find(channel => channel.metadata.name === channelId);
    let { route } = this.props;
    let urlParams = router.resolve(route);
    return (
      <React.Fragment>
        <Text>
          <LinkButton
            onClick={() => {
              router.navigate(
                { ...urlParams, mode: 'detail', resourceName: 'channel' },
                { ...route.queries, resourceIns: channelId }
              );
            }}
            className="tea-text-overflow"
          >
            {channelId}
          </LinkButton>
        </Text>
        {channel && <Text theme="weak">({channel.spec.displayName})</Text>}
      </React.Fragment>
    );
  }
}
