import * as React from "react";

export interface ColumnType {
  key: string,
  name: string,
  render?: any,
}

interface Props {
  className: string,
  style: object,
  columns: Array<ColumnType>,
  list: Array<any>,
  update: Function,
}

interface State {
  checkedAll: boolean
}

export default class InstanceList extends React.Component<Props, State> {

  constructor(props) {
    super(props);
    this.state = {
      checkedAll: props.list.every((item) => item.isChecked),
    };
  }

  onCheck(instance, event) {
    const list = [].concat(this.props.list.map(item => {
      if (item === instance) {
        item.isChecked = event.target.checked;
      }
      return item;
    }));
    const checkedAll = list.every((item) => item.isChecked);
    this.setState({checkedAll}, () => {
      this.props.update(list);
    });
  };

  onCheckAll(event) {
    const list = [].concat(this.props.list.map(item => {
      item.isChecked = event.target.checked;
      return item;
    }));
    this.setState({checkedAll: event.target.checked}, () => {
      this.props.update(list);
    });
  }

  render() {
    const { className, style } = this.props;

    return (
      <div
        className={ className }
        style={ style }
      >
        <div>
          <table className="Tchart_table-box">
            <thead>
            <tr>
              <th className="">
                <div className="tc-15-first-checkbox">
                  <input
                    type="checkbox"
                    className="tc-15-checkbox"
                    checked={this.state.checkedAll}
                    onChange={this.onCheckAll.bind(this)}
                  />
                </div>
              </th>
              {
                this.props.columns.map(item => {
                  return (
                    <th className="" key={item.key}>
                      <div>
                        <span className="text-overflow">
                          { item.name }
                        </span>
                      </div>
                    </th>
                  )
                })
              }
            </tr>
            </thead>
            <tbody>
            {
              this.props.list.map((item, index) => {
                return (
                  <tr key={ index }>
                    <th className="">
                      <div className="tc-15-first-checkbox">
                        <input
                          type="checkbox"
                          className="tc-15-checkbox"
                          checked={item.isChecked}
                          onChange={this.onCheck.bind(this, item)}
                        />
                      </div>
                    </th>
                    {
                      this.props.columns.map(column => {
                        return (
                          <th className="" key={column.key}>
                              <div>
                                {
                                  column.render? column.render(item) : item[column.key]
                                }
                              </div>
                          </th>
                        )
                      })
                    }
                  </tr>
                );
              })
            }
            </tbody>
          </table>
        </div>
      </div>
    )
  }
}