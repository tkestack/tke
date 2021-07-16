import * as React from 'react';

export interface PureInputProps {
  inputValue: string;
  onChange?: (value: Array<any>) => void;
  onSelect?: (value: Array<any>) => void;
  offset: number;
}

const keys = {
  "9" : 'tab',
  "13": 'enter'
};

export class PureInput extends React.Component<PureInputProps, any> {

  componentDidMount() {
    const onChange = this.props.onChange;
    if (onChange) {
      onChange(this.getValue(this.props.inputValue));
    }
  }

  componentWillReceiveProps(nextProps: PureInputProps) {
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
          onSelect(this.getValue(this.props.inputValue).filter(i => !!i.name));
        }
        return false;
    }
  }

  getValue(value: string): Array<any> {
    return value.split('|').map(item => { return { name: item.trim() } });
  }

  render() {
    return null;
  }
}