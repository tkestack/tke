/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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
