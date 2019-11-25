import * as React from 'react';
import { Input, InputProps, InputAdorment, InputAdornmentProps } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

interface FormPanelInputProps extends InputProps {
  inputAdornment?: InputAdornmentProps;
  label?: string;
}

function FormPanelInput({ ...props }: FormPanelInputProps) {
  if (props.label && !props.placeholder) {
    props.placeholder = t('请输入') + props.label;
  }
  if (props.inputAdornment) {
    return (
      <InputAdorment {...props.inputAdornment}>
        <Input {...props} />
      </InputAdorment>
    );
  } else {
    return <Input {...props} />;
  }
}

export { FormPanelInput, FormPanelInputProps };
