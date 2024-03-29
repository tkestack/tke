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

import { ExternalLink, Text } from '@tea/component';
import { FetchState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import {
    FormItem, LinkButton, ResourceSelectorGeneric, ResourceSelectorInfoRow, ResourceSelectorProps
} from '../../common/components';
import { Resource } from '../../notify/models';
import { router } from '../../notify/router';
import { RootProps } from './AlarmPolicyApp';

// 简单显示一个字符串，包含到 title 里的组件
function WithTitle({ children }: { children?: string }) {
  return <span title={String(children)}>{children}</span>;
}

export class EditAlarmPolicyReceiverGroup extends React.Component<RootProps> {
  componentWillReceiveProps(nextPorp: RootProps) {
    if (
      (this.props.receiverGroup.list.fetchState === FetchState.Fetching &&
        nextPorp.receiverGroup.list.fetchState === FetchState.Ready) ||
      this.props.receiverGroup.selections !== nextPorp.receiverGroup.selections
    ) {
      let selected = nextPorp.receiverGroup.selections.map(s => {
        let item = nextPorp.receiverGroup.list.data.records.filter(m => m.metadata.name === s.metadata.name);
        return item.length ? item[0] : s;
      });
      this.setState({
        groupSelection: selected
      });
    }
  }

  renderGroupUserInfoName(userInfo) {
    let str = '',
      len = userInfo ? userInfo.length : 0;
    for (let i = 0; i < len && i < 4; ++i) {
      str += userInfo[i].name + ' ';
    }
    return str;
  }

  render() {
    let { actions, receiverGroup, alarmPolicyEdition } = this.props;

    // 表示 ResourceSelector 里要显示和选择的数据类型是 `Group`
    const ResouceSelector = ResourceSelectorGeneric as new () => ResourceSelectorGeneric<Resource>;
    // 参数配置
    const selectorProps: ResourceSelectorProps<Resource> = {
      /** 要供选择的数据 */
      list: receiverGroup.list.data.records,

      /** 已选中的数据 */
      selection: receiverGroup.selections,

      className: 'new-strategy-warrant-group',

      /** 用户选择发生改变后，应该更新选中的数据状态 */
      onSelectionChanged: selected => actions.resourceActions.receiverGroup.selects(selected),

      /** 选择器标题 */
      selectorTitle: t('当前账户下有以下接收组'),

      /** 如何渲染具体一项的名字 */
      itemNameRender: item => {
        return (
          <div>
            <Text parent="div" className="m-width" overflow>
              <LinkButton
                title={item.metadata.name}
                onClick={() => {
                  router.navigate(
                    {
                      mode: 'detail',
                      resourceName: 'receiverGroup'
                    },
                    { resourceIns: item.metadata.name }
                  );
                }}
                className="tea-text-overflow"
              >
                {item.metadata.name}
              </LinkButton>
            </Text>
            <div>
              <WithTitle>{item.spec.displayName}</WithTitle>
            </div>
          </div>
        );
      },

      /** 如何渲染具体一项的附带说明 */
      itemDescriptionRender: item => <WithTitle>{this.renderGroupUserInfoName(item.userInfo)}</WithTitle>
    };

    return (
      <FormItem
        label={t('接收组')}
        tips={
          <ExternalLink
            href={
              window.location.pathname.indexOf('tkestack-project')
                ? '/tkestack-project/notify/create/receiverGroup'
                : '/tkestack/notify/create/receiverGroup'
            }
          >
            {t('新建接收组')}
          </ExternalLink>
        }
      >
        <div className="form-unit unit-group">
          <ResouceSelector {...selectorProps} style={{ overflow: 'auto' }}>
            {receiverGroup.list.loading && <ResourceSelectorInfoRow>{t('正在加载...')}</ResourceSelectorInfoRow>}
            {!receiverGroup.list.loading && receiverGroup.list.data.records.length <= 0 && (
              <ResourceSelectorInfoRow>{t('暂无接收组')}</ResourceSelectorInfoRow>
            )}
          </ResouceSelector>

          <div className="is-error">
            <p className="form-input-help" style={{ fontSize: '12px', marginTop: '5px' }}>
              {alarmPolicyEdition.v_groupSelection.message}
            </p>
          </div>
        </div>
      </FormItem>
    );
  }
}
