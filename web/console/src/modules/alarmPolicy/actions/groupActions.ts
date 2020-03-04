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
