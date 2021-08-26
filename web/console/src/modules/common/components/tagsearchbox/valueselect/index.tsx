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
import { PureInput } from './PureInput';
import { SingleValueSelect } from './SingleValueSelect';
import { MultipleValueSelect } from './MultipleValueSelect';
import { Loading } from './Loading';

export interface ValueSelectProps {
  type: string;
  values?: Array<any> | Function;
  inputValue?: string;
  onChange?: (value: any) => void;
  onSelect?: (value: any) => void;
  onCancel?: () => void;
  offset: number;
}

export class ValueSelect extends React.Component<ValueSelectProps, any> {
  mount = false;

  constructor(props) {
    super(props);

    let values = [];
    const propsValues = this.props.values;

    if (typeof propsValues !== 'function') {
      values = propsValues;
    }

    this.state = { values };
  }

  componentDidMount() {
    this.mount = true;
    const propsValues = this.props.values;
    if (typeof propsValues === 'function') {
      const res = propsValues();
      // Promise
      if (res && res.then) {
        res.then(values => {
          this.mount && this.setState({ values });
        });
      } else {
        this.mount && this.setState({ values: res });
      }
    }
  }

  componentWillUnmount() {
    this.mount = false;
  }

  handleKeyDown = (keyCode: number): boolean => {
    if (this['select'] && this['select'].handleKeyDown) {
      return this['select'].handleKeyDown(keyCode);
    }
    return true;
  };

  // handleKeyUp = (keyCode: number): boolean => {
  //   if (this['select'] && this['select'].handleKeyUp) {
  //     return this['select'].handleKeyUp(keyCode);
  //   }
  //   return true;
  // }

  render() {
    const values = this.state.values;
    const { type, inputValue, onChange, onSelect, onCancel, offset } = this.props;

    const props = { values, inputValue, onChange, onSelect, onCancel, offset };

    switch (type) {
      case 'input':
        return (
          <PureInput
            ref={select => {
              this['select'] = select;
            }}
            {...props}
          />
        );
      case 'single':
        if (values.length <= 0) {
          return <Loading offset={offset} />;
        }
        return (
          <SingleValueSelect
            ref={select => {
              this['select'] = select;
            }}
            {...props}
          />
        );
      case 'multiple':
        if (values.length <= 0) {
          return <Loading />;
        }
        return (
          <MultipleValueSelect
            ref={select => {
              this['select'] = select;
            }}
            {...props}
          />
        );
    }
    return null;
  }
}
