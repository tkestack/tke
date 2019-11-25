
export type Reducer<TState, TAction> = (state: TState, action: TAction) => TState;

export function serialReducer<TState, TAction>(
    first: Reducer<TState, TAction>,
    ...series: Reducer<TState, TAction>[]): Reducer<TState, TAction> {

    const reducers = [first, ...series];

    const combined: Reducer<TState, TAction> = (state, action) =>
        reducers.reduce((state, reducer) => reducer(state, action), state);

    return combined;
}