import { RecordSet } from '@tencent/qcloud-lib';
import { WorkflowState } from '@tencent/ff-redux';
import { Helm, HelmHistory, HelmHistoryFilter } from './';
import { FetcherState, QueryState } from '@tencent/ff-redux';

type HelmWorkflow = WorkflowState<Helm, string>;
export interface DetailState {
  /**当前的 Cluster详细数据 */
  helm?: Helm;
  isRefresh?: boolean;
  historyQuery?: QueryState<HelmHistoryFilter>;
  histories?: FetcherState<RecordSet<HelmHistory>>;
}
