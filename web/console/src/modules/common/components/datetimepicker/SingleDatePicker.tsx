/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import * as classNames from 'classnames';
import * as React from 'react';

import { Bubble } from '@tea/component';
import { OnOuterClick } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface SingleDatePickerValue {
  year: number;
  month: number;
  day: number;
}

export interface SingleDatePickerRange {
  min?: string;
  max?: string;
}

export interface SingleDatePickerProps {
  defaultValue?: string | SingleDatePickerValue;
  value?: string | SingleDatePickerValue;

  placeHolder?: string;

  version?: string;

  disabled?: boolean;

  range?: SingleDatePickerRange;

  onChange?: (date: string) => void;
}

const keys = {
  '13': 'enter'
};

const MONTH = [
  '',
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December'
];

export class SingleDatePicker extends React.Component<SingleDatePickerProps, any> {
  input = null;
  constructor(props) {
    super(props);

    let defaultValue = this.props.value || this.props.defaultValue;
    defaultValue = this.parse(defaultValue);

    const value = this.check(defaultValue) ? defaultValue : null;

    let cur: any = {};
    if (value) {
      cur = {
        curYear: value.year,
        curMonth: value.month
      };
    } else {
      cur = this.getCurYearAndMonth(this.props.range);
    }

    this.state = {
      value,
      active: false,
      bubbleActive: false,
      curYear: cur.curYear,
      curMonth: cur.curMonth,
      inputValue: this.format(value)
    };

    this.close = this.close.bind(this);
  }

  componentWillReceiveProps(nextProps: SingleDatePickerProps) {
    // 受控组件
    if ('value' in nextProps) {
      const nextValue = this.parse(nextProps.value);
      const value = this.check(nextValue) ? nextValue : null;
      this.setState({ value: nextProps.value, inputValue: this.format(value) });

      if (value) {
        this.setState({ curYear: value.year, curMonth: value.month });
      } else {
        this.setState(this.getCurYearAndMonth(nextProps.range));
      }
    } else {
      this.setState(this.getCurYearAndMonth(nextProps.range));
    }
  }

  /**
   * 获取当前年月
   */
  getCurYearAndMonth = range => {
    const date = new Date();
    let curYear = date.getFullYear(),
      curMonth = date.getMonth();

    if (range) {
      if (range.max) {
        const max = this.parse(range.max);
        curYear = max.year;
        curMonth = max.month;
      }
      if (range.min) {
        const min = this.parse(range.min);
        curYear = min.year;
        curMonth = min.month;
      }
    }
    return { curYear, curMonth };
  };

  /**
   * 格式化
   */
  formatNum = (num: number): string => (num > 9 ? `${num}` : `0${num}`);

  format = (value: string | SingleDatePickerValue | null): string => {
    if (typeof value === 'string') return value;
    if (!value) return '';
    const { year, month, day } = value;
    return `${year}-${this.formatNum(month + 1)}-${this.formatNum(day)}`;
  };

  getGlobalFormat = (year: number, month: number): string => {
    const version = this.props.version || window['VERSION'];
    if (version === 'en') {
      return `${MONTH[month]} ${year}`;
    } else {
      return `${year}年${this.formatNum(month)}月`;
    }
  };

  /**
   * 检验传入value值是否合法
   */
  check = (value: SingleDatePickerValue): boolean => {
    if (!value) return false;
    // TODO
    if (!('year' in value) || !('month' in value) || !('day' in value)) return false;

    const { year, month, day } = value;
    let date = new Date(year, month, day);
    if (date.getFullYear() !== year || date.getMonth() !== month || date.getDate() !== day) return false;
    return true;
  };

