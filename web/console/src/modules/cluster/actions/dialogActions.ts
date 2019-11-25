import { RootState, DialogNameEnum } from '../models';
import { cloneDeep } from '../../common';
import * as ActionType from '../constants/ActionType';

type GetState = () => RootState;

export const dialogActions = {
  updateDialogState: (dialogName: DialogNameEnum) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { dialogState } = getState();
      let newDialogState = cloneDeep(dialogState);

      newDialogState[dialogName] = !newDialogState[dialogName];
      dispatch({
        type: ActionType.UpdateDialogState,
        payload: newDialogState
      });
    };
  }
};
