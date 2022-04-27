import * as React from 'react';
import classNames from 'classnames';
import { OnOuterClick } from "../libs/decorators/OnOuterClick";
import { SingleDatePicker, SingleDatePickerValue, SingleDatePickerRange } from './SingleDatePicker';
import { TimePicker, TimePickerValue, TimePickerRange } from '../timepicker';
import * as languages from '../../i18n';

const version = (window as any).VERSION || "zh";
const language = languages[version];

export interface DateTimePickerTab {
  from: Date | string;
  to: Date  | string;
  label: string;
}

export interface DateTimePickerValue {
  from?: Date | string;
  to?: Date | string;
}

export interface DateTimePickerRange {
  min?: Date | string;
  max?: Date | string;
  maxLength?: number;
}

export interface DateTimeRange {
  min?: Date;
  max?: Date;
}


export interface DateTimePickerProps {
  /**
   * 默认初始值 { from: string/Date, to: string/Date }
   */
  defaultValue?: DateTimePickerValue;

  /**
   * 预定义日期值的选项卡 { from: string/Date, to: string/Date, label: string }
   */
  tabs?: DateTimePickerTab[];

  /**
   * 默认选中的选项卡序号
   */
  defaultSelectedTabIndex?: number;

  /**
   * 选项卡与 Picker 是否联动，默认不联动
   */
  linkage?: boolean;

  /**
   * 选择指定长度的区间 **(单位为秒(s))** **（duration 动态改变时将重新初始化组件）**
   */
  duration?: number;

  /**
   * 未选中日期的时候显示的文案
   */
  placeHolder?: string;

  /**
   * 是否禁用，默认为 false
   */
  disabled?: boolean;

  /**
   * 允许选择的时间范围限制 { from: string/Date, to: string/Date, maxLength: number } **(maxLength 单位为秒(s))** **（range暂不支持动态改变）**
   */
  range?: DateTimePickerRange;

  /**
   * 日期时间被更改时触发
   */
  onChange?: (range: DateTimePickerValue, tabLabel?: string) => void;

  /**
   * 语言（’zh'/'en'），默认为 window.VERSION
   */
  version?: string;
}

/**
 * props中时间单位为s，乘unit转化为ms
 */
const unit = 1000;


/**
 * 日期时间范围选择组件
 * TODO 受控组件支持
 */
export class DateTimePicker extends React.Component<DateTimePickerProps, any> {

  constructor(props: DateTimePickerProps) {
    super(props);

    const { dateFrom, timeFrom, dateTo, timeTo, rangeFrom, rangeTo, pickerValue } = this.getInitDateAndRange(this.props);
    this.state = {
      active: false,
      selectedTabIndex: null,

      pickerValue: pickerValue,

      dateFrom, timeFrom,
      rangeFrom,

      dateTo, timeTo,
      rangeTo
    }
  }


  componentDidMount() {
    const { tabs, defaultSelectedTabIndex } = this.props;
    if ('tabs' in this.props && 'defaultSelectedTabIndex' in this.props && tabs[defaultSelectedTabIndex]) {
      this.handleTabSelect(tabs[defaultSelectedTabIndex], defaultSelectedTabIndex);
    }
  }


  componentWillReceiveProps(nextProps: DateTimePickerProps) {
    // duration可变
    if (!('duration' in nextProps) || nextProps.duration === this.props.duration) return;

    const { dateFrom, timeFrom, dateTo, timeTo, rangeFrom, rangeTo, pickerValue } = this.getInitDateAndRange(nextProps);
    this.setState({ dateFrom, timeFrom, dateTo, timeTo, rangeFrom, rangeTo, pickerValue, selectedTabIndex: null});

    // TODO range可变
  }

