import * as React from 'react';


interface BaseReactProps {
  key?: string;
  defaultValue?: string;
  children?: React.ReactNode;
  className?: string;
  placeholder?: string;
  style?: object;
}

export interface SelectItem {
  disabled?: boolean;
  value: string;
  label: string;
}

interface MetricConfig {
  metricShowName: string;
  calcType: string;
  calcValue: string;
  continueTime: string;
  alarmNotifyType: string;
  alarmNotifyPeriod: string;
  unit: string;
}

export interface SelectProps extends BaseReactProps {
  disabled?: boolean;
  value?: string;
  onChange?: Function;
  options: SelectItem[];
}

interface SelectState {
}

export class Select extends React.Component<SelectProps, SelectState> {
  constructor(props) {
    super(props);
  }

  handleChange(e) {
    const { value } = e.target;
    const { onChange } = this.props;
    onChange && onChange(value);
  }

  render() {
    const { options, value, className, onChange, ...opts } = this.props;

    return (
      <select
        { ...opts }
        value={value}
        className={`tc-15-select ${className || ''}`}
        onChange={this.handleChange.bind(this)}
      >
        {
          options.map(item => (
            <option key={item.value} value={item.value}>{item.label}</option>
          ))
        }
      </select>
    );
  }
}
