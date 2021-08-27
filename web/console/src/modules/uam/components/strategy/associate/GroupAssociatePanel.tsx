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
import { connect } from 'react-redux';
import { RootProps } from '../StrategyApp';
import { TransferTableProps, TransferTable } from '../../../../common/components';
import { GroupPlain } from '../../../models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface Props extends RootProps{
  onChange?: (selection: GroupPlain[]) => void;
}
@connect(state => state, mapDispatchToProps)
export class GroupAssociatePanel extends React.Component<Props, {}> {

  render() {
    let { groupAssociation, actions, groupPlainList } = this.props;
    // 表示 ResourceSelector 里要显示和选择的数据类型是 `GroupPlain`
    const TransferTableSelector = TransferTable as new () => TransferTable<GroupPlain>;

    // 参数配置
    const selectorProps: TransferTableProps<GroupPlain> = {
      /** 要供选择的数据 */
      model: groupPlainList,

      /** 用于改变model的query值等 */
      action: actions.group.associate.groupList,

      /** 已选中的数据 */
      selections: groupAssociation.groups,

      /** 用户选择发生改变后，应该更新选中的数据状态 */
      onChange: (selection: GroupPlain[]) => {
        actions.group.associate.selectGroup(selection);
        this.props.onChange && this.props.onChange(selection);
      },

      /** 选择器标题 */
      title: t(`当前策略可关联以下用户组`),

      columns: [
        {
          key: 'name',
          header: t('名称'),
          render: (group: GroupPlain) => <p>{`${group.displayName}`}</p>
        }
      ],
      recordKey: 'id'
    };
    return <TransferTableSelector {...selectorProps} />;
  }
}
