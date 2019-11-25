import * as React from 'react';
import * as ReactDom from 'react-dom';
import { extend } from '@tencent/qcloud-lib';

export interface DatePickerProps extends React.Props<any> {
  /** 选中时间的起始日期 */
  selected?: {
    /**
     * 选中的起始日期，格式为 YYYY-MM-DD，支持宏：
     *   - '%TODAY' 表示今天
     *   - '%TODAY-1' 表示昨天
     */
    from: string;

    /**
     * 选中的结束日期，格式为 YYYY-MM-DD，支持宏：
     *   - '%TODAY' 表示今天
     *   - '%TODAY-1' 表示昨天
     */
    to: string;
  };

  /**
   * 预定义日期值的选项卡
   */
  tabs?: DatePickerTab[];

  /**
   * 未选中日期的时候显示的文案
   */
  placeHolder?: string;

  /**
   * 允许选择的时间范围限制
   */
  range?: {
    /**
     * 允许选择的最小日期范围
     */
    min?: string;

    /**
     * 允许选择的最大日期范围
     */
    max?: string;

    /**
     * 允许选择的最长时间区间
     */
    maxLength?: number;
  };

  /**
   * 选择模式
   *   - "range"    （默认）多选模式，允许选择指定的日期范围
   *   - "single"   单选模式，只允许选择一天
   *   - "duration" 区间选择模式，允许选择指定长度的区间
   */
  mode?: DatePickerMode;

  /**
   * 区间选择模式的情况下，选择的区间长度，默认为 1
   */
  durationLength?: number;

  /**
   * 区间选择模式下，区间的对其模式
   *    - "start" （默认）选择点击日期之后的 X 天，X 为区间长度
   *    - "end"   选择点击日期之前的 X 天，X 为区间长度
   */
  durationAnchor?: DatePickerDurationAnchor;

  /**
   * 是否禁用，默认为 false
   */
  disabled?: boolean;

  /**
   * 日期被更改时的事件
   */
  onPick?: (value: DatePickerValue) => void;

  /**
   * 弹出浮层的对齐方式
   *   - "left"  左对齐
   *   - "right" 右对齐
   */
  pull?: DatePickerPullMode;
}

export enum DatePickerMode {
  Single = 'single' as any,
  Range = 'range' as any,
  Duration = 'duration' as any
}
export enum DatePickerDurationAnchor {
  Start = 'start' as any,
  End = 'end' as any
}
export enum DatePickerPullMode {
  Left = 'left' as any,
  Right = 'right' as any
}

/**
 * 表示一个快捷日期选项
 */
export interface DatePickerTab {
  /**
   * 选项卡对应的起始日期，格式为 YYYY-MM-DD
   */
  from: string;

  /**
   * 选项卡对应的结束日期，格式为 YYYY-MM-DD
   */
  to: string;

  /**
   * 显示的标签
   */
  label?: string;

  /**
   * 选项卡选中后显示的文本
   */
  display?: string;
}

/**
 * 表示选中的日期范围
 */
export interface DatePickerValue {
  /**
   * 当前选中的起始日期
   */
  from: string;

  /**
   * 当前选中的结束日期
   */
  to: string;

  /**
   * 当前选中的日期区间长度（选中了多少天）
   */
  length: number;

  /**
   * 当前选中的选项卡
   */
  tab: DatePickerTab;
}

interface DatePickerBeeInstance extends nmc.BeeInstance {
  value: DatePickerValue;
}

/**
 * 日历组件，目前从 Bee 组件包装
 */
export class DatePicker extends React.Component<DatePickerProps, any> {
  dom = null;
  private bee: DatePickerBeeInstance;

  public render() {
    return (
      <div
        ref={ref => {
          this.dom = ref;
        }}
      />
    );
  }

  public getValue() {
    return this.bee.value;
  }

  componentDidMount() {
    if (!this.bee) {
      let Bee = seajs.require('qccomponent');
      let target = this.dom;

      target.innerHTML = `<div b-tag="qc-date-picker"></div>`;

      let beeProps = {
        $data: extend({}, this.props)
      };

      this.bee = Bee.mount(target.firstElementChild as HTMLElement, beeProps) as DatePickerBeeInstance;
    }
  }

  componentWillUnmount() {
    if (this.bee) {
      this.bee.$destroy();
      this.bee = null;
    }
  }
}
