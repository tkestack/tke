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
import { t } from '@tencent/tea-app/lib/i18n';
import * as React from 'react';
import { TransferTable, TransferTableProps } from '../../common/components';
import { Manager } from '../models';
import { RootProps } from './ProjectApp';

interface EditProjectManagerPanelProps extends RootProps {
  rowDisabled?: (record: Manager) => boolean;
}
export class EditProjectManagerPanel extends React.Component<EditProjectManagerPanelProps, {}> {
  render() {
    let { projectEdition, actions, manager, rowDisabled } = this.props;
    // 表示 ResourceSelector 里要显示和选择的数据类型是 `Manager`
    const TransferTableSelector = TransferTable as new () => TransferTable<Manager>;
    // 参数配置
    const selectorProps: TransferTableProps<Manager> = {
      /** 要供选择的数据 */
      model: manager,

      action: actions.manager,

      rowDisabled: rowDisabled,

      /** 已选中的数据 */
      selections: projectEdition.members,

      /** 用户选择发生改变后，应该更新选中的数据状态 */
      onChange: (selection: Manager[]) => {
        actions.project.selectManager(selection);
      },

      /** 选择器标题 */
      title: t(`当前业务可分配以下责任人`),

      columns: [
        {
          key: 'name',
          header: t('ID/名称'),
          render: (manager: Manager) => <p>{`${manager.displayName}(${manager.name})`}</p>
        }
      ],
      recordKey: 'id'
    };
    return <TransferTableSelector {...selectorProps} />;
  }
}
