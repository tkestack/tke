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
import { FormPanel } from '@tencent/ff-component';
import { t } from '@tencent/tea-app/lib/i18n';
import { Justify, Tooltip, Select } from '@tencent/tea-component';
import * as React from 'react';
import { RootProps } from './AlarmPolicyApp';
export class AlarmPolicyHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;

    actions.projectNamespace.initProjectList();
  }

  render() {
    let { actions, projectList, namespaceList, projectSelection, namespaceSelection, cluster } = this.props;

    let projectListOptions = projectList.map((p, index) => ({
      text: p.displayName,
      value: p.name
    }));

    const namespaceGroups = namespaceList.data.records.reduce((gr, { clusterDisplayName, clusterName }) => {
      const value = `${clusterDisplayName}(${clusterName})`;
      return { ...gr, [clusterName]: <Tooltip title={value}>{value}</Tooltip> };
    }, {});

    let namespaceOptions = namespaceList.data.records.map(item => {
      const text = `${item.clusterDisplayName}-${item.namespace}`;

      return {
        value: item.name,
        text: <Tooltip title={text}>{text}</Tooltip>,
        groupKey: item.clusterName,
        realText: text
      };
    });
    return (
      <Justify
        left={
          <div style={{ lineHeight: '28px' }}>
            <h2 style={{ float: 'left' }}>{t('告警设置')}</h2>
            <FormPanel.InlineText>{t('项目：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={t('业务')}
              options={projectListOptions}
              value={projectSelection}
              onChange={value => {
                actions.projectNamespace.selectProject(value);
              }}
            ></FormPanel.Select>
            <FormPanel.InlineText>{t('namespace：')}</FormPanel.InlineText>
            <Select
              size="m"
              type="simulate"
              searchable
              filter={(inputValue, { realText }: any) => (realText ? realText.includes(inputValue) : true)}
              appearence="button"
              // label={'namespace'}
              groups={namespaceGroups}
              options={namespaceOptions}
              value={namespaceSelection}
              onChange={value => actions.namespace.selectNamespace(value)}
            />
          </div>
        }
      />
    );
  }
}
