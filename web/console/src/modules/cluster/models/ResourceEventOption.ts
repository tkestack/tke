import { QueryState } from '@tencent/qcloud-redux-query';
import { ResourceFilter, Resource } from './ResourceOption';
import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { RecordSet } from '@tencent/qcloud-lib';
import { Event, EventFilter } from './Event';

export interface ResourceEventOption {
  /** workloadType */
  workloadType?: string;

  /** namespaceSelection */
  namespaceSelection?: string;

  /** workloadquery */
  workloadQuery?: QueryState<ResourceFilter>;

  /** workload的列表 */
  workloadList?: FetcherState<RecordSet<Resource>>;

  /** workloadSelection */
  workloadSelection?: string;

  /** eventQuery */
  eventQuery?: QueryState<EventFilter>;

  /** eventList */
  eventList?: FetcherState<RecordSet<Event>>;

  /** 是否自动刷新 */
  isAutoRenew?: boolean;
}