  /**
   * 根据 props 获得初始日期和范围
   */
  getInitDateAndRange = (props) => {
    let dateFrom, timeFrom, dateTo, timeTo, pickerValue = null;
    let defaultValue = props.defaultValue || {};

    // defaultValue
    if (defaultValue.from) {
      const from = this.parseMacro(defaultValue.from as string);
      dateFrom = this.formatDate(this.getDate(from));
      timeFrom = this.formatTime(this.getTime(from));
      pickerValue = { from, to: from }
    }

    if (defaultValue.to) {
      const to = this.parseMacro(defaultValue.to as string);
      dateTo = this.formatDate(this.getDate(to));
      timeTo = this.formatTime(this.getTime(to));
      if (!pickerValue) {
        pickerValue = { from: to, to }
      } else {
        pickerValue.to = to;
      }
    }

    // 获取初始范围
    const { rangeFrom, rangeTo } = this.getInitRange(props);
    const range = props.range || {};

    // duration 存在时更新起止时间
    if ('duration' in props && pickerValue) {
      pickerValue.to = new Date(pickerValue.from.getTime() + props.duration * unit);
      // 超出范围置空
      if (rangeFrom.max && this.compare(pickerValue.from, rangeFrom.max) > 0) {
        dateFrom = timeFrom = dateTo = timeTo = pickerValue = null;
        defaultValue = {};
      } else {
        dateTo = this.formatDate(this.getDate(pickerValue.to));
        timeTo = this.formatTime(this.getTime(pickerValue.to));
      }
    }

    // TODO
    // DefaultValue Range 判断

    // 根据 defaultValue 限定开始结束范围
    if (defaultValue.from) {
      const date = this.parse(`${dateFrom} ${timeFrom}`);
      if ((rangeFrom.min && this.compare(date, rangeFrom.min) < 0) || (rangeFrom.max && this.compare(date, rangeFrom.max) > 0)) return;
      rangeTo.min = date;
      if ('maxLength' in range) {
        rangeTo.max = new Date(date.getTime() +  range.maxLength * unit);
      }
    }
    if (defaultValue.to) {
      const date = this.parse(`${dateTo} ${timeTo}`);
      if ((rangeTo.min && this.compare(date, rangeTo.min) < 0) || (rangeTo.max && this.compare(date, rangeTo.max) > 0)) return;
      rangeFrom.max = date;
      // if ('maxLength' in range) {
      //     rangeFrom.min = new Date(date.getTime() - range.maxLength * unit);
      // }
    }
    return { dateFrom, timeFrom, dateTo, timeTo, rangeFrom, rangeTo, pickerValue };
  }

  /**
   * 根据 props 获得初始范围
   */
  getInitRange = (props) => {
    const { duration, range = {} } = props;

    let rangeFrom: DateTimeRange = {}, rangeTo: DateTimeRange = {};

    if ('min' in range) {
      const rangeMin = this.parseMacro(range['min'] as string);
      rangeFrom.min = rangeTo.min = rangeMin;
    }

    if ('max' in range) {
      const rangeMax = this.parseMacro(range['max'] as string);
      rangeTo.max = rangeMax;
      if ('duration' in props) {
        rangeFrom.max = new Date(rangeMax.getTime() - duration * unit);
      } else {
        rangeFrom.max = rangeMax;
      }
    }
    return { rangeFrom, rangeTo };
  }


  /**
   * 根据起止日期重置范围
   *   type - 'from'/'to' 表示起止
   */
  resetRange = (type: string, date: Date): void => {
    let { rangeFrom, rangeTo, dateTo, timeTo } = Object.assign({}, this.state);
    const { duration, range = {} } = this.props;

    if (type === 'from') {

      if ((rangeFrom.min && this.compare(date, rangeFrom.min) < 0) ||
        (rangeFrom.max && this.compare(date, rangeFrom.max) > 0)) {
        return;
      }

      rangeTo.min = date;

      const rangeMax = this.parseMacro(range['max'] as string);
      const rangeMin = this.parseMacro(range['min'] as string);

      if ('maxLength' in range) {
        const max = new Date(date.getTime() + range['maxLength'] * unit);
        rangeTo.max = this.compare(max, rangeMax) < 0 ? max : rangeMax;
        if (dateTo) {
          const to = this.parse(`${dateTo} ${timeTo}`);
          if (this.compare(to, rangeTo.min) < 0 || this.compare(to, rangeTo.max) > 0) {
            this.setState({ dateTo: null, timeTo: null });
          }
        }
      } else {
        rangeTo.max = rangeMax;
      }

      rangeFrom.min = rangeMin;

      if ('duration' in this.props) {
        rangeFrom.max = new Date(rangeMax.getTime() - duration * unit);
      }

    }

    if (type === 'to') {

      if ((rangeTo.min && this.compare(date, rangeTo.min) < 0) ||
        (rangeTo.max && this.compare(date, rangeTo.max) > 0)) {
        return;
      }

      if (!('maxLength' in range)) {
        rangeFrom.max = date;
      }

      // maxLength存在时不再限制rangeFrom.min

      // if ('maxLength' in range) {
      //   const min = new Date(date.getTime() - range.maxLength * unit);
      //   rangeFrom.min = this.compare(min, range.min) > 0 ? min : range.min;
      // }
    }

    this.setState({ rangeFrom, rangeTo });
  }

