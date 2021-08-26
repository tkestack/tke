/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { FetcherState, FFListModel, QueryState, RecordSet } from '@tencent/ff-redux';

import { Region, RegionFilter, Resource, ResourceFilter } from '../../common/models';
import {
    ClusterHelmStatus, Helm, HelmFilter, InstallingHelm, InstallingHelmDetail,
    TencenthubChartVersion
} from './';
import { HelmKeyValue } from './HelmCreation';

export interface HelmListUpdateValid {
  otherChartUrl?: string;
  otherUserName?: string;
  otherPassword?: string;
}

export interface ListState {
  region?: FFListModel<Region, RegionFilter>;

  /** 集群列表 */
  cluster?: FFListModel<Resource, ResourceFilter>;

  clusterHelmStatus?: ClusterHelmStatus;

  /** 集群列表 */
  helmList?: FetcherState<RecordSet<Helm>>;

  /** 集群查询 */
  helmQuery?: QueryState<HelmFilter>;

  /** 集群选择 */
  helmSelection?: Helm;

  token?: string;

  installingHelmList?: FetcherState<RecordSet<InstallingHelm>>;
  installingHelmSelection?: InstallingHelm;
  installingHelmDetail?: InstallingHelmDetail;

  tencenthubChartVersionList?: FetcherState<RecordSet<TencenthubChartVersion>>;
  tencenthubChartVersionSelection?: TencenthubChartVersion;

  otherChartUrl?: string;
  otherTypeSelection?: string;
  otherUserName?: string;
  otherPassword?: string;

  isValid?: HelmListUpdateValid;

  kvs?: HelmKeyValue[];
}
