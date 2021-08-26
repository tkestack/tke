/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../actions';
import { router } from '../../router';
import { UserCreate } from './create/UserCreate';
import { UserList } from './list/UserList';
import { UserDetail } from './detail/UserDetail';

export const UserPanel = (props) => {
  const state = useSelector((state) => state);
  // const dispatch = useDispatch();
  // const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { route } = state;
  const { sub, action } = router.resolve(route);

  let content;
  if (!action && (!sub || sub === 'normal')) {
    content = <UserList />;
  } else if (action === 'create') {
    content = <UserCreate />;
  } else if (action === 'detail') {
    content = <UserDetail />;
  }
  return content;
};
