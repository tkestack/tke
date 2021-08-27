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

import { Table, TableColumn, Text } from '@tea/component';
import { radioable } from '@tea/component/table/addons/radioable';
import { stylize } from '@tea/component/table/addons/stylize';
import { Identifiable } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { CodeMirrorEditor } from '../../../common/components';

interface Resource extends Identifiable {
  /**名称 */
  name?: string;

  /**内容 */
  content?: string;
}

interface ResourceListState {
  /**当前选中Resource */
  selected?: Resource;
}

export interface ResourceListProps {
  /**列表数据 */
  list: Resource[];
}

export class ResourceList extends React.Component<ResourceListProps, ResourceListState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      selected: props.list[0] || {}
    };
  }

  componentWillReceiveProps(nextProps) {
    if (!this.props.list.length && nextProps.list.length) {
      this.setState({ selected: nextProps.list[0] });
    } else if (this.props.list.length && !nextProps.list.length) {
      this.setState({ selected: null });
    }
  }

  render() {
    let { list } = this.props,
      { selected } = this.state;

    return (
      <div>
        <div className="configuration-box" style={{ overflow: 'hidden' }}>
          <div className="version-wrap">{this.renderResourceTable()}</div>
          <div className="rich-textarea simple-mod">
            <CodeMirrorEditor
              title={t('内容')}
              height={360}
              width={480}
              dHeight={420}
              isShowHeader={true}
              isOpenClip={true}
              isOpenDialogEditor={true}
              theme="monokai"
              value={selected ? selected.content : ''}
              readOnly={true}
              isForceRefresh={true}
            />
          </div>
        </div>
      </div>
    );
  }

  private renderResourceTable() {
    let { list } = this.props,
      { selected } = this.state;

    const colunms: TableColumn<Resource>[] = [
      {
        key: 'resourceName',
        header: t('名称'),
        render: x => {
          return (
            <Text parent="div" overflow>
              {x.name}
            </Text>
          );
        }
      }
    ];

    return (
      <div>
        <Table
          columns={colunms}
          records={list}
          recordKey="id"
          addons={[
            stylize({
              className: 'version-list update-cont'
            }),
            radioable({
              value: selected ? (selected.id as string) : '',
              onChange: key => {
                let record = list.find(record => record.id === key);
                this.setState({ selected: record });
              },
              rowSelect: true,
              width: '5%'
            })
          ]}
        />
      </div>
    );
  }
}
