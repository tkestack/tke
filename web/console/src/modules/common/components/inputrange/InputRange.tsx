import * as classnames from 'classnames';
import * as React from 'react';

import { Bubble } from '@tea/component';
import { BaseReactProps } from '@tencent/ff-redux';

import { Rule } from '../../../../../helpers/Validator';
import { Validation } from '../../models';

export interface InputRangeProps extends BaseReactProps {
  /**最小值 */
  minValue?: number;

  /**最大值 */
  maxValue?: number;

  /**校验态 */
  minValidator?: Validation;
  maxValidator?: Validation;

  /**输入事件 */
  onMinInput?: (value) => void;

  /**输入事件 */
  onMaxInput?: (value) => void;

  /**失去焦点事件 */
  onMinBlur?: (value) => void;

  /**失去焦点事件 */
  onMaxBlur?: (value) => void;

  /**placeholder */
  minPlaceholder?: string;
  maxPlaceholder?: string;

  /**输入框后面放置的操作 */
  ops?: string | JSX.Element;

  /**校验规则 */
  rule?: Rule;
}

interface InputRangeState {
  /**最小值 */
  min?: number;

  /**最大值 */
  max?: number;
}

export class InputRange extends React.Component<InputRangeProps, InputRangeState> {
  constructor(props, context) {
    super(props, context);

    this.state = {
      min: props.minValue,
      max: props.maxValue
    };
  }

  render() {
    let {
        minValidator,
        maxValidator,
        onMinInput,
        onMaxInput,
        onMinBlur,
        onMaxBlur,
        style,
        ops,
        minPlaceholder,
        maxPlaceholder
      } = this.props,
      { min, max } = this.state;

    let isError = (minValidator && minValidator.status === 2) || (maxValidator && maxValidator.status === 2);

    return (
      <div className="form-unit" style={{ fontSize: '12px' }}>
        <span className={classnames({ 'is-error': minValidator && minValidator.status === 2 })}>
          <Bubble placement="right" content={minValidator && minValidator.status === 2 ? minValidator.message : null}>
            <input
              type={'text'}
              className="tc-15-input-text m"
              style={{ width: '120px' }}
              placeholder={minPlaceholder}
              value={min.toString()}
              onChange={e => this.setState({ min: e.target.value })}
              onBlur={this.handleMinBlur.bind(this)}
            />
          </Bubble>
        </span>
        <span className="text"> ~ </span>
        <span className={classnames({ 'is-error': maxValidator && maxValidator.status === 2 })}>
          <Bubble placement="right" content={maxValidator && maxValidator.status === 2 ? maxValidator.message : null}>
            <input
              type={'text'}
              className="tc-15-input-text m"
              style={{ width: '120px' }}
              placeholder={maxPlaceholder}
              value={max.toString()}
              onChange={e => this.setState({ max: e.target.value })}
              onBlur={this.handleMaxBlur.bind(this)}
            />
          </Bubble>
        </span>
        {ops ? <span className="inline-help-text">{ops}</span> : <noscript />}
      </div>
    );
  }

  private handleMinBlur() {
    let { onMinBlur, onMinInput, rule } = this.props,
      { min } = this.state;

    onMinInput && onMinInput(min);
  }

  private handleMaxBlur() {
    let { onMinBlur, onMaxInput, rule } = this.props,
      { max } = this.state;

    onMaxInput && onMaxInput(max);
  }
}
