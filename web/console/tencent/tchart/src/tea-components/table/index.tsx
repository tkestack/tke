import * as React from 'react';

import { BaseRecord, TableProps, TableBodyProps } from './tableProps';
import { Pagination } from '../pagination';

import { Header } from './Header';
import { Body } from './Body';

interface TableState {
}

export class Table<T extends BaseRecord> extends React.Component<TableProps<T>, TableState> {
  constructor(props) {
    super(props);
    this.state = {
    };
  }

  handleSelectedAll(e) {
    const { checked } = e.target;
    const { onChange } = this.props.rowSelection;
    onChange(checked ? (this.props.records || []).map(item => item.key) : []);
  }

  render() {
    const {
      style,
      check,
      prefix,
      records = [],
      loading,
      className,
      bodyStyle,
      rowSelection,
      bodyClassName,
      columns: cols,
      pageOption = false,
      onRowClick,
      onRowMouseEnter,
      onRowMouseLeave,
      onBodyMouseLeave,
    } = this.props;
    const { onChange: onSelectedChange = undefined, selectedRowKeys = [] } = rowSelection || {};
    const columns = (cols || []).filter(item => item.visible);

    return (
      <div className={`tc-15-table-panel ${className}`} style={style}>
        <Header columns={columns} check={check}
          checkedAll={records.length > 0 && (selectedRowKeys || []).length === records.length}
          onSelectedAll={this.handleSelectedAll.bind(this)}
        />
        <Body columns={columns}
          className={bodyClassName}
          style={bodyStyle}
          records={records}
          loading={loading}
          prefix={prefix}
          total={pageOption ? pageOption.total : records.length}
          check={check}
          selectedRowKeys={selectedRowKeys}
          onRowSelected={onSelectedChange}
          onRowClick={onRowClick}
          onRowMouseEnter={onRowMouseEnter}
          onRowMouseLeave={onRowMouseLeave}
          onBodyMouseLeave={onBodyMouseLeave}
        />
        {
          pageOption && <Pagination {...pageOption} />
        }
      </div>
    );
  }
}
