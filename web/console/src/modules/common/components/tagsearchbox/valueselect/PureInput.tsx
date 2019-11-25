import * as React from 'react';
import * as classNames from 'classnames';

export interface PureInputtProps {
  inputValue: string;
  onChange?: (value: Array<any>) => void;
  onSelect?: (value: Array<any>) => void;
  offset: number;
}

const keys = {
  '9': 'tab',
  '13': 'enter'
};

export class PureInput extends React.Component<PureInputtProps, any> {
  componentDidMount() {
    const onChange = this.props.onChange;
    if (onChange) {
      onChange(this.getValue(this.props.inputValue));
    }
  }

  componentWillReceiveProps(nextProps: PureInputtProps) {
    if (this.props.inputValue !== nextProps.inputValue) {
      const onChange = nextProps.onChange;
      if (onChange) {
        onChange(this.getValue(nextProps.inputValue));
      }
    }
  }

  // 父组件调用
  handleKeyDown = (keyCode: number): boolean => {
    if (!keys[keyCode]) return;
    const { onSelect, inputValue } = this.props;

    switch (keys[keyCode]) {
      case 'tab':
      case 'enter':
        if (inputValue.length <= 0) return false;
        if (onSelect) {
          onSelect(this.getValue(this.props.inputValue));
        }
        return false;
    }
  };

  getValue(value: string): Array<any> {
    return value
      .split('|')
      .filter(item => item.trim().length > 0)
      .map(item => {
        return { name: item.trim() };
      });
  }

  render() {
    return null;
  }
}
