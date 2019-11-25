import * as clusterApi from './ClusterApi';
import * as regionApi from './RegionApi';
import * as configApi from './ConfigApi';
import * as versionApi from './VersionApi';
import * as bandwidthApi from './BandwidthApi';

export const WebAPI = {
  region: regionApi,
  cluster: clusterApi,
  config: configApi,
  version: versionApi,
  bandwidth: bandwidthApi
};
