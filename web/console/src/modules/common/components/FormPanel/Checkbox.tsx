import * as React from 'react';
import { Checkbox, CheckboxProps } from '@tea/component';

interface FormPanelCheckboxProps extends CheckboxProps {
  text?: String;
}

function FormPanelCheckbox({ ...props }: FormPanelCheckboxProps) {
  return <Checkbox {...props}>{props.text}</Checkbox>;
}

export { FormPanelCheckbox, FormPanelCheckboxProps };
