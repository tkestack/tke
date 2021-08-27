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
import * as classnames from 'classnames';
import * as React from 'react';

import { Bubble } from '@tea/component';
import { BaseReactProps } from '@tencent/ff-redux';

import { Rule, Validate } from '../../../../../helpers/Validator';
import { Validation } from '../../models';

let deepEqual = require('deep-equal');

export interface InputFieldProps extends BaseReactProps {
  /**输入框类型 */
  type?: string;

  /**提示模式 inline: 单行显示, popup: 气泡提示   */
  tipMode?: 'inline' | 'popup';

  /**输入值 */
  value?: string | number;

  /**校验态 */
  validator?: Validation;

  /**输入事件 */
  onChange?: (value) => void;

  /**聚焦事件 */
  onFocus?: () => void;

  /**placeholder */
  placeholder?: string;

  /**失去焦点事件 */
  onBlur?: (value) => void;

  /**样式 */
  style?: React.CSSProperties;

  /**提示 */
  tip?: string | JSX.Element;

  /**输入框后面放置的操作 */
  ops?: string | JSX.Element;

  /**校验规则 */
  rule?: Rule;

  /**是否可编辑 */
  disabled?: boolean;

  /**不可填写tip */
  disabeldTip?: string;

  /**是否多行输入 */
  isMultiline?: boolean;

  bodyStyle?: React.CSSProperties;

  /** 冒泡展示的方向，default: left */
  popDirection?: 'left' | 'top' | 'bottom' | 'right';
}

interface InputFieldState {
  /**输入值 */
  data?: any;
}

export class InputField extends React.Component<InputFieldProps, InputFieldState> {
  constructor(props, context) {
    super(props, context);

    this.state = {
      data: props.value
    };
  }

  shouldComponentUpdate(nextProps, nextState) {
    return (
      this.props.value !== nextProps.value ||
      this.props.type !== nextProps.type ||
      !deepEqual(this.props.validator, nextProps.validator) ||
      !deepEqual(nextState, this.state)
    );
  }

  componentWillReceiveProps(nextProps) {
    if (!deepEqual(this.props.value, nextProps.value)) {
      this.setState({ data: nextProps.value });
    }
  }

  render() {
    const {
      type,
      tipMode,
      value,
      validator,
      onChange,
      onBlur,
      onFocus,
      placeholder,
      style,
      tip,
      ops,
      disabled = false,
      isMultiline = false,
      bodyStyle,
      popDirection = 'right',
      disabeldTip,
      className
    } = this.props;
    let { data } = this.state;

    let editor: JSX.Element;
    if (type === 'textarea') {
      editor = isMultiline ? (
        <div className="search-box multi-search-box" style={{ borderWidth: '0px' }}>
          <div className="search-input-wrap">
            <textarea
              className="tc-15-input-text search-input"
              placeholder={placeholder}
              onChange={e => this.setState({ data: e.target.value })}
              onBlur={this.handleBlur.bind(this)}
              style={style}
              value={data}
            />
          </div>
        </div>
      ) : (
        <textarea
          className="tc-15-input-textarea"
          value={data}
          placeholder={placeholder}
          onChange={e => this.setState({ data: e.target.value })}
          onBlur={this.handleBlur.bind(this)}
          style={style}
          disabled={disabled}
        />
      );
    } else {
      editor = (
        <input
          type={type || 'text'}
          className={className || 'tc-15-input-text m'}
          placeholder={placeholder}
          value={data}
          onChange={e => this.setState({ data: e.target.value })}
          onBlur={this.handleBlur.bind(this)}
          autoComplete="new-password"
          onFocus={onFocus}
          style={style}
          disabled={disabled}
        />
      );
    }

    let error = (
      <p className={tipMode === 'popup' ? 'text' : 'form-input-help'} style={{ fontSize: '12px' }}>
        {validator && (validator.status === 2 || validator.status === 3) && validator.message}
      </p>
    );

    let bubbleContent: string | React.ReactNode = '';
    if (validator && validator.status === 2) {
      bubbleContent = error;
    }
    if (disabled && disabeldTip) {
      bubbleContent = disabeldTip;
    }
    return (
      <div
        className={classnames('form-unit', { 'is-error': validator && validator.status === 2 })}
        style={bodyStyle ? bodyStyle : { display: 'inline-block', fontSize: '12px' }}
      >
        {tipMode === 'popup' ? (
          <Bubble placement={popDirection} content={bubbleContent || null}>
            {editor}
          </Bubble>
        ) : (
          editor
        )}
        {ops ? <span className="inline-help-text">{ops}</span> : <noscript />}
        {tip ? (
          <p className="form-input-help text-weak" style={{ fontSize: '12px' }}>
            {tip}
          </p>
        ) : (
          <noscript />
        )}
        {tipMode === 'popup' ? <noscript /> : error}
      </div>
    );
  }

  private handleBlur() {
    let { onBlur, onChange, rule } = this.props,
      { data } = this.state;

    onChange && onChange(data);

    if (rule) {
      let va = Validate(data, rule);
      onBlur && onBlur(va);
    } else {
      onBlur && onBlur(data);
    }
  }
}
