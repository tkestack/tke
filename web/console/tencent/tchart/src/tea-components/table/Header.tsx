import * as React from 'react';
import { MultiDropdown } from '../multidropdown';
import { ColGroup } from './ColGroup';
import { BaseRecord, HeaderProps, ColumnProps } from './tableProps';



interface HeaderState {}

export class Header<T extends BaseRecord> extends React.Component<HeaderProps, HeaderState> {
  constructor(props) {
    super(props);
    this.state = {};
  }

  render() {
    const {onSelectedAll, columns, check, checkedAll} = this.props;

    return (
      <div className="tc-15-table-fixed-head">
        <table className="tc-15-table-box">
          <ColGroup columns={ columns } check={ check }/>
          <thead>
          <tr>
            {
              check &&
              <th>
                <div className="tc-15-first-checkbox">
                  <input type="checkbox" className="tc-15-checkbox"
                         checked={ checkedAll }
                         onChange={ e => onSelectedAll(e) }
                  />
                </div>
              </th>
            }
            {
              (columns || []).map(item => (
                <th key={ item.key || item.dataIndex }>
                  {
                    item.filter
                      ? (
                        <MultiDropdown { ...item.filter.props } style={ {marginTop: 0} }>
                          <span className={ item.nowrap ? '' : 'text-overflow' }
                                style={ {verticalAlign: 'middle'} }>
                            { item.title }
                          </span>
                          <i className="filtrate-icon"></i>
                        </MultiDropdown>
                      )
                      : (
                        item.sorter ?
                          <div onClick={ e => item.onChange(
                            item.sortOrder === 'desc' ? 'asc' : 'desc',
                          ) }>
                            <span className={ item.nowrap ? '' : 'text-overflow' }
                                  style={ {cursor: 'pointer'} }>
                              { item.title }
                              { item.sortOrder === 'asc' ?
                                <i className="up-sort-icon"></i> : (
                                  item.sortOrder === 'desc' ?
                                    <i className="down-sort-icon"></i> :
                                    <i className="sort-icon"></i>
                                )
                              }
                            </span>
                          </div> : typeof item.title === 'string'
                          ? <div>
                              <span className={ item.nowrap ? '' : 'text-overflow' }>
                                { item.title }
                              </span>
                          </div>
                          : item.title
                      ) }
                </th>
              ))
            }
          </tr>
          </thead>
        </table>
      </div>
    );
  }
}
