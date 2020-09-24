/**
 * hpa这个scope内的全局数据
 */
import React, { useReducer } from 'react';
import { useImmerReducer } from 'use-immer';

/**
 * 常量定义
 */
const CHANGE_NAMESPACE = 'CHANGE_NAMESPACE';
const PANEL = 'PANEL';
// const SET_HPADATA = 'SET_HPADATA';

/**
 * 初始值
 */
const initialState = {
  // 选中的namespace value
  namespaceValue: '',
  panel: '',
  // HPAData: {}
};

/**
 * 全局数据定义
 */
const StateContext = React.createContext(null);
const DispatchContext = React.createContext(null);

/**
 * 结合immer实现的reducer
 * @param draft
 * @param action
 */
function immerReducer(draft, action) {
  switch (action.type) {
    case CHANGE_NAMESPACE:
      draft.namespaceValue = action.payload.namespaceValue;
      break;
    // case SET_HPADATA:
    //   draft.HPAData = action.payload.HPAData;
    //   break;
  }
  return draft;
}

/**
 * 包裹函数
 * @param props
 * @constructor
 */
function HpaScopeProvider(props) {
  const [hpaScopeState, hpaScopeDispatch] = useImmerReducer(immerReducer, initialState);
  return (
    <DispatchContext.Provider value={hpaScopeDispatch}>
      <StateContext.Provider value={hpaScopeState}>
        {props.children}
      </StateContext.Provider>
    </DispatchContext.Provider>
  );
}

export {
  CHANGE_NAMESPACE,
  // SET_HPADATA,
  StateContext,
  DispatchContext,
  HpaScopeProvider
};
