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
