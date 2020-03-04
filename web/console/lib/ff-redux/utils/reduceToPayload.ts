export interface ReduxAction<TPayload> {
  /**
   * The action type, the `number` type is to support enum.
   * */
  type: string | number;
  payload?: TPayload;
  error?: boolean;
  meta?: any;
}

export function reduceToPayload<T>(actionType: string | number, initialState: T) {
  return (state: T = initialState, action: ReduxAction<T>) => {
    if (action.type === actionType) {
      return action.payload;
    }
    return state;
  };
}
