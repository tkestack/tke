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

import { FetcherState, FetchState, FFListAction, FFListModel, RecordSet } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { LoadingTip, Table, TableColumn, TableProps, Transfer } from '@tencent/tea-component';
import { StyledProps } from '@tencent/tea-component/lib/_type';
import { removeable } from '@tencent/tea-component/lib/table/addons';

import { TablePanel, TablePanelTagSearchBox, TablePanelTagSearchProps } from '../tablepanel';

interface FormPanelTransferTableTableProps<TResource> extends StyledProps {
  /**列表项配置 */
  columns: TableColumn<TResource>[];
  /**target列表项配置 */
  targetColumns?: TableColumn<TResource>[];
  /**tagSearchOptions */
  tagSearch: TablePanelTagSearchProps<TResource>;
  /**action */
  action: FFListAction;
  /**列表 */
  model: FFListModel<TResource>;

  recordKey?: TableProps['recordKey'];

  /**header */
  header?: React.ReactNode;
  /**title */
  title?: string;
  /**是否禁用的回调 */
  rowDisabled?: (record: TResource) => boolean;
  /**是否需要滚动加载 */
  isNeedScollLoding?: boolean;
  /**sort */
  sortFn?: (a: TResource, b: TResource) => number;
}

interface TransferTableState<TResource> {
  filter?: {
    /**搜索项 */
    keyword?: string;
    /**每次刷新页码大小 */
    pageSize?: number;
  };

  /**列表项 */
  records?: FetcherState<RecordSet<TResource>>[];

  /**左列表是否需要loding */
  isNeedLoading?: boolean;
}
function getFieldValue(record, recordKey) {
  if (typeof recordKey === 'function') {
    return recordKey(record);
  } else {
    return record[recordKey];
  }
}
function FormPanelTransferTable<TResource = any>({ ...props }: FormPanelTransferTableTableProps<TResource>) {
  let [callbackFlag, setCallbackFlag] = React.useState(false);
  let {
    className,
    style,
    columns,
    targetColumns,
    recordKey,
    model: { list, query, selections },
    header,
    rowDisabled = () => false,
    title,
    action,
    sortFn,
    tagSearch,
    isNeedScollLoding = false
  } = props;

  tagSearch.model = props.model;
  tagSearch.action = props.action;

  React.useEffect(() => {
    if (props.model.list.fetchState === FetchState.Ready) {
      setCallbackFlag(false);
    }
  }, [props.model.list.fetchState]);

  let finallist: FetcherState<RecordSet<TResource>>;
  //如果是滚动更新，使用state里面的record否则是ffredux中的list
  if (isNeedScollLoding) {
    let records: FetcherState<RecordSet<TResource>>[] = list.pages;
    let allRecords = [];

    records.forEach(record => {
      if (sortFn) {
        record.data.records = record.data.records.sort(sortFn);
      }
      allRecords = allRecords.concat(record.data.records);
    });

    finallist = Object.assign({}, list, { data: Object.assign({}, list.data, { records: allRecords }) });
  } else {
    finallist = list;
  }

  let scrollableOption = {
    maxHeight: 310,
    minHeight: 310
  };
  if (isNeedScollLoding) {
    scrollableOption['onScrollBottom'] = () => {
      if (!callbackFlag && !finallist.loading && finallist.data.continue) {
        action.next();
        setCallbackFlag(true);
      }
    };
  }

  return (
    <Transfer
      className={className}
      style={style}
      header={header}
      leftCell={
        <Transfer.Cell
          title={
            title +
            t('   共{{totalCount}}项 已加载 {{count}} 项', {
              totalCount: finallist.data.recordCount,
              count: finallist.data.records.length
            })
          }
          scrollable={false}
          tip={t('支持按住shift键进行多选')}
          header={<TablePanelTagSearchBox {...tagSearch} />}
        >
          <TablePanel
            model={{
              list: finallist,
              query: query
            }}
            isNeedCard={false}
            action={action}
            columns={columns}
            recordKey={recordKey}
            rowDisabled={rowDisabled}
            bottomTip={
              isNeedScollLoding && finallist.data.records.length < finallist.data.recordCount ? (
                <LoadingTip />
              ) : (
                  undefined
                )
            }
            scrollable={scrollableOption}
            selectable={{
              value: selections.map(selection => getFieldValue(selection, recordKey)),
              onChange: keys => {
                let listSelection = finallist.data.records.filter(record => {
                  return keys.indexOf(getFieldValue(record, recordKey)) !== -1;
                });
                keys.forEach(key => {
                  if (listSelection.findIndex(item => getFieldValue(item, recordKey) === key) === -1) {
                    let finder = selections.find(item => getFieldValue(item, recordKey) === key);
                    finder && listSelection.push(finder);
                  }
                });
                action.selects(listSelection);
              },
              rowSelect: true
            }}
          />
        </Transfer.Cell>
      }
      rightCell={
        <Transfer.Cell title={t('已选择 {{count}} 项', { count: selections.length })}>
          <Table
            records={selections}
            recordKey={recordKey}
            columns={targetColumns ? targetColumns : columns}
            addons={[
              removeable({
                onRemove: key => {
                  action.selects(
                    selections.filter(record => {
                      return getFieldValue(record, recordKey) !== key;
                    })
                  );
                }
              })
            ]}
          />
        </Transfer.Cell>
      }
    />
  );
}

export { FormPanelTransferTable, FormPanelTransferTableTableProps };
