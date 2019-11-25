/* eslint-disable */
import * as React from 'react';
import * as classNames from 'classnames';
import { OnOuterClick } from '@tencent/qcloud-lib';
import { Select, SelectRange } from './Select';
import * as ReactDOM from 'react-dom';
import { Bubble } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
export interface TimePickerValue {
  hour: number;
  minute: number;
  second: number;
}

export interface TimePickerRange {
  min?: string;
  max?: string;
}

export interface TimePickerProps {
  /**
   * 默认时间 - 'HH:mm:ss'
   */
  defaultValue?: string;

  /**
   * 指定时间（作为受控组件时使用）- 'HH:mm:ss'
   */
  value?: string;

  /**
   * 未选中日期的时候显示的文案
   */
  placeHolder?: string;

  /**
   * 是否禁用，默认为 false
   */
  disabled?: boolean;

  /**
   * 允许选择的时间范围限制
   */
  range?: TimePickerRange;

  /**
   * 时间被更改时触发
   */
  onChange?: (time: string) => void;
}

// disabledHours?: () => Array<number>;
// disabledMinutes?: (selectedHour: number) => Array<number>;
// disabledSeconds?: (selectedHour: number, disabledMinutes: number) => Array<number>;

const keys = {
  '38': 'up',
  '40': 'down',
  '13': 'enter'
};

export class TimePicker extends React.Component<TimePickerProps, any> {
  input = null;
  constructor(props: TimePickerProps) {
    super(props);

    let defaultValue = props.value || props.defaultValue;
    const value = this.check(this.parse(defaultValue)) ? this.parse(defaultValue) : null;

    this.state = {
      bubbleActive: false,
      active: false,
      inputValue: defaultValue ? this.format(defaultValue) : '',
      value
    };

    this.close = this.close.bind(this);
  }

  componentWillReceiveProps(nextProps: TimePickerProps) {
    // 受控组件
    if ('value' in nextProps) {
      const nextValue = this.parse(nextProps.value);
      let value = this.check(nextValue) ? nextValue : null;

      const range: TimePickerRange = nextProps.range || {};

      if (value !== null && range.min && this.format(value) < range.min) {
        value = this.parse(range.min);
        this.confirm(range.min);
        return;
      }

      if (value !== null && range.max && this.format(value) > range.max) {
        value = this.parse(range.max);
        this.confirm(range.max);
        return;
      }

      this.setState({ value, inputValue: this.format(value) });
    }
  }

  // getCurTime = (props: TimePickerProps) => {

  // }

  /**
   * 格式化
   */
  formatNum = (num: number): string => (num > 9 ? `${num}` : `0${num}`);

  format = (value: string | TimePickerValue | null): string => {
    if (typeof value === 'string') return value;
    if (!value) return '';
    const { hour, minute, second } = value;
    return `${this.formatNum(hour)}:${this.formatNum(minute)}:${this.formatNum(second)}`;
  };

  /**
   * 检验传入value值是否合法
   */
  check = (value: TimePickerValue): boolean => {
    if (!value) return false;
    if (!('hour' in value) || value.hour < 0 || value.hour > 23) return false;
    if (!('minute' in value) || value.minute < 0 || value.minute > 59) return false;
    if (!('second' in value) || value.second < 0 || value.second > 59) return false;
    return true;
  };

  /**
   * 将字符串解析为 TimePickerValue
   */

