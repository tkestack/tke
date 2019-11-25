import * as React from 'react';
import { FormText, Icon } from '@tea/component';
import { StyledProps } from '@tencent/tea-component/lib/_type';

interface FormPanelTextProps extends StyledProps {
  children?: React.ReactNode;
  onEdit?: () => void;
}

function FormPanelText({ children, ...props }: FormPanelTextProps) {
  return (
    <FormText {...props}>
      {children}
      {props.onEdit && <Icon onClick={() => props.onEdit()} style={{ cursor: 'pointer' }} type="pencil" />}
    </FormText>
  );
}

export { FormPanelText, FormPanelTextProps };
