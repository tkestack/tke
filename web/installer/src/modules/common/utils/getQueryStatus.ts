import { FetcherState, FetchState, RecordSet } from '@tencent/ff-redux';

export const getQueryStatus = (fetcher: FetcherState<RecordSet<any>>, search?: any) => {
  let status: any = null;
  if (fetcher.fetchState === FetchState.Fetching) {
    status = 'loading';
  } else if (search) {
    status = 'found';
  } else if (fetcher.fetched && !fetcher.data.recordCount) {
    status = 'empty';
  } else if (fetcher.error) {
    status = 'error';
  }

  return status;
};