  /**
   * 字符串解析为 SingleDatePickerValue
   */
  parse = (value: string | SingleDatePickerValue): SingleDatePickerValue => {
    if (typeof value !== 'string') return value;
    if (!/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}$/.test(value)) return null;
    return {
      year: +value.substr(0, 4),
      month: +value.substr(5, 2) - 1,
      day: +value.substr(8, 2)
    };
  };

  /**
   * 确认value，合法将调用 onChange
   */
  confirm = (value: SingleDatePickerValue): void => {
    const range: SingleDatePickerRange = this.props.range || {};

    if (range.min && this.format(value) < range.min) return;
    if (range.max && this.format(value) > range.max) return;

    if (!('value' in this.props)) {
      this.setState({ value });
    }

    if (this.props.onChange) {
      this.props.onChange(this.format(value));
    }

    this.setState({
      inputValue: this.format(value),
      curYear: value.year,
      curMonth: value.month
    });
  };

  handleInputChange = (e): void => {
    const inputValue = e.target.value;

    if (/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}$/.test(inputValue)) {
      const value = {
        year: +inputValue.substr(0, 4),
        month: +inputValue.substr(5, 2) - 1,
        day: +inputValue.substr(8, 2)
      };
      if (this.check(value)) {
        this.confirm(value);
        this.setState({ bubbleActive: false });
      }
    } else {
      this.setState({ bubbleActive: true });
    }

    if (inputValue.length === 0) {
      this.setState({ bubbleActive: false });
    }

    this.setState({ inputValue });
  };

  open = () => {
    const value = this.parse(this.props.value);
    if ('value' in this.props) {
      this.setState({ value });
      const curValue = this.check(value) ? value : null;
      if (curValue) {
        this.setState({ curYear: curValue.year, curMonth: curValue.month });
      } else {
        this.setState(this.getCurYearAndMonth(this.props.range));
      }
    } else {
      this.setState(this.getCurYearAndMonth(this.props.range));
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

    if (!/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}$/.test(inputValue)) {
      this.setState({ inputValue: this.format(value) });
    }

    const curValue = {
      year: +inputValue.substr(0, 4),
      month: +inputValue.substr(5, 2) - 1,
      day: +inputValue.substr(8, 2)
    };

    if (!this.check(curValue)) {
      this.setState({ inputValue: this.format(value) });
    }

    this.setState({ active: false, bubbleActive: false });
  }

  handleKeyDown = (e): void => {
    if (!keys[e.keyCode]) return;
    e.preventDefault();

    switch (keys[e.keyCode]) {
      case 'enter':
        this.close();
        this.input.blur();
        break;
    }
  };

  handleSelect = (year: number, month: number, day: number): void => {
    this.confirm({ year, month, day });
    setTimeout(() => this.close(), 100);
  };

  calendarRender = (year: number, month: number, range: SingleDatePickerRange) => {
    let value = this.state.value;

    if ('value' in this.props) {
      value = this.parse(this.props.value);
    }

    if (!this.check(value)) {
      value = {};
    }

    // 本月第一天
    const firstDate = new Date(year, month, 1);
    let day = firstDate.getDay();

    // 获取本月有多少天
    const lastDate = new Date(year, month + 1, 0);
    const count = lastDate.getDate();

    const weeks = [];
    for (let i = 1; i <= count; i = i + 0) {
      const week = [];
      // 第一周补全
      for (let j = 0; j < day % 7; ++j) {
        week.push(<td key={`dis-${j}`} className="tc-15-calendar-dis" />);
      }
      // 填充一周
      do {
        (day => {
          if (day === value.day && month === value.month && year === value.year) {
            week.push(
              <td key={day} className="tc-15-calendar-today" onClick={() => this.handleSelect(year, month, day)}>
                {day}
              </td>
            );
          } else {
            const cur = this.format({ year, month, day });
            if ((range.min && cur < range.min) || (range.max && cur > range.max)) {
              week.push(
                <td key={day} className="tc-15-calendar-dis">
                  {day}
                </td>
              );
            } else {
              week.push(
                <td key={day} onClick={() => this.handleSelect(year, month, day)}>
                  {day}
                </td>
              );
            }
          }
        })(i++);
        if (i > count) break;
      } while (++day % 7 !== 0);

      weeks.push(<tr key={`week-${i / 7}`}>{week}</tr>);
    }
    return weeks;
  };

  prevBtnRender = (range: SingleDatePickerRange) => {
    const { curYear, curMonth } = this.state;

    if (range.min) {
      const minRange = this.parse(range.min);
      if (minRange.year >= curYear && minRange.month >= curMonth) {
        return (
          <i tabIndex={0} className="tc-15-calendar-i-pre-m disabled">
            <b>
              <span>{t('过去时间不可选')}</span>
            </b>
          </i>
        );
      }
    }

    return (
      <i
        tabIndex={0}
        className="tc-15-calendar-i-pre-m"
        onClick={() => {
          if (curMonth > 0) this.setState({ curMonth: curMonth - 1 });
          else this.setState({ curMonth: 11, curYear: curYear - 1 });
        }}
      >
        <b>
          <span>{t('转到上个月')}</span>
        </b>
      </i>
    );
  };

  nextBtnRender = (range: SingleDatePickerRange) => {
    const { curYear, curMonth } = this.state;

    if (range.max) {
      const maxRange = this.parse(range.max);
      if (maxRange.year <= curYear && maxRange.month <= curMonth) {
        return (
          <i tabIndex={0} className="tc-15-calendar-i-next-m disabled">
            <b>
              <span>{t('未来时间不可选')}</span>
            </b>
          </i>
        );
      }
    }

    return (
      <i
        tabIndex={0}
        className="tc-15-calendar-i-next-m"
        onClick={() => {
          if (curMonth < 11) this.setState({ curMonth: curMonth + 1 });
          else this.setState({ curMonth: 0, curYear: curYear + 1 });
        }}
      >
        <b>
          <span>{t('转到下个月')}</span>
        </b>
      </i>
    );
  };

  render() {
    const { curYear, curMonth } = this.state;

    let range: SingleDatePickerRange = {};
    if (this.props.range) {
      range.min = this.props.range.min;
      range.max = this.props.range.max;
    }

    return (
      <div className="tc-15-calendar-select-wrap tc-15-calendar2-hook">
        <div className={classNames('tc-15-calendar-select', 'tc-15-calendar-single', { show: this.state.active })}>
          <Bubble
            placement="bottom"
            className="error"
            style={{ display: this.state.bubbleActive ? '' : 'none' }}
            content={t('格式错误，应为yyyy-MM-dd')}
          >
            <input
              ref={ref => {
                this.input = ref;
              }}
              disabled={this.props.disabled}
              className="tc-15-simulate-select m show"
              value={this.state.inputValue}
              onChange={this.handleInputChange}
              onClick={this.open}
              onFocus={this.open}
              placeholder={this.props.placeHolder || t('日期选择')}
              onKeyDown={this.handleKeyDown}
              maxLength={10}
            />
          </Bubble>
          <div className="tc-15-calendar-triangle-wrap" />
          <div className="tc-15-calendar-triangle" />
          <div className="tc-15-calendar tc-15-calendar2">
            <div className="tc-15-calendar-cont">
              <table cellSpacing={0} className="tc-15-calendar-left">
                <caption>{this.getGlobalFormat(curYear, curMonth + 1)}</caption>
                <thead>
                  <Trans>
                    <tr>
                      <th>日</th>
                      <th>一</th>
                      <th>二</th>
                      <th>三</th>
                      <th>四</th>
                      <th>五</th>
                      <th>六</th>
                    </tr>
                  </Trans>
                </thead>
                <tbody>
                  <tr>
                    <td colSpan={3}>{this.prevBtnRender(range)}</td>
                    <td colSpan={4}>{this.nextBtnRender(range)}</td>
                  </tr>
                  {this.calendarRender(curYear, curMonth, range)}
                </tbody>
              </table>
            </div>
            <div className="tc-15-calendar-for-style" />
          </div>
        </div>
      </div>
    );
  }
}
