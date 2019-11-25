import { RecordSet } from '@tencent/qcloud-lib';
import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { QueryState } from '@tencent/qcloud-redux-query';
import { WorkflowState } from '@tencent/qcloud-redux-workflow';
import { Helm, HelmHistory, HelmHistoryFilter } from './';

type HelmWorkflow = WorkflowState<Helm, string>;
export interface DetailState {
  /**当前的 Cluster详细数据 */
  helm?: Helm;
  isRefresh?: boolean;
  historyQuery?: QueryState<HelmHistoryFilter>;
  histories?: FetcherState<RecordSet<HelmHistory>>;
}
