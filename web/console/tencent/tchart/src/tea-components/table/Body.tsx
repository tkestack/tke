import * as React from 'react';

import { BaseRecord, TableBodyProps } from './tableProps';

import { ColGroup } from './ColGroup';

interface BodyState {
}

export class Body<T extends BaseRecord> extends React.Component<TableBodyProps<T>, BodyState> {
  constructor(props) {
    super(props);
    this.state = {};
  }

  handleRowChange(e) {
    const {onRowSelected, selectedRowKeys} = this.props;
    const {value, checked} = e.target;
    const index = selectedRowKeys.indexOf(value);
    if (checked && index === - 1) {
      selectedRowKeys.push(value);
    }
    else if (! checked && index > - 1) {
      selectedRowKeys.splice(index, 1);
    }

    onRowSelected(selectedRowKeys);
  }

  render() {
    const {
      selectedRowKeys,
      onRowSelected,
      check,
      columns,
      records,
      className,
      style,
      prefix,
      total,
      loading = false,
      onRowClick,
      onRowMouseEnter,
      onRowMouseLeave,
      onBodyMouseLeave,
    } = this.props;

    return (
      <div onMouseLeave={ onBodyMouseLeave as any } className={ `tc-15-table-fixed-body ${className}` } style={ style }>
        <table className="tc-15-table-box tc-15-table-rowhover" ref={ el => el && el.parentElement && (
          el.parentElement.style.overflowY = el.offsetHeight < el.parentElement.offsetHeight ? 'visible' : 'auto') }>
          <ColGroup columns={ columns } check={ check }/>
          <tbody>
          {
            (loading || prefix || (records || []).length === 0) &&
            <tr>
              <td colSpan={ columns.length + + (! ! check) }>
                <div className="text-center">
                  { loading ?
                    <span className="text-overflow">
                        <i className="n-loading-icon"></i>加载中
                      </span> : (prefix ?
                        <span className="text-overflow">
                          { prefix(total) }
                        </span> :
                        <span className="text-overflow">
                          { '没有记录' }
                        </span>
                    )
                  }
                </div>
              </td>
            </tr>
          }
          {
            (records || []).map((item, index) => (
              item.trRender ? item.trRender(item, index, records) :
                <tr key={ item.key || index }
                    onClick={ onRowClick && (() => onRowClick(item, index)) }
                    onMouseEnter={ onRowMouseEnter && (() => onRowMouseEnter(item, index)) }
                    onMouseLeave={ onRowMouseLeave && (() => onRowMouseLeave(item, index)) }
                >
                  {
                    check &&
                    <td>
                      <div className="tc-15-first-checkbox">
                        <input type="checkbox"
                               className="tc-15-checkbox"
                               value={ item.key }
                               checked={ selectedRowKeys.indexOf(item.key) > - 1 }
                               onChange={ this.handleRowChange.bind(this) }
                        />
                      </div>
                    </td>
                  }
                  {
                    columns.map(col => (
                      <td key={ col.key || col.dataIndex }>
                        { col.render ? col.render(item, index, records) :
                          <div>
                          <span className="text-overflow">
                            { item[col.dataIndex] }
                          </span>
                          </div>
                        }
                      </td>
                    ))
                  }
                </tr>
            ))
          }
          </tbody>
        </table>
      </div>
    );
  }
}
