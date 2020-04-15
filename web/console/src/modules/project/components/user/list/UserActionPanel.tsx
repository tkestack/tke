import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Justify, Icon, Table, Button, SearchBox } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';

export const UserActionPanel = (props) => {
  const state = useSelector((state) => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { route, userList } = state;
  console.log('route:', route, route.queries);
  return (
    <Table.ActionPanel>
      <Justify
        left={
          <Button
            type="primary"
            onClick={(e) => {
              e.preventDefault();
              router.navigate({ sub: 'detail', tab: 'member', action: 'create' }, route.queries);
            }}
          >
            {t('新建')}
          </Button>
        }
        right={
          <React.Fragment>
            <SearchBox
              value={userList.query.keyword || ''}
              onChange={actions.user.changeKeyword}
              onSearch={actions.user.performSearch}
              onClear={() => {
                actions.user.performSearch('');
              }}
              placeholder={t('请输入用户名称')}
            />
          </React.Fragment>
        }
      />
    </Table.ActionPanel>
  );
};
