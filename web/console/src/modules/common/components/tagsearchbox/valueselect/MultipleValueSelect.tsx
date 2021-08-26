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
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface MultipleValueSelectProps {
  values: Array<any>;
  inputValue: string;
  onChange: (value: Array<any>) => void;
  onSelect: (value: Array<any>) => void;
  onCancel: () => void;
  offset: number;
}

export interface MultipleValueSelectState {
  curIndex: number;
  select: Array<number>;
}

const keys = {
  '8': 'backspace',
  '9': 'tab',
  '13': 'enter',
  '37': 'left',
  '38': 'up',
  '39': 'right',
  '40': 'down'
};

export class MultipleValueSelect extends React.Component<MultipleValueSelectProps, any> {
  constructor(props) {
    super(props);

    const list = this.props.inputValue.split('|').map(i => i.trim());
    const select = [],
      values = this.props.values;

    values.forEach((item, index) => {
      if (list.indexOf(item.name) >= 0) {
        select.push(index);
      }
    });

    this.state = {
      curIndex: -1,
      select
    };
  }

  componentDidMount() {
    const select = this.state.select;
    if (select.length <= 0 && this.props.onSelect) {
      this.props.onSelect(this.getValue(select));
    }
  }

  componentWillReceiveProps(nextProps: MultipleValueSelectProps) {
    if (this.props.inputValue !== nextProps.inputValue) {
      const list = nextProps.inputValue.split('|').map(i => i.trim());
      const select = [],
        values = nextProps.values;
      values.forEach((item, index) => {
        if (list.indexOf(item.name) >= 0) {
          select.push(index);
        }
      });
      this.setState({ select });
    }
  }

  // 父组件调用
  handleKeyDown = (keyCode: number): boolean => {
    if (!keys[keyCode]) return;
    const { onSelect, onChange } = this.props;
    const { curIndex, select } = this.state;

    const pos = select.indexOf(curIndex);
    switch (keys[keyCode]) {
      case 'tab':
        if (curIndex < 0) return false;
        if (pos >= 0) {
          select.splice(pos, 1);
        } else {
          select.push(curIndex);
        }

        if (onChange) {
          onChange(this.getValue(select));
        }
        return false;

      case 'enter':
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
  };

  getValue(select: Array<number>): Array<any> {
    const { values } = this.props;
    return select.map(i => values[i]);
  }

  move = (step: number): void => {
    const curIndex = this.state.curIndex;
    const { values, inputValue } = this.props;
    if (values.length <= 0) return;
    this.setState({ curIndex: (curIndex + step + values.length) % values.length });
  };

  handleClick = (e, index: number): void => {
    e.stopPropagation();
    const select = this.state.select;
    const onChange = this.props.onChange;

    const pos = select.indexOf(index);
    if (pos >= 0) {
      select.splice(pos, 1);
    } else {
      select.push(index);
    }

    if (onChange) {
      onChange(this.getValue(select));
    }
  };

  handleSubmit = (e): void => {
    e.stopPropagation();
    const onSelect = this.props.onSelect;
    const select = this.state.select;
    if (onSelect) {
      onSelect(this.getValue(select));
    }
  };

  handleCancel = (e): void => {
    e.stopPropagation();
    const onCancel = this.props.onCancel;
    if (onCancel) {
      onCancel();
    }
  };

  render() {
    const { select, curIndex } = this.state;
    const { inputValue, values, offset } = this.props;

    const list = values.map((item, index) => {
      let input = (
        <label className="form-ctrl-label">
          <input type="checkbox" readOnly checked={false} className="tc-15-checkbox" />
          {item.name}
        </label>
      );
      if (select.indexOf(index) >= 0) {
        input = (
          <label className="form-ctrl-label">
            <input type="checkbox" readOnly checked={true} className="tc-15-checkbox" />
            {item.name}
          </label>
        );
      }
      if (curIndex === index) {
        return (
          <li
            role="presentation"
            key={index}
            className="autocomplete-cur"
            onMouseDown={e => this.handleClick(e, index)}
          >
            {input}
          </li>
        );
      }
      return (
        <li role="presentation" key={index} onMouseDown={e => this.handleClick(e, index)}>
          {input}
        </li>
      );
    });

    if (list.length === 0) return null;

    let maxHeight = document.body.clientHeight ? document.body.clientHeight - 400 : 450;
    maxHeight = maxHeight > 240 ? maxHeight : 240;

    return (
      <div className="tc-15-autocomplete" style={{ left: `${offset}px` }}>
        <ul className="tc-15-autocomplete-menu" role="menu" style={{ maxHeight: `${maxHeight}px` }}>
          {list}
        </ul>
        <div className="tc-autocomplete-ft">
          <a href="javascript:;" className="autocomplete-btn" onClick={this.handleSubmit}>
            {t('完成')}
          </a>
          <a href="javascript:;" className="autocomplete-btn" onClick={this.handleCancel}>
            {t('取消')}
          </a>
        </div>
      </div>
    );
  }
}
