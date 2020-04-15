import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../actions';
import { router } from '../../router';
import { UserCreate } from './create/UserCreate';
import { UserList } from './list/UserList';

import { RootState } from '../../models';
export interface RootProps extends RootState {
  actions?: typeof allActions;
}
export const UserPanel = (props) => {
  const state = useSelector((state) => state);
  // const dispatch = useDispatch();
  // const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { route } = state;
  const { action } = router.resolve(route);

  let content;
  if (!action) {
    content = <UserList />;
  } else if (action === 'create') {
    content = <UserCreate />;
  }

  return content;
};
