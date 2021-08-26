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

import { FormPanelSelect, FormPanelSelectProps } from '@tencent/ff-component';
import { t } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component';
import * as React from 'react';
import { RootProps } from './AlarmPolicyApp';
export class AlarmPolicyHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    // actions.region.fetch();
    actions.cluster.applyFilter({ regionId: 1 });
  }

  render() {
    let { actions, regionList, regionSelection, cluster } = this.props;

    let selectProps: FormPanelSelectProps = {
      type: 'simulate',
      appearence: 'button',
      label: '集群',
      model: cluster,
      action: actions.cluster,
      valueField: record => record.metadata.name,
      displayField: record => `${record.metadata.name} (${record.spec.displayName})`,
      onChange: (clusterId: string) => {
        actions.cluster.selectCluster(cluster.list.data.records.find(c => c.metadata.name === clusterId));
      }
    };
    return (
      <Justify
        left={
          <div style={{ lineHeight: '28px' }}>
            <h2 style={{ float: 'left' }}>{t('告警设置')}</h2>
            <div className="tc-15-dropdown" style={{ marginLeft: '20px', display: 'inline-block', minWidth: '30px' }}>
              {t('集群')}
            </div>
            <FormPanelSelect {...selectProps} />
          </div>
        }
      />
    );
  }
}
