import * as React from 'react';

export interface SingleValueSelectProps {
  values: Array<any>;
  inputValue: string;
  onChange?: (value: Array<any>) => void;
  onSelect?: (value: Array<any>) => void;
  offset: number;
}

export interface SingleValueSelectState {
  select: number;
}

const keys = {
  "8": 'backspace',
  "9": 'tab',
  "13": 'enter',
  "37": 'left',
  "38": 'up',
  "39": 'right',
  "40": 'down'
};

export class SingleValueSelect extends React.Component<SingleValueSelectProps, any> {

  constructor(props) {
    super(props);
    const { values, inputValue, onSelect } = this.props;

    let select = -1;
    values.forEach((item, index) => {
      if (item.name === inputValue) {
        select = index;
      }
    });

    this.state = {
      select,
    }
  }

  componentDidMount() {
    const select = this.state.select;
    if (select < 0 && this.props.onSelect) {
      this.props.onSelect(this.getValue(select));
    }
  }

  componentWillReceiveProps(nextProps: SingleValueSelectProps) {
    const { values, inputValue } = nextProps;
    const list = values.map(item => item.name);
    const select = list.indexOf(inputValue);
    this.setState({ select });
  }

  // 父组件调用
  handleKeyDown = (keyCode: number): boolean => {

    if (!keys[keyCode]) return;
    const { onSelect } = this.props;
    const select = this.state.select;

    switch (keys[keyCode]) {

      case 'enter':
      case 'tab':
        if (onSelect) {
          onSelect(this.getValue(select));
        }
        return false;

      case 'up':
        this.move(-1);
        break;

      case 'down':
        this.move(1);
        break;
    }
  }

  getValue(select: number): Array<any> {
    const { values, inputValue } = this.props;
    if (select < 0) {
      // return [];
      if (inputValue)
        return [{ name: inputValue }];
      return []
    }

    const list = values;
    if (select < list.length) {
      return [list[select]];
    } else {
      const select = list.map(item => item.name).indexOf(inputValue);
      this.setState({ select });
      if (select < 0) return [];
      return [list[select]];
    }
  }

  move = (step: number): void => {
    const select = this.state.select;
    const { values, inputValue } = this.props;
    const list = values;
    if (list.length <= 0) return;
    this.setState({ select: (select + step + list.length) % list.length });
  }


  handleClick = (e, index: number): void => {
    e.stopPropagation();
    if (this.props.onSelect) {
      this.props.onSelect(this.getValue(index));
    }
  }

  render() {
    const select = this.state.select;
    const { inputValue, values, offset } = this.props;

    const list = values.map((item, index) => {
      if (select === index) {
        return (
          <li role="presentation" key={index} className="autocomplete-cur" onClick={(e) => this.handleClick(e, index)}>
            <a className="text-truncate" role="menuitem" href="javascript:;" title={item.name} style={item.style || {}}>{item.name}</a>
          </li>
        )
      }
      return (
        <li role="presentation" key={index} onClick={(e) => this.handleClick(e, index)}>
          <a className="text-truncate" role="menuitem" href="javascript:;" title={item.name} style={item.style || {}}>{item.name}</a>
        </li>
      )
    });

    if (list.length === 0) {
      list.push(
        <li role="presentation" key={0}>
          <a className="autocomplete-empty" role="menuitem" href="javascript:;">相关值不存在</a>
        </li>
      );
    }

    let maxHeight = document.body.clientHeight ? document.body.clientHeight - 450 : 400;
    maxHeight = maxHeight > 240 ? maxHeight : 240;

    return (
      <div className="tc-15-autocomplete" style={{ left: `${offset}px`, width: 'auto', minWidth: 180 }}>
        <ul className="tc-15-autocomplete-menu" role="menu" style={{ maxHeight: `${maxHeight}px` }}>
          {list}
        </ul>
      </div>
    )
  }
}