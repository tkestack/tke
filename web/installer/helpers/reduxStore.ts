import { Reducer } from 'redux';

/**
 * 重置redux store，用于离开页面时清空状态
 */
export const ResetStoreAction = 'ResetStore';

/**
 * 生成可重置的reducer，用于rootReducer简单包装
 * @return 可重置的reducer，当接收到 ResetStoreAction 时重置之
 */
export const generateResetableReducer: (rootReducer: Reducer) => Reducer = rootReducer => {
  return (state, action) => {
    let newState = state;
    // 销毁页面
    if (action.type === ResetStoreAction) {
      newState = undefined;
    }
    return rootReducer(state, action);
  };
};
