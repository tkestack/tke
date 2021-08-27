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

import { FetchState, FFListAction, FFListModel, uuid, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import {
  Bubble,
  Card,
  CardBodyProps,
  CardProps,
  Dropdown,
  Icon,
  Justify,
  JustifyProps,
  List,
  Pagination,
  PaginationProps,
  Table,
  TableColumn,
  TableProps,
  Text,
  Button,
  ExternalLink
} from 'tea-component';
import {
  expandable,
  ExpandableAddonOptions,
  filterable,
  FilterableConfig,
  FilterOption,
  radioable,
  RadioableOptions,
  scrollable,
  ScrollableAddonOptions,
  selectable,
  SelectableOptions,
  sortable,
  SortBy,
  stylize,
  StylizeOption
} from 'tea-component/es/table/addons';
import { StatusTip } from 'tea-component';

import { CamBox, isCamRefused } from '../cam';

const { autotip } = Table.addons;

insertCSS(
  '@tencent/ff-component/operationListButton',
  `
  .tke-operation-list-button button{
  width: 100%;
  text-align: left;
}
`
);

export interface TablePanelProps extends TablePanelBodyProps, TablePanelHeaderProps, StylizeOption {
  title?: React.ReactNode;
  operation?: React.ReactNode;

  // className?: string;
  // style?: React.CSSProperties;

  warringTips?: React.ReactNode;
  errorText?: React.ReactNode;
  retryText?: React.ReactNode;
  onSort?: (sorts: SortBy[]) => void;
}

export function TablePanel<Record = any>({ ...props }: TablePanelProps) {
  return (
    <React.Fragment>
      <TablePanelHeader {...props} />
      <TablePanelBody {...props} />
    </React.Fragment>
  );
}

interface TablePanelHeaderProps {
  left?: React.ReactNode;
  right?: React.ReactNode;
  headerClass?: string;
  headerStyle?: React.CSSProperties;
}
function TablePanelHeader({ left, right, headerClass, headerStyle }: TablePanelHeaderProps) {
  if (left || right) {
    return (
      <Table.ActionPanel>
        <Justify className={headerClass} style={headerStyle} left={left} right={right} />
      </Table.ActionPanel>
    );
  } else {
    return <React.Fragment />;
  }
}

export interface TablePanelColumnProps<Record = any> extends TableColumn<Record> {
  sortable?: boolean;
  filterable?: {
    // type: 'single' | 'multiple';
    options: FilterOption[];
    all?: FilterOption | false;
    onChange?: (value: string) => void;
  };
  headerTips?: React.ReactNode;
}
interface TablePanelBodyProps<Record = any> extends TableProps, TablePanelPaginationProps {
  isNeedCard?: boolean;
  cardProps?: CardProps;
  cardBodyProps?: CardBodyProps;
  columns: TablePanelColumnProps<Record>[];

  //支持单选，多选，排序，筛选，操作按钮

  getOperations?: (
    record: Record,
    rowKey: string,
    recordIndex: number,
    column: TableColumn<Record>
  ) => React.ReactNode[];
  operationsWidth?: number;

  /** 列表的相关配置，包含list、query等 */
  model: FFListModel;
  /** fetcher、query相关的action，包含select、selects、clear等 */
  action: FFListAction;

  /** 空列表的相关提示，不传默认为 暂无数据 */
  emptyTips?: React.ReactNode;

  /** autoTips的onRetry，自定义重试的逻辑 */
  onRetry?: () => void;

  selectable?: SelectableOptions;
  radioable?: RadioableOptions;
  expandable?: ExpandableAddonOptions;
  scrollable?: ScrollableAddonOptions;
  isNeedPagination?: boolean;
  isNeedContinuePagination?: boolean;
}
function TablePanelBody({
  isNeedCard = true,
  cardProps,
  cardBodyProps,
  getOperations,
  //保证3个2字按钮能在一行内显示
  operationsWidth = 140,
  emptyTips,

  onRetry,

  className,
  style,
  headClassName,
  headStyle,
  bodyClassName,
  bodyStyle,

  topTip,

  ...props
}: TablePanelProps) {
  const {
    model: { list, query },
    action
  } = props;
  const isCAMError = list.fetchState === FetchState.Failed && isCamRefused(list.error);

  const [isNeedLoading, setIsNeedLoading] = React.useState(false);

  // query.search 为纯文本搜索， query.filter.searchBoxValues 为tagSearchbox  tagSearch 在searchFilter字段
  const searchFilterKeys = query.searchFilter ? Object.keys(query.searchFilter) : [];

  const search =
    query.search ||
    (query.filter &&
      query.filter.searchBoxValues &&
      query.filter.searchBoxValues.length &&
      query.filter.searchBoxValues.map(
        ({ attr, values }) => `${attr.name}:${values.map(({ name }) => name).join(',')}`
      )) ||
    (query.searchFilter &&
      searchFilterKeys.length !== 0 &&
      searchFilterKeys
        .filter(key => query.searchFilter[key] !== null)
        .map(key => `${key}:${query.searchFilter[key]}`)
        .join(','));
  /**
   * 判断是否需要展示loading态
   * 1. list.fetchState 为 Fetching => 表示列表正在拉取
   * 2. list.data.recordCount => 结合1、2判断，列表需要轮询，则不需要展示loading
   * 3. search => 如果是有搜索关键词，则需要展示loadint
   * 4. this.state.isNeedLoadint => 结合4，清除搜索条件之后，也需要展示loading
   */
  const isShowLoading: boolean =
    list.fetchState === FetchState.Fetching &&
    (list.data.recordCount === 0 ||
      !!(
        query.search ||
        (query.filter && query.filter.searchBoxValues && query.filter.searchBoxValues.length) ||
        (query.searchFilter &&
          searchFilterKeys.length &&
          searchFilterKeys.some(key => query.searchFilter[key] !== null))
      ) ||
      isNeedLoading);

  props.records = list.data.records;
  props.recordKey = props.recordKey || 'id';

  props.addons = props.addons || [];

  const [sorts, setSorts] = React.useState([]);
  const [filters, setFilters] = React.useState([]);

  const filteredRecords = list.data.records.slice();
  // 如果要在前端排序，可以用 sortable.comparer 生成默认的排序方法
  if (!props.onSort) {
    filteredRecords.sort(sortable.comparer(sorts));
  }

  if (props.scrollable) {
    props.addons.push(scrollable(props.scrollable));
  }

  if (props.columns) {
    const { columns, addons } = formatColumn(props, {
      sorts,
      filters,
      setSorts,
      setFilters
    });
    props.columns = columns;
    props.addons = props.addons.concat(addons);
  }

  if (props.selectable) {
    props.addons.push(selectable(props.selectable));
  }
  if (props.radioable) {
    props.addons.push(radioable(props.radioable));
  }

  if (props.expandable) {
    props.addons.push(expandable(props.expandable));
  }

  if (getOperations) {
    props.columns.push(createOperationColumn(getOperations, operationsWidth));
  }

  if (className || style || headClassName || headStyle || bodyClassName || bodyStyle) {
    props.addons.unshift(
      stylize({
        className,
        style,
        headClassName,
        headStyle,
        bodyClassName,
        bodyStyle
      })
    );
  }
  const table = (
    <Table
      {...props}
      records={filteredRecords}
      topTip={topTip || (isCAMError ? <CamBox message={list.error.message || list.error.data.message} /> : null)}
      addons={
        topTip || isCAMError
          ? null
          : props.addons.concat([
              autotip({
                isLoading: isShowLoading,
                isError: list.fetchState === FetchState.Failed,
                isFound: !!search,
                onClear: () => {
                  // 清除搜索条件，需要展示loading态，表示正在拉取数据
                  setIsNeedLoading(true);
                  // 清楚搜索关键词
                  query.search && action.performSearch('');

                  query.filter.searchBoxValues &&
                    query.filter.searchBoxValues.length &&
                    action.applyFilter({ searchBoxValues: [] });

                  if (query.searchFilter && searchFilterKeys.some(key => query.searchFilter[key] !== null)) {
                    const nextFilter = {};
                    searchFilterKeys.forEach(key => {
                      nextFilter[key] = null;
                    });
                    action.applySearchFilter(nextFilter);
                  }

                  // 这里需要重置一下 isNeedLoading的状态，不然清除搜索条件之后，列表还会展示loading的状态
                  setTimeout(() => {
                    setIsNeedLoading(false);
                  }, 500);
                },
                onRetry: () => {
                  onRetry ? onRetry() : action.fetch();
                },
                foundKeyword: search,
                emptyText: search ? null : <StatusTip status="empty" emptyText={emptyTips} />,
                errorText: props.errorText,
                retryText: props.retryText
              })
            ])
      }
    />
  );

  // let table = createTable(props);
  return isNeedCard ? (
    <Card {...cardProps}>
      <Card.Body {...cardBodyProps}>
        {table}
        {props.isNeedPagination && <TablePanelPagination {...props} />}
        {props.isNeedContinuePagination && !query.search && <TablePanelContinuePagination {...props} />}
      </Card.Body>
    </Card>
  ) : (
    <React.Fragment>
      {table}
      {props.isNeedPagination && <TablePanelPagination {...props} />}
      {props.isNeedContinuePagination && !query.search && <TablePanelContinuePagination {...props} />}
    </React.Fragment>
  );
}

function formatColumn<Record = any>({ columns, onSort, action }: TablePanelProps, stateObj) {
  const columnsFormat: TableColumn<Record>[] = [],
    addons = [],
    sortableColumns = [];
  columns.forEach((config, index) => {
    const columnInfo = {
      key: config.key,
      header: config.header,
      width: config.width,
      render: config.render
        ? config.render
        : (record, rowKey, recordIndex, column) => {
            return (
              <Text overflow parent="div" title={record[column.key] + ''}>
                {record[column.key] + ''}
              </Text>
            );
          }
    };
    if (config.headerTips) {
      if (config.header) {
        if (typeof config.header === 'function') {
          columnInfo.header = columns => {
            return (
              <React.Fragment>
                {(config.header as Function)(columns)}
                <Bubble content={config.headerTips || null}>
                  <Icon type="info" />
                </Bubble>
              </React.Fragment>
            );
          };
        } else {
          columnInfo.header = columns => (
            <React.Fragment>
              <Text overflow>{config.header}</Text>
              <Bubble content={config.headerTips || null}>
                <Icon type="info" />
              </Bubble>
            </React.Fragment>
          );
        }
      }
    }
    columnsFormat.push(columnInfo);
    if (config.filterable) {
      addons.push(
        filterable({
          type: 'single',
          column: config.key,
          value: stateObj.filters[config.key],
          onChange: value => {
            stateObj.filters[config.key] = value;
            stateObj.setFilters(stateObj.filters);

            if (config.filterable.onChange) {
              config.filterable.onChange(value);
            } else {
              action.applyFilter({ [config.key]: value });
            }
          },
          all: config.filterable.all,
          options: config.filterable.options,
          searchable: config.filterable.options.length > 10
        })
      );
    }

    if (config.sortable) {
      sortableColumns.push(config.key);
    }
  });
  if (sortableColumns.length) {
    addons.push(
      sortable({
        columns: sortableColumns,
        value: stateObj.sorts,
        onChange: value => {
          stateObj.setSorts(value);
          onSort && onSort(value);
        }
      })
    );
  }
  return {
    columns: columnsFormat,
    addons
  };
}

function createOperationColumn<Record>(
  getOperations: (
    record: Record,
    rowKey: string,
    recordIndex: number,
    column: TableColumn<Record>
  ) => React.ReactNode[],
  operationsWidth: number
) {
  const column4Operations: TableColumn<Record> = {
    key: 'operations',
    header: t('操作'),
    width: operationsWidth,
    render: (record: Record, rowKey: string, recordIndex: number, column: TableColumn<Record>) => {
      let ops = getOperations(record, rowKey, recordIndex, column);
      if (ops.length > 3) {
        const nodes = ops.splice(0, 2);
        const more = (
          <Dropdown button={t('更多')}>
            <List type="option">
              {ops.map((operation, index) => {
                return (
                  <List.Item key={index} className="tke-operation-list-button">
                    {operation}
                  </List.Item>
                );
              })}
            </List>
          </Dropdown>
        );
        ops = [...nodes, more];
      }
      return <React.Fragment>{ops}</React.Fragment>;
    }
  };
  return column4Operations;
}

interface TablePanelPaginationProps {
  // recordCount?: number;
  paginationProps?: {
    /**
     * 切换页长时，是否自动把页码重设为 1
     * @default true
     */
    isPagingReset?: boolean;
    /** 分页组件显示的说明信息，不传将渲染数据个数信息 */
    stateText?: React.ReactNode;
    /**
     * 支持的页长设置
     * @default [10, 20, 30, 50]
     */
    pageSizeOptions?: number[];
    /**
     * 是否显示状态文案
     * @default true
     */
    stateTextVisible?: boolean;
    /**
     * 是否显示页长选择
     * @default true
     */
    pageSizeVisible?: boolean;
    /**
     * 是否显示页码输入
     * @default true
     */
    pageIndexVisible?: boolean;
    /**
     * 是否显示切页按钮（上一页/下一页）
     * @default true
     */
    jumpVisible?: boolean;
    /**
     * 是否显示第一页和最后一页按钮
     * @default true
     */
    endJumpVisible?: boolean;
  };
}

function TablePanelPagination({ model, action, paginationProps }: TablePanelProps) {
  const { query, list } = model,
    { pageIndex, pageSize } = query.paging;

  return (
    <Pagination
      pageIndex={pageIndex}
      pageSize={pageSize}
      recordCount={list.data.recordCount}
      onPagingChange={query => {
        if (query.pageIndex > Math.ceil(list.data.recordCount / query.pageSize)) {
          query.pageIndex = 1;
        }
        action.changePaging(query);
      }}
      pageSizeOptions={(paginationProps && paginationProps.pageSizeOptions) || [10, 20, 30, 40, 50]}
      {...paginationProps}
    />
  );
}

function TablePanelContinuePagination({ model, action, paginationProps }: TablePanelProps) {
  const { query, list } = model,
    { pageIndex, pageSize } = query.paging;

  let totalCount = 0;
  const finalPages = list.pages ? list.pages : [];
  finalPages.forEach(listItem => {
    totalCount += listItem.data.recordCount;
  });

  const stateText = (
    <Text>
      <Text verticalAlign="middle">{t('第 {{pageIndex}} 页', { pageIndex: pageIndex })}</Text>
      {pageIndex > 1 && (
        <Trans>
          <Text theme="warning" className="tea-ml-2n" verticalAlign="middle">
            分页信息来源于数据快照，进行修改操作后，请
          </Text>
          <Button
            type="link"
            className="tea-ml-1n tea-mr-1n"
            onClick={() => {
              action.resetPaging();
            }}
          >
            返回列表首页
          </Button>
          <Text theme="warning" verticalAlign="middle">
            进行刷新。
          </Text>
          <ExternalLink
            style={{ verticalAlign: 'middle' }}
            href="https://kubernetes.io/docs/reference/using-api/api-concepts/#retrieving-large-results-sets-in-chunks"
          >
            查看原因
          </ExternalLink>
        </Trans>
      )}
    </Text>
  );

  return (
    <Pagination
      pageIndex={pageIndex}
      pageSize={pageSize}
      pageIndexVisible={false}
      endJumpVisible={false}
      stateText={stateText}
      recordCount={
        (list.data.continue || list.data.continueToken) && list.fetchState !== FetchState.Fetching
          ? Number.MAX_SAFE_INTEGER
          : totalCount
      }
      onPagingChange={query => {
        if (pageSize === query.pageSize) {
          action.changePagingIndex(query.pageIndex);
        } else {
          action.changePaging(
            Object.assign({}, query, {
              append: false,
              clear: true
            })
          );
        }
      }}
      pageSizeOptions={(paginationProps && paginationProps.pageSizeOptions) || [5, 10, 20, 50, 100]}
      {...paginationProps}
    />
  );
}
