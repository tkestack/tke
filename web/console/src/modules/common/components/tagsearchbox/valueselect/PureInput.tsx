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
