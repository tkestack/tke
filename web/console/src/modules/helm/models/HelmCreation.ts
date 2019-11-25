import { RecordSet } from '@tencent/qcloud-lib';
import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { QueryState } from '@tencent/qcloud-redux-query';
import { Region, RegionFilter, Resource, ResourceFilter } from '../../common/models';
import { TencenthubNamespace, TencenthubChart, TencenthubChartVersion, TencenthubChartReadMe } from './';
import { ListModel } from '@tencent/redux-list';

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
  region?: ListModel<Region, RegionFilter>;

  /** 集群列表 */
  cluster?: ListModel<Resource, ResourceFilter>;

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
