import * as React from 'react';
import { Justify } from '@tea/component';

interface FormPanelFooterProps {
  children: React.ReactNode;
}
function FormPanelFooter({ children }: FormPanelFooterProps) {
  return (
    <li className="pure-text-row fixed formpanel-footer" style={{ marginLeft: -20, paddingRight: 26 }}>
      <Justify left={children} />
    </li>
  );
}

export { FormPanelFooter, FormPanelFooterProps };
