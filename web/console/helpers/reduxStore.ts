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

/**
 * 重置redux store，用于离开页面时清空状态
 */
export const ResetStoreAction = 'ResetStore';

/**
 * 生成可重置的reducer，用于rootReducer简单包装
 * @return 可重置的reducer，当接收到 ResetStoreAction 时重置之
 */
export const generateResetableReducer: (rootReducer: Redux.Reducer) => Redux.Reducer = rootReducer => {
  return (state, action) => {
    let newState = state;
    // 销毁页面
    if (action.type === ResetStoreAction) {
      newState = undefined;
    }
    return rootReducer(newState, action);
  };
};
