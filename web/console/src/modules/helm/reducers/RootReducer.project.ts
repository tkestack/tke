import { Namespace } from './../models/Namespace';
import { Resource } from './../../common/models/Resource';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import { HelmCreationReducer } from './HelmCreationReducer';
import { DetailReducer } from './DetailReducer';
import { ListReducer } from './ListReducer';
import { combineReducers } from 'redux';
import * as ActionType from '../constants/ActionType';
import { router } from '../router';

export const RootReducer = combineReducers({
  route: router.getReducer(),
  helmCreation: HelmCreationReducer,
  listState: ListReducer,
  detailState: DetailReducer,

  namespaceQuery: generateQueryReducer({
    actionType: ActionType.QueryNamespaceList
  }),
  namespaceSelection: reduceToPayload(ActionType.SelectNamespace, ''),
  projectList: reduceToPayload(ActionType.InitProjectList, []),
  projectSelection: reduceToPayload(ActionType.ProjectSelection, ''),
  projectNamespaceQuery: generateQueryReducer({
    actionType: ActionType.QueryProjectNamespace
  }),
  projectNamespaceList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchProjectNamespace,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),
  namespaceList: generateFetcherReducer<RecordSet<Namespace>>({
    actionType: ActionType.FetchNamespaceList,
    initialData: {
      recordCount: 0,
      records: [] as Namespace[]
    }
  })
});
