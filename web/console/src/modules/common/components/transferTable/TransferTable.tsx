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

import { SearchBox } from '@tea/component/searchbox';
import { Table, TableColumn } from '@tea/component/table';
import { removeable } from '@tea/component/table/addons/removeable';
import { selectable } from '@tea/component/table/addons/selectable';
import { Transfer } from '@tea/component/transfer';
import { FetcherState, FetchState, FFListAction, FFListModel, RecordSet } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { LoadingTip, StatusTip } from '@tencent/tea-component';
import { autotip, scrollable } from '@tencent/tea-component/lib/table/addons';

function SourceTable({
  dataSource,
  targetKeys,
  onChange,
  disabled,
  columns,
  recordKey,
  bottomTip,
  nextToPage,
  autotipOptions
}) {
  let scrollableOption = {
    maxHeight: 308,
    minHeight: 308
  };
  if (nextToPage) {
    scrollableOption['onScrollBottom'] = nextToPage;
  }

  return (
    <Table
      records={dataSource}
      recordKey={recordKey}
      rowDisabled={disabled}
      columns={columns}
      bottomTip={bottomTip}
      addons={[
        scrollable(scrollableOption),
        selectable({
          value: targetKeys,
          onChange,
          rowSelect: true
        }),
        autotip(autotipOptions)
      ]}
    />
  );
}

function TargetTable({ columns, dataSource, onRemove, recordKey, disabled }) {
  return (
    <Table
      records={dataSource}
      recordKey={recordKey}
      columns={columns}
      addons={[removeable({ onRemove })]}
      rowDisabled={disabled}
    />
  );
}

export interface TransferTableProps<TResource> {
  /**列表项配置 */
  columns: TableColumn<TResource>[];
  /**action */
  action: FFListAction;
  /**列表 */
  model: FFListModel<TResource>;
  /**id */
  recordKey: string;
  /**header */
  header?: React.ReactNode;
  /**title */
  title?: string;
  /**搜索操作 */
  placeholder?: string;
  /**是否禁用的回调 */
  rowDisabled?: (record: TResource) => boolean;
  /**选择列表项的回调 */
  onChange: (selection: TResource[]) => void;
  /**选中项 */
  selections: TResource[];
  /**是否需要滚动加载 */
  isNeedScollLoding?: boolean;
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
export class TransferTable<T> extends React.Component<TransferTableProps<T>, TransferTableState<T>> {
  //防止回调被执行多次
  callbackFlag = false;
  constructor(props: TransferTableProps<T>, context) {
    super(props, context);
    this.state = {
      filter: {
        keyword: '',
        pageSize: 20
      },
      isNeedLoading: false,
      records: []
    };
  }
  appendData(nextProps: TransferTableProps<T>) {
    let {
      recordKey,
      model: {
        query: {
          paging: { pageSize, pageIndex },
          search
        },
        list
      }
    } = nextProps;
    let { records, filter } = this.state;
    if (this.state.filter.keyword !== search || this.state.filter.pageSize !== pageSize) {
      //如果查询条件发生任何改变，先清理数据
      filter = {
        keyword: search,
        pageSize
      };
      records = [];
    }
    //将回调标志变量置为false
    if (list.fetched) {
      this.callbackFlag = false;
      if (list.error) {
        return;
      }
    }
    //如果当前列表没有列表项，直接返回
    if (list.data.records.length === 0) {
      return;
    }
    try {
      for (let i = 0; i < records.length; i++) {
        if (
          records[i].data.records.length &&
          records[i].data.records[0][recordKey] === list.data.records[0][recordKey]
        ) {
          return;
        }
      }
    } catch (error) {}
    if (list.fetched && !records[pageIndex - 1]) {
      records[pageIndex - 1] = list;
      this.setState({ records, filter });
    }
  }

  setIsNeedLoading(value) {
    this.setState({
      isNeedLoading: value
    });
  }
  componentWillMount() {
    if (this.props.isNeedScollLoding) {
      this.appendData(this.props);
    }
  }

