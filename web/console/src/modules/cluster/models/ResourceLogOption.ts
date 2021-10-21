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
