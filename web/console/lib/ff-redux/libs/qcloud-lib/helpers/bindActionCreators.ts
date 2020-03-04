﻿/**
 * like redux.bindActionCreators but do it recurisivly
 * */
export function bindActionCreators<T>(actions: T, dispatch: Redux.Dispatch): T {
  if (typeof actions !== 'object' || !actions) {
    throw new RangeError('invalid actions!');
  }

  let result: any = {};

  for (let key in actions) {
    if (!actions.hasOwnProperty(key)) {
      continue;
    }
    const creator = actions[key];
    if (typeof creator === 'object' && creator) {
      result[key] = bindActionCreators(creator, dispatch);
    } else if (typeof creator === 'function') {
      result[key] = (...args: any[]) => dispatch(creator(...args));
    }
  }

  return result as T;
}