  /**
   * 判断时间是否选择完成
   */
  isCompleted = (): boolean => {
    const { dateFrom, timeFrom, dateTo, timeTo } = this.state;
    return dateFrom && timeFrom && dateTo && timeTo;
  }

  /**
   * Date比较，精确到s
   */
  compare = (value1: Date, value2: Date): number => {
    if (value1.getFullYear() > value2.getFullYear()) return 1;
    if (value1.getFullYear() < value2.getFullYear()) return -1;
    if (value1.getMonth() > value2.getMonth()) return 1;
    if (value1.getMonth() < value2.getMonth()) return -1;
    if (value1.getDate() > value2.getDate()) return 1;
    if (value1.getDate() < value2.getDate()) return -1;
    if (value1.getHours() > value2.getHours()) return 1;
    if (value1.getHours() < value2.getHours()) return -1;
    if (value1.getMinutes() > value2.getMinutes()) return 1;
    if (value1.getMinutes() < value2.getMinutes()) return -1;
    if (value1.getSeconds() > value2.getSeconds()) return 1;
    if (value1.getSeconds() < value2.getSeconds()) return -1;
    return 0;
  }

  /**
   * 格式化
   */
  formatNum = (num:number):string => num > 9 ? `${num}` : `0${num}`;
  formatDate = (value: string | SingleDatePickerValue | null): string => {
    if (typeof value === 'string') return value;
    if (!value) return "";
    const { year, month, day } = value;
    return `${year}-${this.formatNum(month+1)}-${this.formatNum(day)}`;
  }

  formatTime = (value: string | TimePickerValue | null): string => {
    if (typeof value === 'string') return value;
    if (!value) return "";
    const { hour, minute, second } = value;
    return `${this.formatNum(hour)}:${this.formatNum(minute)}:${this.formatNum(second)}`;
  }

  format = (date: Date): string => {
    return `${this.formatDate(this.getDate(date))}  ${this.formatTime(this.getTime(date))}`
  }

  /**
   * 解析宏指令
   */
  parseMacro = (macro: string, type?: string): Date => {

    if (typeof macro !== 'string') {
      return macro;
    }

    if (/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}\s[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(macro)) {
      return this.parse(macro);
    }

    const today = new Date();
    let unit, offset, stamp;

    if (macro.indexOf('%TODAY') >= 0) {
      macro = macro.replace('%TODAY', '');
      unit = 86400 * 1000;
      offset = +macro;
      const date = new Date(today.getTime() + (offset * unit));
      if (type === 'from') {
        date.setHours(0);
        date.setMinutes(0);
        date.setSeconds(0);
      }
      if (type === 'to') {
        date.setHours(23);
        date.setMinutes(59);
        date.setSeconds(59);
      }
      return date;
    }

    if (macro.indexOf('%NOW') >= 0) {
      macro = macro.replace('%NOW', '');
      unit = 1000;

      if (macro.indexOf('h') >= 0) {
        macro = macro.replace('h', '');
        unit = 3600 * 1000;
      }

      if (macro.indexOf('m') >= 0) {
        macro = macro.replace('m', '');
        unit = 60 * 1000;
      }

      if (macro.indexOf('s') >= 0) {
        macro = macro.replace('s', '');
      }

      offset = +macro;
      const date = new Date(today.getTime() + (offset * unit));
      return date;
    }
  }


  /**
   * 'yyyy-MM-dd HH:mm:ss'解析为 Date 对象
   * 解决IE不能使用 new Date('yyyy-MM-dd HH:mm:ss')
   */
  parse = (str: string) => {
    const date = new Date(NaN);
    if (!/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}\s[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(str)) return null;

    const year = +str.substr(0, 4),
      month = +str.substr(5, 2)-1,
      day = +str.substr(8, 2),
      hour = +str.substr(11, 2),
      minute = +str.substr(14, 2),
      second = +str.substr(17, 2);

    date.setFullYear(year, month, day);
    date.setHours(hour);
    date.setMinutes(minute);
    date.setSeconds(second);

    return date;
  }


