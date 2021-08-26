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

import { Placement } from 'popper.js';
import * as React from 'react';

import { insertCSS } from '@tencent/ff-redux';
import { Validation, ValidationIns, ValidatorModel } from '@tencent/ff-validator';
import { Bubble, Icon } from '@tencent/tea-component';

import { classNames } from '../lib/classname';
import { FormPanel } from './FormPanel';

insertCSS(
  '@tencent/ff-component/ff-component-validatable',
  `
.ff-component-validatable {
  display : inline-block;
}
`
);

interface FormPanelValidatableProps {
  //
  formvalidator?: ValidatorModel;
  vkey?: string;
  vactions?: ValidationIns;

  validator?: Validation;

  errorTipsStyle?: 'Icon' | 'Bubble' | 'Message';
  bubblePlacement?: Placement;

  children?: React.ReactNode;
}

interface FormPanelValidatablePropsWhiteoutChildren {
  //
  formvalidator?: ValidatorModel;
  vkey?: string;
  vactions?: ValidationIns;

  //如果只
  validator?: Validation;

  errorTipsStyle?: 'Icon' | 'Bubble' | 'Message';
  bubblePlacement?: Placement;
}

function FormPanelValidatable({
  formvalidator,
  vkey,
  // vactions,

  validator = formvalidator && vkey ? (formvalidator[vkey] as Validation) : null,

  errorTipsStyle = 'Icon',
  bubblePlacement = 'right',

  children
}: FormPanelValidatableProps) {
  if (!validator && !vkey) {
    return <React.Fragment>{children}</React.Fragment>;
  }
  let iserror = validator ? validator.status > 1 : false;
  // if (!iserror) {
  //   return <React.Fragment>{children}</React.Fragment>;
  // } else {
  return (
    <div className={classNames('ff-component-validatable', { 'is-error': iserror })}>
      {errorTipsStyle === 'Icon' && (
        <React.Fragment>
          {children}
          {iserror && (
            <Bubble placement={bubblePlacement} content={validator ? validator.message || null : null}>
              <Icon type="error" style={{ marginLeft: 5 }} />
            </Bubble>
          )}
        </React.Fragment>
      )}
      {errorTipsStyle === 'Bubble' && (
        <React.Fragment>
          <Bubble placement={bubblePlacement} content={validator ? validator.message || null : null}>
            {children}
          </Bubble>
        </React.Fragment>
      )}
      {errorTipsStyle === 'Message' && (
        <React.Fragment>
          {children}
          <FormPanel.HelpText theme="danger" parent="div">
            {validator ? validator.message || '' : ''}
          </FormPanel.HelpText>
        </React.Fragment>
      )}
    </div>
  );
}

export { FormPanelValidatable, FormPanelValidatableProps, FormPanelValidatablePropsWhiteoutChildren };
