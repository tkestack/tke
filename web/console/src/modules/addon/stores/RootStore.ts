/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { createStore } from '@tencent/ff-redux';

import { generateResetableReducer } from '../../../../helpers';
import { RootReducer } from '../reducers/RootReducer';

export function configStore() {
  const store = createStore(generateResetableReducer(RootReducer));

  // hot reloading
  if (typeof module !== 'undefined' && (module as any).hot) {
    (module as any).hot.accept('../reducers/RootReducer', () => {
      store.replaceReducer(generateResetableReducer(require('../reducers/RootReducer').RootReducer));
    });
  }

  return store;
}
