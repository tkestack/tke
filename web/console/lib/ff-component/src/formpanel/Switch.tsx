import * as React from 'react';

import { Switch, SwitchProps } from '@tencent/tea-component';

interface FormPanelSwitchProps extends SwitchProps {
  text?: String;
}

function FormPanelSwitch({ ...props }: FormPanelSwitchProps) {
  return <Switch {...props}>{props.text}</Switch>;
}

export { FormPanelSwitch, FormPanelSwitchProps };
