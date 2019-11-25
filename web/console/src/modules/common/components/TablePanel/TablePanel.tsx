import * as React from 'react';
import { FetchState } from '@tencent/qcloud-redux-fetcher';
import { autotip } from '@tencent/tea-component/lib/table/addons/autotip';
import { StatusTip } from '@tencent/tea-component/lib/tips';
import { ListModel, ListAction } from '@tencent/redux-list';

import {
  CardProps,
  CardBodyProps,
  Card,
  Justify,
  Text,
  Table,
  TableColumn,
  TableProps,
  Dropdown,
  List,
  Pagination,
  Bubble,
  Icon
} from '@tea/component';
import {
  stylize,
  StylizeOption,
  sortable,
  SortBy,
  selectable,
  radioable,
  filterable,
  expandable,
  ExpandableAddonOptions,
  SelectableOptions,
  RadioableOptions,
  FilterOption,
  scrollable,
  ScrollableAddonOptions
} from '@tea/component/table/addons';

import { t, Trans } from '@tencent/tea-app/lib/i18n';

interface TablePanelProps extends TablePanelBodyProps, TablePanelHeaderProps, StylizeOption {
  title?: React.ReactNode;
  operation?: React.ReactNode;

  // className?: string;
  // style?: React.CSSProperties;

  warringTips?: React.ReactNode;
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
  model: ListModel;
  /** fetcher、query相关的action，包含select、selects、clear等 */
  action: ListAction;

  /** 空列表的相关提示，不传默认为 暂无数据 */
  emptyTips?: React.ReactNode;

  selectable?: SelectableOptions;
  radioable?: RadioableOptions;
  expandable?: ExpandableAddonOptions;
  scrollable?: ScrollableAddonOptions;
  isNeedPagination?: boolean;
}
function TablePanelBody({
  isNeedCard = true,
  cardProps,
  cardBodyProps,
  getOperations,
  //保证3个2字按钮能在一行内显示
  operationsWidth = 140,
  emptyTips,

  className,
  style,
  headClassName,
  headStyle,
  bodyClassName,
  bodyStyle,

  ...props
}: TablePanelProps) {
  let {
    model: { list, query },
    action
  } = props;

  let [isNeedLoading, setIsNeedLoading] = React.useState(false);
  /**
   * 判断是否需要展示loading态
   * 1. list.fetched 不为true => 适用于刚进来页面，没有数据
   * 2. list.fetchState 为 Fetching => 表示列表正在拉取
   * 3. list.data.recordCount => 结合1、2判断，列表需要轮询，则不需要展示loading
   * 4. query.search => 如果是有搜索关键词，则需要展示loadint
   * 5. this.state.isNeedLoadint => 结合4，清除搜索条件之后，也需要展示loading
   */
  let isShowLoading: boolean =
    (list.fetched !== true || list.fetchState === FetchState.Fetching) &&
    (list.data.recordCount === 0 || !!query.search || isNeedLoading);

  props.records = list.data.records;
  props.recordKey = props.recordKey || 'id';

  props.addons = props.addons || [];

  let [sorts, setSorts] = React.useState([]);
  let [filters, setFilters] = React.useState([]);

  let filteredRecords = list.data.records.slice();
  // 如果要在前端排序，可以用 sortable.comparer 生成默认的排序方法
  if (!props.onSort) {
    filteredRecords.sort(sortable.comparer(sorts));
  }

  if (props.scrollable) {
    props.addons.push(scrollable(props.scrollable));
  }

  if (props.columns) {
    let { columns, addons } = formatColumn(props, {
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
  let table = (
    <Table
      records={filteredRecords}
      {...props}
      addons={props.addons.concat([
        autotip({
          isLoading: isShowLoading,
          isError: list.fetchState === FetchState.Failed,
          isFound: !!query.search,
          onClear: () => {
            // 清除搜索条件，需要展示loading态，表示正在拉取数据
            setIsNeedLoading(true);
            // 清楚搜索关键词
            action.performSearch('');
            // 这里需要重置一下 isNeedLoading的状态，不然清除搜索条件之后，列表还会展示loading的状态
            setTimeout(() => {
              setIsNeedLoading(false);
            }, 500);
          },
          onRetry: () => {
            action.fetch();
          },
          foundKeyword: query.search,
          emptyText: query.search ? null : <StatusTip status="empty" emptyText={emptyTips} />
        })
      ])}
    />
  );

  // let table = createTable(props);
  return isNeedCard ? (
    <Card {...cardProps}>
      <Card.Body {...cardBodyProps}>
        {table}
        {props.isNeedPagination && <TablePanelPagination {...props} />}
      </Card.Body>
    </Card>
  ) : (
    <React.Fragment>
      {table}
      {props.isNeedPagination && <TablePanelPagination {...props} />}
    </React.Fragment>
  );
}

function formatColumn<Record = any>({ columns, onSort, action }: TablePanelProps, stateObj) {
  let columnsFormat: TableColumn<Record>[] = [],
    addons = [],
    sortableColumns = [];
  columns.forEach((config, index) => {
    let columnInfo = {
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
              {config.header}
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
          options: config.filterable.options
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
  let column4Operations: TableColumn<Record> = {
    key: 'operations',
    header: t('操作'),
    width: operationsWidth,
    render: (record: Record, rowKey: string, recordIndex: number, column: TableColumn<Record>) => {
      let ops = getOperations(record, rowKey, recordIndex, column);
      if (ops.length > 3) {
        let nodes = ops.splice(0, 2);
        let more = (
          <Dropdown button={t('更多')}>
            <List type="option">
              {ops.map((operation, index) => {
                return <List.Item key={index}>{operation}</List.Item>;
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
  let { query, list } = model,
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
