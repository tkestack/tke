import { FetchState } from './FetchState';
import { RecordSet } from './RecordSet';

/** state for data fetcher */
export interface FetcherState<TData> {
  /**
   * current fetch state
   * */
  fetchState: FetchState;

  /**
   * 请求是否已完成
   */
  fetched?: boolean;

  /**
   * data fetched from the last time
   * */
  data?: TData;

  /**
   * error object when in fail state
   */
  error?: any;

  /**
   * If the fetch started for a while, the loading will be true.
   * You can specific the duration by passing `loadingTolerance` when generating action creator.
   * If the duration is not specific, loading will be true as well as the fetchState gets to `Fetching`
   * */
  loading?: boolean;

  pages?: FetcherState<TData>[];
}
