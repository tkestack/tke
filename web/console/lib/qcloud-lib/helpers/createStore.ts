import { createStore as createReduxStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import { createLogger } from 'redux-logger';

export const createStore =
  process.env.NODE_ENV === 'development'
    ? applyMiddleware(thunk, createLogger({ collapsed: true }))(createReduxStore)
    : applyMiddleware(thunk)(createReduxStore);
