import { FetcherTrigger } from './FetcherTrigger';
import { ReduxAction } from './ReduxAction';

/** action payload for trigger */
export interface FetcherPayload<TData> {
  trigger: FetcherTrigger;
  data?: TData;
  error?: Error;
  append?: boolean;
  pageIndex?: number;
  clear?: boolean;
}

export type FetcherAction<TData> = ReduxAction<FetcherPayload<TData>>;