  componentWillReceiveProps(nextProps: TransferTableProps<T>) {
    if (nextProps.isNeedScollLoding) {
      this.appendData(nextProps);
    }
  }
  render() {
    let {
      columns,
      model: { list, query },
      recordKey,
      header,
      rowDisabled = () => false,
      onChange,
      selections,
      title,
      action,
      isNeedScollLoding = false
    } = this.props;
    let finallist: FetcherState<RecordSet<T>>;
    //如果是滚动更新，使用state里面的record否则是ffredux中的list
    if (isNeedScollLoding) {
      let records: FetcherState<RecordSet<T>>[] = this.state.records;
      let lastRecord = records && records.length ? records[records.length - 1] : null;
      let allRecords = [];

      records.forEach(record => {
        allRecords = allRecords.concat(record.data.records);
      });

      finallist = {
        fetchState: lastRecord ? lastRecord.fetchState : FetchState.Fetching,
        fetched: lastRecord ? lastRecord.fetched : false,
        data: {
          recordCount: lastRecord ? lastRecord.data.recordCount : 0,
          records: allRecords
        }
      };
    } else {
      finallist = list;
    }

    /**
     * 判断是否需要展示loading态
     * 1. list.fetched 不为true => 适用于刚进来页面，没有数据
     * 2. list.fetchState 为 Fetching => 表示列表正在拉取
     * 3. list.data.recordCount => 结合1、2判断，列表需要轮询，则不需要展示loading
     * 4. query.search => 如果是有搜索关键词，则需要展示loading
     * 5. this.state.isNeedLoading => 结合4，清除搜索条件之后，也需要展示loading
     */
    let isShowLoading: boolean =
      (list.fetched !== true || list.fetchState === FetchState.Fetching) &&
      (list.data.recordCount === 0 || !!query.search || this.state.isNeedLoading);
    return (
      <Transfer
        header={header}
        leftCell={
          <Transfer.Cell
            scrollable={false}
            title={title}
            tip={t('支持按住 shift 键进行多选')}
            header={
              <SearchBox
                value={query.keyword || ''}
                onChange={keyword => {
                  action.changeKeyword((keyword || '').trim());
                }}
                onSearch={keyword => {
                  action.performSearch((keyword || '').trim());
                }}
                onClear={() => {
                  action.changeKeyword('');
                  action.performSearch('');
                }}
              />
            }
          >
            <SourceTable
              dataSource={finallist.data.records}
              targetKeys={selections.map(selection => selection[recordKey])}
              disabled={rowDisabled}
              onChange={keys => {
                let listSelection = finallist.data.records.filter(record => {
                  return keys.indexOf(record[recordKey] as string) !== -1;
                });
                keys.forEach(key => {
                  if (listSelection.findIndex(item => item[recordKey] === key) === -1) {
                    let finder = selections.find(item => item[recordKey] === key);
                    finder && listSelection.push(finder);
                  }
                });
                onChange(listSelection);
              }}
              columns={columns}
              recordKey={recordKey}
              bottomTip={
                isNeedScollLoding && finallist.data.records.length < finallist.data.recordCount ? (
                  <LoadingTip />
                ) : (
                  undefined
                )
              }
              nextToPage={
                isNeedScollLoding
                  ? () => {
                      if (
                        !this.callbackFlag &&
                        !finallist.loading &&
                        finallist.data.records.length < finallist.data.recordCount
                      ) {
                        action.changePaging({
                          pageSize: query.paging.pageSize,
                          pageIndex: query.paging.pageIndex + 1
                        });
                        this.callbackFlag = true;
                      }
                    }
                  : undefined
              }
              autotipOptions={{
                isLoading: isShowLoading,
                isError: list.fetchState === FetchState.Failed,
                isFound: !!query.search,
                onClear: () => {
                  // 清除搜索条件，需要展示loading态，表示正在拉取数据
                  this.setIsNeedLoading(true);
                  // 清楚搜索关键词
                  action.performSearch('');
                  // 这里需要重置一下 isNeedLoading的状态，不然清除搜索条件之后，列表还会展示loading的状态
                  setTimeout(() => {
                    this.setIsNeedLoading(false);
                  }, 500);
                },
                onRetry: () => {
                  action.fetch();
                },
                foundKeyword: query.search,
                emptyText: query.search ? null : (
                  <StatusTip status="empty" emptyText={<div className="text-center">{t('暂无数据')}</div>} />
                )
              }}
            />
          </Transfer.Cell>
        }
        rightCell={
          <Transfer.Cell title={t('已选择 {{count}} 项', { count: selections.length })}>
            <TargetTable
              columns={columns}
              dataSource={selections}
              onRemove={key => {
                onChange(
                  selections.filter(record => {
                    return record[recordKey] !== key;
                  })
                );
              }}
              recordKey={recordKey}
              disabled={rowDisabled}
            />
          </Transfer.Cell>
        }
      />
    );
  }
}
