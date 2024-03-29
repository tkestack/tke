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
export enum FetchState {
  /** indicates the data is up to date and ready to use */
  Ready = 'Ready',

  /** indicates the data is out of date, and the new data is fetching */
  Fetching = 'Fetching',

  /**
   * indicates the data is out of date, and the new data fetches failed
   */
  Failed = 'Failed'
}
