import * as React from 'react';
import { inspect } from 'util';
import { Modal } from '../modal';
import { Alert } from '../alert';
import { ColumnProps } from './tableProps';


function saveCache(key: string, value: object) {
  const uin = /(^|;)\s*uin=(\w+);?/.exec(document.cookie)[1];
  localStorage.setItem(`cat-${uin}-${key}`, inspect(value));
}

export interface FieldsManageProps {
  cacheKey: string;
  visible: boolean;
  columns: ColumnProps[];
  onOk: Function;
  onCancel: Function;
  selectedFieldsMaxNumber?: number;
}

interface FieldsManageState {
  checkedMap: object;
}

const MAX_COL = 6;

export class FieldsManage extends React.Component<FieldsManageProps, FieldsManageState> {
  static defaultProps = {
    selectedFieldsMaxNumber: 9,
  };

  constructor(props) {
    super(props);
    const checkedMap = (props.columns || [])
      .reduce((acc, cur) => ({...acc, [cur.key]: cur.visible}), {});
    this.state = {
      checkedMap,
    };
  }

  handleOk() {
    const {columns} = this.props;
    const {checkedMap} = this.state;
    saveCache(
      this.props.cacheKey,
      checkedMap,
    );
    this.props.onOk(checkedMap);
  }

  handleChange(e) {
    const {value, checked} = e.target;
    const {checkedMap} = this.state;
    checkedMap[value] = checked;

    this.setState({
      checkedMap,
    });
  }

  render() {
    const {visible, onCancel, columns} = this.props;
    const {checkedMap} = this.state;

    const columnGroups = (columns || []).reduce(
      (acc, item, idx) => {
        if (idx % MAX_COL === 0) {
          acc.push([item]);
        }
        else {
          acc[acc.length - 1].push(item);
        }

        return acc;
      },
      [],
    );

    const modalProps = {
      visible,
      onCancel,
      onOk: this.handleOk.bind(this),
      title: '自定义列表字段',
    };

    return (
      <Modal { ...modalProps }>
        <div className="customize-column">
          <Alert>
            { `请选择您想显示的列表详细信息，最多勾选${Math.min(columns.length,
              this.props.selectedFieldsMaxNumber)}个字段，已勾选${Object.keys(checkedMap).filter(
              key => checkedMap[key]).length}个。` }
          </Alert>
          <div className="list-wrap clearfix">
            {
              columnGroups.map((group, index) => (
                <ul className="list-mod" key={ `group${index}` }>
                  {
                    group.map((item, idx) => (
                      <li key={ idx }>
                        <input
                          type="checkbox"
                          className="tc-15-checkbox"
                          value={ item.key }
                          onChange={ e => this.handleChange(e) }
                          checked={ item.isRequired || checkedMap[item.key] }
                          disabled={ item.isRequired }
                        />
                        <label>{ item.title }</label>
                      </li>
                    ))
                  }
                </ul>
              ))
            }
          </div>
        </div>
      </Modal>
    );
  }
}
