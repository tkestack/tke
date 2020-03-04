import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { Resource } from '../../common/models/Resource';
import * as ActionType from '../constants/ActionType';
import { Namespace } from '../models/Namespace';
import { router } from '../router';
import { DetailReducer } from './DetailReducer';
import { HelmCreationReducer } from './HelmCreationReducer';
import { ListReducer } from './ListReducer';

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
