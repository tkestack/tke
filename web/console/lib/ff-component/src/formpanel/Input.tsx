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

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Input, InputAdorment, InputAdornmentProps, InputProps } from '@tencent/tea-component';

import { FormPanelValidatable, FormPanelValidatablePropsWhiteoutChildren } from './Validatable';

interface FormPanelInputProps extends InputProps, FormPanelValidatablePropsWhiteoutChildren {
  inputAdornment?: InputAdornmentProps;
  label?: string;
}

function FormPanelInput({
  validator,
  formvalidator,
  vkey,
  vactions,
  errorTipsStyle,
  bubblePlacement,

  onBlur,

  inputAdornment,
  ...props
}: FormPanelInputProps) {
  let rOnBlur = onBlur;
  if (props.label && !props.placeholder) {
    props.placeholder = t('请输入') + props.label;
  }

  let validatableProps = {
    validator,
    formvalidator,
    vkey,
    vactions,
    errorTipsStyle,
    bubblePlacement
  };

  let onBlurWrap =
    vactions && vkey
      ? event => {
          rOnBlur && rOnBlur(event);
          vactions && vkey && vactions.validate(vkey);
        }
      : rOnBlur;

  if (inputAdornment) {
    //添加一个div.style=inline-block，为了外面包裹bubble时能正常工作
    return (
      <FormPanelValidatable {...validatableProps}>
        <div style={{ display: 'inline-block' }}>
          <InputAdorment {...inputAdornment}>
            <Input {...props} onBlur={onBlurWrap} />
          </InputAdorment>
        </div>
      </FormPanelValidatable>
    );
  } else {
    return (
      <FormPanelValidatable {...validatableProps}>
        <Input {...props} onBlur={onBlurWrap} />
      </FormPanelValidatable>
    );
  }
}

export { FormPanelInput, FormPanelInputProps };
