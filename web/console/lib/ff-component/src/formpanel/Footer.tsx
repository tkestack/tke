import * as React from 'react';

import { insertCSS } from '@tencent/ff-redux';
import { Justify } from '@tencent/tea-component';

insertCSS(
  '@tencent/ff-component/formpanel/footer',
  `
  .ff-formpanel-footer {    
    position: fixed;
    bottom: 0;
    z-index: 11;
    width: 1224px;
    padding: 15px 36px 15px 100px;
    border-top: 1px solid #dbe3e4;
    background-color: #fff;
    box-shadow: 0 2px 3px 0 rgba(0,0,0,.2);
  }

  .ff-formpanel-footer .app-tke-fe-btn {
    margin-right:20px !important;
  }
`
);

interface FormPanelFooterProps {
  children: React.ReactNode;
  width?: number;
  cardRef?: React.RefAttributes<HTMLDivElement>['ref'];
}
function FormPanelFooter({ children, width, cardRef }: FormPanelFooterProps) {
  let [marginLeft, setMarginLeft] = React.useState(-20);
  let [footWidth, setFootWidth] = React.useState(1224);

  let cardEl = cardRef || (React.createRef() as any);

  if (cardEl) {
    React.useEffect(() => {
      let refreshWidth = () => {
        if (cardEl && cardEl.current) {
          setFootWidth(cardEl.current.offsetWidth - 136);
        }
      };
      window.addEventListener('resize', refreshWidth);
      refreshWidth();
      return () => {
        window.removeEventListener('resize', refreshWidth);
      };
    });
  }

  React.useEffect(() => {
    let refreshMarginLeft = () => {
      setMarginLeft(-20 - window.scrollX);
    };
    window.addEventListener('scroll', refreshMarginLeft);
    refreshMarginLeft();
    return () => {
      window.removeEventListener('scroll', refreshMarginLeft);
    };
  });
  return (
    <div className="ff-formpanel-footer" style={{ width: width ? width : footWidth, marginLeft: marginLeft }}>
      <Justify left={children} />
    </div>
  );
}

export { FormPanelFooter, FormPanelFooterProps };
