import { extend } from '@tencent/ff-redux';
import { RootState, UserInfo } from '../../models';
import * as ActionTypes from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;

const restActions = {
  fetchUserInfo: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let response = await WebAPI.fetchUserInfo();
      let info: UserInfo = response;
      dispatch({
        type: ActionTypes.UpdateUserInfo,
        payload: info
      });
    };
  }
};
export const detailActions = extend({}, restActions);
