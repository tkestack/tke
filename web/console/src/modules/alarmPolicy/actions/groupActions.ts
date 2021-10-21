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
import { RootState } from '../models';

type GetState = () => RootState;

// /**
//  * 查询用户组
//  */
// const fetchGroupActions = generateFetcherActionCreator({
//   actionType: ActionType.FetchGroupList,
//   fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
//     const response = await WebAPI.fetchGroupList(getState().groupQuery);
//     /**初始化接受组 */
//     if (fetchOptions && fetchOptions.data) {
//       let groupSelection = [];
//       fetchOptions.data.forEach(item => {
//         let finder = response.records.find(group => group.groupId === item);
//         finder && groupSelection.push(finder);
//       });
//       dispatch({
//         type: ActionType.SelectGroup,
//         payload: groupSelection
//       });
//     }
//     return response;
//   }
// });

// /**
//  * 查询用户组
//  */
// const queryGroupActions = generateQueryActionCreator<GroupFilter>({
//   actionType: ActionType.QueryGroupList,
//   bindFetcher: fetchGroupActions
// });

/**
 * 选择用户组
 */
// const selectActions = {
//   selectGroup: group => {
//     return (dispatch, getState) => {
//       dispatch({
//         type: ActionType.SelectGroup,
//         payload: group
//       });
//     };
//   }
// };

// export const groupActions = selectActions;
