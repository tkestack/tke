import * as React from 'react';

import { Checkbox, CheckboxProps } from '@tencent/tea-component';

interface FormPanelCheckboxProps extends CheckboxProps {
  text?: React.ReactNode;
}

function FormPanelCheckbox({ ...props }: FormPanelCheckboxProps) {
  return <Checkbox {...props}>{props.text}</Checkbox>;
}

export { FormPanelCheckbox, FormPanelCheckboxProps };
