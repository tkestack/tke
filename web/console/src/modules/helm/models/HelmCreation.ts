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
