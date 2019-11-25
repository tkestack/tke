import { RecordSet } from '@tencent/qcloud-lib';
import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { QueryState } from '@tencent/qcloud-redux-query';
import { Region, RegionFilter, Resource, ResourceFilter } from '../../common/models';
import { ClusterHelmStatus, Helm, HelmFilter, InstallingHelm, InstallingHelmDetail, TencenthubChartVersion } from './';
import { HelmKeyValue } from './HelmCreation';
import { ListModel } from '@tencent/redux-list';

export interface HelmListUpdateValid {
  otherChartUrl?: string;
  otherUserName?: string;
  otherPassword?: string;
}

export interface ListState {
  region?: ListModel<Region, RegionFilter>;

  /** 集群列表 */
  cluster?: ListModel<Resource, ResourceFilter>;

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