  parse = (value: string | TimePickerValue): TimePickerValue => {
    if (typeof value !== 'string') return value;
    if (!/^[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(value)) return null;
    return {
      hour: +value.substr(0, 2),
      minute: +value.substr(3, 2),
      second: +value.substr(6, 2)
    };
  };

  /**
   * 确认value，合法将调用 onChange
   */
  confirm = value => {
    if (!('value' in this.props)) {
      this.setState({ value });
    }
    if (this.props.onChange) {
      this.props.onChange(this.format(value));
    }

    this.setState({ inputValue: this.format(value) });
  };

  handleInputChange = (e): void => {
    const inputValue = e.target.value;

    if (/^[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(inputValue)) {
      const value = {
        hour: +inputValue.substr(0, 2),
        minute: +inputValue.substr(3, 2),
        second: +inputValue.substr(6, 2)
      };
      if (this.check(value)) {
        const range: TimePickerRange = this.props.range || {};
        if (
          (!('min' in range) || this.format(value) >= range.min) &&
          (!('max' in range) || this.format(value) <= range.max)
        ) {
          this.confirm(value);
        }
      }
      this.setState({ bubbleActive: false });
    } else {
      this.setState({ bubbleActive: true });
    }

    if (inputValue.length === 0) {
      this.setState({ bubbleActive: false });
    }

    this.setState({ inputValue });
  };

  open = () => {
    if ('value' in this.props) {
      this.setState({ value: this.parse(this.props.value) });
    }
    if (!this.state.active && !this.props.disabled) {
      this.setState({ active: true });
    }
  };

  @OnOuterClick
  public close() {
    const { inputValue } = this.state;
    let value = this.state.value;

    if ('value' in this.props) {
      value = this.props.value;
      this.setState({ inputValue: this.format(value) });
    }

    // 非法输入检测
    if (!/^[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(inputValue)) {
      this.setState({ inputValue: this.format(value) });
    }

    const curValue = {
      hour: +inputValue.substr(0, 2),
      minute: +inputValue.substr(3, 2),
      second: +inputValue.substr(6, 2)
    };
    if (!this.check(curValue)) {
      this.setState({ inputValue: this.format(value) });
    }

    this.setState({ active: false, bubbleActive: false });
  }

  handleSelect = (type: string, val: number): void => {
    let value = this.state.value || { hour: 0, minute: 0, second: 0 };
    switch (type) {
      case 'hour':
        value.hour = val;
        break;
      case 'minute':
        value.minute = val;
        break;
      case 'second':
        value.second = val;
    }

    // 根据范围自调
    const range: TimePickerRange = this.props.range || {};
    let times = 0;
    while (('min' in range && this.format(value) < range.min) || ('max' in range && this.format(value) > range.max)) {
      if (type === 'hour') value.minute = (value.minute + 1) % 60;
      if (type === 'minute') value.second = (value.second + 1) % 60;
      if (++times > 60) break;
    }

    this.confirm(value);
  };

  /**
   * 处理键位按键事件
   */
  handleKeyDown = (e, hourRange, minuteRange, secondRange) => {
    if (!keys[e.keyCode]) return;
    e.preventDefault();

    const input = e.currentTarget as HTMLInputElement;

    let start = input.selectionStart,
      end = input.selectionEnd;
    setImmediate(() => input.setSelectionRange(start, end));

    let value = this.state.value || { hour: 0, minute: 0, second: 0 };
    let type;

    switch (keys[e.keyCode]) {
      case 'up':
        if (start >= 0 && end <= 2) {
          type = 'hour';
          value.hour = (value.hour + 1) % 24;
          if (
            ('min' in hourRange && value.hour < hourRange.min) ||
            ('max' in hourRange && value.hour > hourRange.max)
          ) {
            value.hour = hourRange.min || 0;
          }
        }
        if (start >= 3 && end <= 5) {
          type = 'minute';
          value.minute = (value.minute + 1) % 60;
          if (
            ('min' in minuteRange && value.minute < minuteRange.min) ||
            ('max' in minuteRange && value.minute > minuteRange.max)
          ) {
            value.minute = minuteRange.min || 0;
          }
        }
        if (start >= 6 && end <= 8) {
          type = 'second';
          value.second = (value.second + 1) % 60;
          if (
            ('min' in secondRange && value.second < secondRange.min) ||
            ('max' in secondRange && value.second > secondRange.max)
          ) {
            value.second = secondRange.min || 0;
          }
        }
        break;
      case 'down':
        if (start >= 0 && end <= 2) {
          type = 'hour';
          value.hour = (value.hour - 1 + 24) % 24;
          if (
            ('min' in hourRange && value.hour < hourRange.min) ||
            ('max' in hourRange && value.hour > hourRange.max)
          ) {
            value.hour = hourRange.max || 24;
          }
        }
        if (start >= 3 && end <= 5) {
          type = 'minute';
          value.minute = (value.minute - 1 + 60) % 60;
          if (
            ('min' in minuteRange && value.minute < minuteRange.min) ||
            ('max' in minuteRange && value.minute > minuteRange.max)
          ) {
            value.minute = minuteRange.max || 60;
          }
        }
        if (start >= 6 && end <= 8) {
          type = 'second';
          value.second = (value.second - 1 + 60) % 60;
          if (
            ('min' in secondRange && value.second < secondRange.min) ||
            ('max' in secondRange && value.second > secondRange.max)
          ) {
            value.second = secondRange.max || 60;
          }
        }
        break;
      case 'enter':
        this.close();
        (ReactDOM.findDOMNode(this.refs['input']) as HTMLElement).blur();
        return;
    }

    const range: TimePickerRange = this.props.range || {};

    while (('min' in range && this.format(value) < range.min) || ('max' in range && this.format(value) > range.max)) {
      // 根据范围修正
      if (type === 'hour') value.minute = (value.minute + 1) % 60;
      if (type === 'minute') value.second = (value.second + 1) % 60;
    }

    this.confirm(value);
  };

  render() {
    const value = this.state.value || { hour: 0, minute: 0, second: 0 };
    const { hour, minute, second } = value;

    const range: TimePickerRange = this.props.range || {};
    const minRange: TimePickerValue = this.parse(range.min) || { hour: 0, minute: 0, second: 0 };
    const maxRange: TimePickerValue = this.parse(range.max) || { hour: 23, minute: 59, second: 59 };

    let hourRange: SelectRange = { min: minRange.hour, max: maxRange.hour };
    let minuteRange: SelectRange = {},
      secondRange: SelectRange = {};

    if (hour === hourRange.min) {
      minuteRange.min = minRange.minute;
      if ('min' in minuteRange && minute === minuteRange.min) {
        secondRange.min = minRange.second;
      }
    }
    if (hour === hourRange.max) {
      minuteRange.max = maxRange.minute;
      if ('max' in minuteRange && minute === minuteRange.max) {
        secondRange.max = maxRange.second;
      }
    }

    const combobox = this.state.active ? (
      <div className="tc-time-picker-combobox">
        <Select
          from={0}
          to={23}
          value={hour}
          range={hourRange}
          onChange={value => {
            this.handleSelect('hour', value);
          }}
        />
        <Select
          from={0}
          to={59}
          value={minute}
          range={minuteRange}
          onChange={value => {
            this.handleSelect('minute', value);
          }}
        />
        <Select
          from={0}
          to={59}
          value={second}
          range={secondRange}
          onChange={value => {
            this.handleSelect('second', value);
          }}
        />
      </div>
    ) : null;

    return (
      <div className={classNames('tc-time-picker', { active: this.state.active })}>
        <div className="tc-time-picker-input-wrap">
          <Bubble
            placement="bottom"
            className="error"
            style={{ display: this.state.bubbleActive ? '' : 'none' }}
            content={t('格式错误，应为HH:mm:ss')}
          >
            <input
              type="text"
              ref={ref => {
                this.input = ref;
              }}
              className="tc-15-input-text shortest"
              onClick={this.open}
              onFocus={this.open}
              placeholder={this.props.placeHolder || t('时间选择')}
              value={this.state.inputValue}
              onChange={this.handleInputChange}
              disabled={this.props.disabled}
              onKeyDown={e => this.handleKeyDown(e, hourRange, minuteRange, secondRange)}
              maxLength={8}
            />
          </Bubble>
        </div>
        {combobox}
      </div>
    );
  }
}
