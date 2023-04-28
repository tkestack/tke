import { apiVersion, ApiVersionKeyName } from './apiVersion';
export { ApiVersionKeyName } from './apiVersion';

export function apiPath(clusterVersion: string, resourceType: ApiVersionKeyName) {
  const versionInfo = apiVersion(clusterVersion);

  const { basicEntry, group, version } = versionInfo[resourceType];
  return `/${basicEntry}/${group}/${version}`;
}
