import { applyMiddleware, createStore as createReduxStore } from 'redux';
import { createLogger } from 'redux-logger';
import thunk from 'redux-thunk';

export const createStore =
  process.env.NODE_ENV === 'development'
    ? applyMiddleware(thunk, createLogger({ collapsed: true }))(createReduxStore)
    : applyMiddleware(thunk)(createReduxStore);
