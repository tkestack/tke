import { regionActions } from './regionActions';
import { RootState } from '../models';
import * as ActionType from '../constants/ActionType';

type GetState = () => RootState;

/**
 * 切换模式
 */
export const modeActions = {
  /**
   * 切换模式
   */
  changeMode: mode => {
    return (dispatch, getState) => {
      let { route } = getState();
      dispatch({
        type: ActionType.ChangeMode,
        payload: mode
      });

      dispatch(regionActions.applyFilter({}));
    };
  }
};
