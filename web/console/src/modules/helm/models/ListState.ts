import { RecordSet } from '@tencent/qcloud-lib';
import { Region, RegionFilter, Resource, ResourceFilter } from '../../common/models';
import { ClusterHelmStatus, Helm, HelmFilter, InstallingHelm, InstallingHelmDetail, TencenthubChartVersion } from './';
import { HelmKeyValue } from './HelmCreation';
import { FFListModel, FetcherState, QueryState } from '@tencent/ff-redux';

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
