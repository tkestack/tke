import { FetcherState, QueryState, RecordSet } from '@tencent/ff-redux';

import { Pod } from './Pod';
import { PodLogFilter } from './ResourceDetailState';
import { Resource, ResourceFilter } from './ResourceOption';

export interface ResourceLogOption {
  /** workloadType */
  workloadType?: string;

  /** workloadSelection */
  workloadSelection?: string;

  /** namespaceSelection */
  namespaceSelection?: string;

  /** workloadquery */
  workloadQuery?: QueryState<ResourceFilter>;

  /** workload的列表 */
  workloadList?: FetcherState<RecordSet<Resource>>;

  /** pod的查询 */
  podQuery?: QueryState<ResourceFilter>;

  /** pod的列表 */
  podList?: FetcherState<RecordSet<Pod>>;

  /** podSelection */
  podSelection?: string;

  /** container */
  containerSelection?: string;

  /** log的查询 */
  logQuery?: QueryState<PodLogFilter>;

  /** log的列表 */
  logList?: FetcherState<RecordSet<string>>;

  /** tailLines */
  tailLines?: string;

  /** 是否开启自动刷新 */
  isAutoRenew?: boolean;
}
