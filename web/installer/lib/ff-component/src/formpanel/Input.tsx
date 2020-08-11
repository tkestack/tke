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
