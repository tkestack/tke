import * as React from 'react';
import { PageProps } from '../libs/paginationProps';

interface SelectItem {
  disabled?: boolean;
  value: string;
  label: string;
}

export interface HeaderProps {
  checkedAll?: boolean;
  columns: ColumnProps[];
  check?: boolean;
  onSelectedAll?: Function;
}

export interface RowSelection {
  onChange?: Function;
  selectedRowKeys?: string[];
}

export interface TableBodyProps<T> {
  style?: object;
  className?: string;
  loading?: boolean;
  check?: boolean;
  columns: ColumnProps[];
  records: T[];
  onRowSelected: Function;
  prefix?: Function;
  total?: number;
  selectedRowKeys: string[];
  onRowClick?: Function;
  onRowMouseEnter?: Function;
  onRowMouseLeave?: Function;
  onBodyMouseLeave?: Function;
}

export interface TableProps<T> {
  loading?: boolean;
  prefix?: Function;
  check?: boolean;
  style?: object;
  bodyStyle?: object;
  className?: string;
  bodyClassName?: string;
  pageOption?: PageProps;
  columns: ColumnProps[];
  records: T[];
  rowSelection?: RowSelection;
  onRowClick?: Function;
  onRowMouseEnter?: Function;
  onRowMouseLeave?: Function;
  onBodyMouseLeave?: Function;
}

export interface BaseRecord {
  key?: string;
  [x: string]: any;
}

export interface ColumnProps {
  type?: string;
  title: string | React.ReactNode;
  nowrap?: boolean;
  align?: string;
  visible: boolean;
  isRequired?: boolean;
  dataIndex?: string;
  key?: string | number;
  render?: Function;
  sorter?: boolean;
  sortOrder?: string;
  width?: number | string;
  onChange?: Function;
  onFilterVisibleChange?: Function;
  filter?: {
    props: {
      value: any;
      options: SelectItem[];
      onChange: Function;
    };
    multi: boolean;
  };
}
