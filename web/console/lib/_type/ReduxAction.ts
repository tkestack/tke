export interface ReduxAction<TPayload> {
  /**
   * The action type, the `number` type is to support enum.
   * */
  type: string | number;
  payload?: TPayload;
  error?: boolean;
  meta?: any;
}
