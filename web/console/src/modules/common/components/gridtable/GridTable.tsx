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

import { Card, Pagination, Table, TableProps } from '@tea/component';
import { FetchState, FFListAction, FFListModel } from '@tencent/ff-redux';
import { autotip } from 'tea-component/es/table/addons/autotip';
import { StatusTip } from '@tencent/tea-component/lib/tips';

interface GridTableProps extends TableProps {
  /** 列表的相关配置，包含list、query等 */
  listModel: FFListModel;

  /** fetcher、query相关的action，包含select、selects、clear等 */
  actionOptions: FFListAction;

  /** 空列表的相关提示，不传默认为 暂无数据 */
  emptyTips?: React.ReactNode;

  /** 是否需要展示翻页 */
  isNeedPagination?: boolean;

  /** 表格是否需要展示边框，默认为false */
  isNeedTableBorder?: boolean;

  /** 是否需要 Card */
  isNeedCard?: boolean;

  cardStyle?: React.CSSProperties;
}

interface GridTableState {
  /** 判断是否需要展示loading态，主要适用于清除搜索条件之后展示 */
  isNeedLoading?: boolean;
}

export class GridTable extends React.Component<GridTableProps, GridTableState> {
  constructor(props) {
    super(props);
    this.state = {
      isNeedLoading: false
    };
  }

  render() {
    let {
        columns,
        rowDisabled,
        addons = [],
        emptyTips = '暂无数据',
        isNeedPagination = false,
        actionOptions,
        listModel,
        isNeedTableBorder = false,
        recordKey = 'id',
        records,
        isNeedCard = true,
        cardStyle
      } = this.props,
      { performSearch, fetch } = actionOptions,
      { list, query } = listModel;

    /**
     * 判断是否需要展示loading态
     * 1. list.fetched 不为true => 适用于刚进来页面，没有数据
     * 2. list.fetchState 为 Fetching => 表示列表正在拉取
     * 3. list.data.recordCount => 结合1、2判断，列表需要轮询，则不需要展示loading
     * 4. query.search => 如果是有搜索关键词，则需要展示loadint
     * 5. this.state.isNeedLoadint => 结合4，清除搜索条件之后，也需要展示loading
     */
    const isShowLoading: boolean =
      (list.fetched !== true || list.fetchState === FetchState.Fetching) &&
      (list.data.recordCount === 0 || !!query.search || this.state.isNeedLoading);

    // 数据的来源，如果自己传了records，则采取传入的records，主要针对排序
    const finalRecords = records ? records : list.data.records;

    const content: JSX.Element = (
      <React.Fragment>
        <Table
          columns={columns}
          records={finalRecords}
          recordKey={recordKey}
          rowDisabled={rowDisabled}
          topTip={null}
          addons={addons.concat([
            autotip({
              isLoading: isShowLoading,
              isError: list.fetchState === FetchState.Failed,
              isFound: !!query.search,
              onClear: () => {
                // 清除搜索条件，需要展示loading态，表示正在拉取数据
                this.setState({ isNeedLoading: true });
                // 清楚搜索关键词
                performSearch('');
                // 这里需要重置一下 isNeedLoading的状态，不然清除搜索条件之后，列表还会展示loading的状态
                setTimeout(() => {
                  this.setState({ isNeedLoading: false });
                }, 500);
              },
              onRetry: () => {
                fetch();
              },
              foundKeyword: query.search,
              emptyText: query.search ? null : <StatusTip status="empty" emptyText={emptyTips} />
            })
          ])}
        />
        {isNeedPagination && this._renderPagination()}
      </React.Fragment>
    );

    return isNeedCard ? (
      <Card bordered={isNeedTableBorder} style={cardStyle}>
        <Card.Body>{content}</Card.Body>
      </Card>
    ) : (
      content
    );
  }

  /** 渲染翻页 */
  private _renderPagination() {
    let { actionOptions, listModel } = this.props,
      { changePaging } = actionOptions,
      { query, list } = listModel,
      { pageIndex, pageSize } = query.paging;

    return (
      <Pagination
        pageIndex={pageIndex}
        pageSize={pageSize}
        pageSizeOptions={[10, 20, 30]}
        recordCount={list.data.recordCount}
        onPagingChange={query => {
          if (query.pageIndex > Math.ceil(list.data.recordCount / query.pageSize)) {
            query.pageIndex = 1;
          }
          changePaging(query);
        }}
      />
    );
  }
}
