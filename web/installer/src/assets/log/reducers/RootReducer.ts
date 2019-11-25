import { router } from './../router';
import { generateWorkflowReducer } from '@tencent/qcloud-redux-workflow';
import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import * as ActionType from '../constants/ActionType';
import { Log } from '../../common/models';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  logList: generateFetcherReducer<RecordSet<Log>>({
    actionType: ActionType.FetchLogList,
    initialData: {
      recordCount: 0,
      records: [] as Log[]
    }
  }),

  logQuery: generateQueryReducer({
    actionType: ActionType.QueryLogList,
    initialState: {
      paging: {
        pageIndex: 0,
        pageSize: 100
      }
    }
  })
});
