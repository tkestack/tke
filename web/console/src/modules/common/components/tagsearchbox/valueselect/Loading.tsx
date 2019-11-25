import * as React from 'react';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export const Loading = props => {
  return (
    <div className="tc-15-autocomplete" style={{ left: `${props.offset}px` }}>
      <ul className="tc-15-autocomplete-menu" role="menu">
        <li role="presentation">
          <a className="autocomplete-empty" role="menuitem" href="javascript:;">
            {t('加载中 ..')}
          </a>
        </li>
      </ul>
    </div>
  );
};
