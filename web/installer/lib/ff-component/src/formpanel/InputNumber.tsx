import * as React from 'react';

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import {
    Input, InputAdorment, InputAdornmentProps, InputNumber, InputNumberProps, InputProps
} from '@tencent/tea-component';

import { FormPanelValidatable, FormPanelValidatablePropsWhiteoutChildren } from './Validatable';

interface FormPanelInputNumberProps extends InputNumberProps, FormPanelValidatablePropsWhiteoutChildren {
  inputAdornment?: InputAdornmentProps;
  label?: string;
}

function FormPanelInputNumber({
  validator,
  formvalidator,
  vkey,
  vactions,
  errorTipsStyle,
  bubblePlacement,

  onChange,

  inputAdornment,
  ...props
}: FormPanelInputNumberProps) {
  let rOnChange = onChange;

  let validatableProps = {
    validator,
    formvalidator,
    vkey,
    vactions,
    errorTipsStyle,
    bubblePlacement
  };

  let onChangeWrap =
    vactions && vkey
      ? (value, context) => {
          rOnChange && rOnChange(value, context);
          vactions && vkey && vactions.validate(vkey);
        }
      : rOnChange;
  if (inputAdornment) {
    //添加一个div.style=inline-block，为了外面包裹bubble时能正常工作
    return (
      <FormPanelValidatable {...validatableProps}>
        <div style={{ display: 'inline-block' }}>
          <InputAdorment {...inputAdornment}>
            <InputNumber {...props} onChange={onChangeWrap} />
          </InputAdorment>
        </div>
      </FormPanelValidatable>
    );
  } else {
    return (
      <FormPanelValidatable {...validatableProps}>
        <InputNumber {...props} onChange={onChangeWrap} />
      </FormPanelValidatable>
    );
  }
}

export { FormPanelInputNumber, FormPanelInputNumberProps };
