import { ReduxAction } from '../';

export function reduceToPayload<T>(actionType: string | number, initialState: T) {
  return (state: T = initialState, action: ReduxAction<T>) => {
    if (action.type === actionType) {
      return action.payload;
    }
    return state;
  };
}
