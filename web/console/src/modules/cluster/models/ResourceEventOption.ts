import { ResourceFilter, Resource } from './ResourceOption';
import { RecordSet } from '@tencent/qcloud-lib';
import { Event, EventFilter } from './Event';
import { QueryState, FetcherState } from '@tencent/ff-redux';

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