  getDate = (date: Date) => {
    if (!date) return null;
    return {
      year: +date.getFullYear(),
      month: +date.getMonth(),
      day: +date.getDate()
    }
  }

  getTime = (date: Date) => {
    if (!date) return null;
    return {
      hour: +date.getHours(),
      minute: +date.getMinutes(),
      second: +date.getSeconds()
    }
  }

  /**
   * 展开选择框
   */
  open = (): void => {
    if (this.props.disabled) return;
    const value = this.state.pickerValue;
    if (value) {
      this.setState({
        dateFrom: this.formatDate(this.getDate(value.from)),
        timeFrom: this.formatTime(this.getTime(value.from)),
        dateTo: this.formatDate(this.getDate(value.to)),
        timeTo: this.formatTime(this.getTime(value.to))
      }, () => {
        this.resetRange('from', value.from);
        if (!('duration' in this.props)) {
          this.resetRange('to', value.to);
        }
      })
    }
    this.setState({ active: true });
  }

  /**
   * 点击确认时调用
   */
  handleSubmit = (): void => {
    if (!this.isCompleted()) return;
    const { dateFrom, timeFrom, dateTo, timeTo } = this.state;
    const pickerValue = {
      from: this.parse(`${dateFrom} ${timeFrom}`),
      to: this.parse(`${dateTo} ${timeTo}`),
    }

    this.setState({ active: false });

    const preValue = this.state.pickerValue;

    if (preValue && preValue.from && preValue.to &&
      this.compare(pickerValue.from, preValue.from) === 0 &&
      this.compare(pickerValue.to, preValue.to) === 0) return;

    if (this.compare(pickerValue.from, pickerValue.to) > 0) {
      return;
    }

    if (this.props.onChange) {
      this.props.onChange(pickerValue);
    }

    this.setState({
      pickerValue,
      selectedTabIndex: null,
      value: pickerValue
    });

  }

  /**
   * 点击取消时调用
   */
  @OnOuterClick
  handleCancel(): void {
    this.setState({ active: false });
  }

  /**
   * 点击Tab时调用
   */
  handleTabSelect = (tab: DateTimePickerTab, index: number): void => {
    const value = {
      from: this.parseMacro(tab.from as string, 'from'),
      to: this.parseMacro(tab.to as string, 'to')
    }

    if (this.props.onChange) {
      this.props.onChange(value, tab.label);
    }

    // 不联动
    if (!this.props.linkage) {
      this.setState({
        active: false,
        selectedTabIndex: index,
        pickerValue: null
      });
      return;
    }

    let dateFrom, timeFrom, dateTo, timeTo;
    dateFrom = this.formatDate(this.getDate(value.from));
    timeFrom = this.formatTime(this.getTime(value.from));
    this.resetRange('from', value.from);

    dateTo = this.formatDate(this.getDate(value.to));
    timeTo = this.formatTime(this.getTime(value.to));
    this.resetRange('to', value.to);

    this.setState({
      active: false,
      selectedTabIndex: index,
      dateFrom, timeFrom, dateTo, timeTo,
      pickerValue: value
    });
  }


