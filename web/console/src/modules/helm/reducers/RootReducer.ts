import { HelmCreationReducer } from './HelmCreationReducer';
import { DetailReducer } from './DetailReducer';
import { ListReducer } from './ListReducer';
import { combineReducers } from 'redux';
import { router } from '../router';

export const RootReducer = combineReducers({
  route: router.getReducer(),
  helmCreation: HelmCreationReducer,
  listState: ListReducer,
  detailState: DetailReducer
});
