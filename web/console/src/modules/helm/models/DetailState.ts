import { FetcherState, QueryState, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { Helm, HelmHistory, HelmHistoryFilter } from './';

type HelmWorkflow = WorkflowState<Helm, string>;
export interface DetailState {
  /**当前的 Cluster详细数据 */
  helm?: Helm;
  isRefresh?: boolean;
  historyQuery?: QueryState<HelmHistoryFilter>;
  histories?: FetcherState<RecordSet<HelmHistory>>;
}