  render() {
    const { pickerValue, dateFrom, timeFrom, dateTo, timeTo, selectedTabIndex, rangeFrom, rangeTo } = this.state;
    const isCompleted = this.isCompleted();

    // picker内容显示
    let text = this.props.placeHolder || language.OptionDate;
    if (pickerValue) {
      text = `${this.format(pickerValue.from)} ${language.To} ${this.format(pickerValue.to)}`;
    }

    // Tabs
    const tabs = this.props.tabs ? this.props.tabs.map((tab, index) => {
      if (selectedTabIndex !== index) {
        return <span role="tab" tabIndex={index} key={index} onClick={() => this.handleTabSelect(tab, index)}>{tab.label}</span>
      } else {
        return <span role="tab" tabIndex={index} key={index} className="current" onClick={() => this.handleTabSelect(tab, index)}>{tab.label}</span>
      }
    }) : [];

    // 时间起止范围
    const timeFromRange: TimePickerRange = {};
    if (dateFrom && dateFrom === this.formatDate(this.getDate(rangeFrom.min))) {
      timeFromRange.min = this.formatTime(this.getTime(rangeFrom.min));
    }
    if (dateFrom && dateFrom === this.formatDate(this.getDate(rangeFrom.max))) {
      timeFromRange.max = this.formatTime(this.getTime(rangeFrom.max));
    }

    const timeToRange: TimePickerRange = {};
    if (dateTo && dateTo === this.formatDate(this.getDate(rangeTo.min))) {
      timeToRange.min = this.formatTime(this.getTime(rangeTo.min));
    }
    if (dateTo && dateTo === this.formatDate(this.getDate(rangeTo.max))) {
      timeToRange.max = this.formatTime(this.getTime(rangeTo.max));
    }

    return (
      <div className="tc-15-calendar-select-wrap tc-15-calendar2-hook">
        <div role="tablist">
          {tabs}
        </div>
        <div className={classNames("tc-15-dropdown tc-15-dropdown-btn-style date-dropdown", { "tc-15-menu-active": this.state.active, "disabled": this.props.disabled})}>
          <a className="tc-15-dropdown-link" onClick={this.open} onFocus={this.open} ><i className="caret"></i>{text}</a>
          <div className="tc-15-dropdown-menu" role="menu">
            <div className="tc-custom-date">
              <div className="custom-date-wrap">
                <em>从</em>
                <div className="calendar-box">
                  <SingleDatePicker
                    value={dateFrom}
                    range={{min: this.formatDate(this.getDate(rangeFrom.min)), max: this.formatDate(this.getDate(rangeFrom.max))}}
                    onChange={(dateFrom) => {
                      const time = timeFrom ? timeFrom : (this.formatTime(this.getTime(rangeFrom.min)) || '00:00:00');
                      if ('duration' in this.props) {
                        const to = new Date(this.parse(`${dateFrom} ${time}`).getTime() + this.props.duration * unit);

                        if (this.compare(to, rangeTo.max) > 0) {
                          const from = new Date(rangeTo.max.getTime() - this.props.duration * unit);
                          this.setState({
                            dateFrom: this.formatDate(this.getDate(from)), timeFrom: this.formatTime(this.getTime(from)),
                            dateTo: this.formatDate(this.getDate(rangeTo.max)), timeTo: this.formatTime(this.getTime(rangeTo.max))
                          });
                        } else {
                          this.setState({
                            dateFrom, timeFrom: time,
                            dateTo: this.formatDate(this.getDate(to)), timeTo: this.formatTime(this.getTime(to))
                          });
                        }

                      } else {
                        this.resetRange('from', this.parse(`${dateFrom} ${time}`));
                        this.setState({ dateFrom, timeFrom: time })
                      }
                    }}
                    version={this.props.version}
                  />
                  <TimePicker
                    value={timeFrom}
                    range={timeFromRange}
                    onChange={(timeFrom) => {
                      if (dateFrom) {
                        this.resetRange('from', this.parse(`${dateFrom} ${timeFrom}`));
                        if ('duration' in this.props) {
                          const to = new Date(this.parse(`${dateFrom} ${timeFrom}`).getTime() + this.props.duration * unit);
                          this.setState({ dateTo: this.formatDate(this.getDate(to)), timeTo: this.formatTime(this.getTime(to)) });
                        }
                      }
                      this.setState({ timeFrom });
                    }}
                  />
                </div>
              </div>
              <div className="custom-date-wrap">
                <em>至</em>
                <div className="calendar-box">
                  <SingleDatePicker
                    value={dateTo}
                    disabled={'duration' in this.props}
                    range={{min: this.formatDate(this.getDate(rangeTo.min)), max: this.formatDate(this.getDate(rangeTo.max))}}
                    onChange={(dateTo) => {
                      const time = timeTo ? timeTo : (this.formatTime(this.getTime(rangeTo.min)) || '00:00:00');
                      this.resetRange('to', this.parse(`${dateTo} ${time}`));
                      this.setState({ dateTo, timeTo: time });
                    }}
                    version={this.props.version}
                  />
                  <TimePicker
                    value={timeTo}
                    disabled={'duration' in this.props}
                    range={timeToRange}
                    onChange={(timeTo) => {
                      if (dateTo) {
                        this.resetRange('to', this.parse(`${dateTo} ${timeTo}`));
                      }
                      this.setState({ timeTo });
                    }}
                  />
                </div>
              </div>
            </div>
            <div className="custom-date-ft">
              <button type="button" className="tc-15-btn m" onClick={this.handleSubmit}>{language.Confirm}</button>&nbsp;
              <button type="button" className="tc-15-btn m weak" onClick={this.handleCancel.bind(this)}>{language.Cancel}</button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}