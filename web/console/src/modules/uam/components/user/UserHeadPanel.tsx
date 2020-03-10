import * as React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { Justify, Icon } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../router';

export const UserHeadPanel = () => {
  const state = useSelector(state => state);
  const { route } = state;
  let urlParam = router.resolve(route);
  const { sub } = urlParam;
  return (
    <Justify
      left={
        <h2>
          {sub ? (
            <React.Fragment>
              <a href="javascript:history.go(-1);">
                <Icon type="btnback" />
              </a>
              <span style={{ marginLeft: '10px' }}>{sub}</span>
            </React.Fragment>
          ) : (
            <Trans>用户管理</Trans>
          )}
        </h2>
      }
    />
  );
};
