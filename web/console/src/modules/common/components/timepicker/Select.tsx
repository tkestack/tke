/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import * as React from 'react';
import * as classNames from 'classnames';
import * as ReactDOM from 'react-dom';

export interface SelectProps {
  /**
   * 当前选择的值
   */
  value: number;
  /**
   * 起止值
   */
  from: number;
  to: number;

  range?: SelectRange;
  onChange?: (value: number) => void;
}

export interface SelectRange {
  min?: number;
  max?: number;
}

/**
 * （受控组件）
 */
export class Select extends React.Component<SelectProps, any> {
  dom = null;
  list = null;

  componentDidMount() {
    this.scrollToSelected(0);
  }

  componentDidUpdate() {
    this.scrollToSelected(150);
  }

  /**
   * 滚动到指定元素
   */
  scrollTo = (element, to, duration) => {
    if (duration <= 0) {
      element.scrollTop = to;
      return;
    }
    let difference = to - element.scrollTop;
    let perTick = (difference / duration) * 10;

    setTimeout(() => {
      element.scrollTop = element.scrollTop + perTick;
      if (element.scrollTop === to) return;
      this.scrollTo(element, to, duration - 10);
    }, 10);
  };

  /**
   * 滚动到当前选择元素
   */
  scrollToSelected = duration => {
    const index = this.props.value;
    const topOption = this.list.children[index] as HTMLElement;
    const to = topOption.offsetTop - this.dom.offsetTop;

    this.scrollTo(this.dom, to, duration);
  };

  /**
   * 根据范围生成列表
   */
  genRangeList = (start: number, end: number): Array<String> =>
    Array(end - start + 1)
      .fill(0)
      .map((e, i) => {
        const num = i + start;
        return num > 9 ? `${num}` : `0${num}`;
      });

  handleSelect = (e, val: number): void => {
    e.stopPropagation();
    if (this.props.onChange) this.props.onChange(val);
  };

  render() {
    const { from, to, value, range } = this.props;
    const list = this.genRangeList(from, to).map((item, i) => {
      if (range && 'min' in range && i < range.min) {
        return (
          <li key={i} className="disabled">
            {item}
          </li>
        );
      }
      if (range && 'max' in range && i > range.max) {
        return (
          <li key={i} className="disabled">
            {item}
          </li>
        );
      }
      return (
        <li key={i} className={+item === value ? 'current' : ''} onClick={e => this.handleSelect(e, +item)}>
          {item}
        </li>
      );
    });

    return (
      <div
        className="tc-time-picker-select"
        ref={ref => {
          this.dom = ref;
        }}
      >
        <ul
          ref={ref => {
            this.list = ref;
          }}
        >
          {list}
        </ul>
      </div>
    );
  }
}
