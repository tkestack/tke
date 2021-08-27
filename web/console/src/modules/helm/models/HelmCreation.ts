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
import { FetcherState, FFListModel, RecordSet } from '@tencent/ff-redux';

import { Region, RegionFilter, Resource, ResourceFilter } from '../../common/models';
import {
    TencenthubChart, TencenthubChartReadMe, TencenthubChartVersion, TencenthubNamespace
} from './';

export interface HelmKeyValue {
  key?: string;
  value?: string;
}

export interface HelmCreationValid {
  name?: string;
  otherChartUrl?: string;
  otherUserName?: string;
  otherPassword?: string;
}

export interface HelmCreation {
  region?: FFListModel<Region, RegionFilter>;

  /** 集群列表 */
  cluster?: FFListModel<Resource, ResourceFilter>;

  name?: string;

  isValid?: HelmCreationValid;

  resourceSelection?: string;

  token?: string;

  tencenthubTypeSelection?: string;
  tencenthubNamespaceList?: FetcherState<RecordSet<TencenthubNamespace>>;
  tencenthubNamespaceSelection?: string;
  tencenthubChartList?: FetcherState<RecordSet<TencenthubChart>>;
  tencenthubChartSelection?: TencenthubChart;
  tencenthubChartVersionList?: FetcherState<RecordSet<TencenthubChartVersion>>;
  tencenthubChartVersionSelection?: TencenthubChartVersion;
  tencenthubChartReadMe?: TencenthubChartReadMe;

  otherChartUrl?: string;
  otherTypeSelection?: string;
  otherUserName?: string;
  otherPassword?: string;

  kvs?: HelmKeyValue[];
}
