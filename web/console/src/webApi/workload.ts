import { apiPath, ApiVersionKeyName } from './apiPath';
import { Request } from './request';

interface IProps {
  clusterVersion: string;
  resourceId: string;
  namespace: string;
  clusterId: string;
  kind: ApiVersionKeyName;
  data?: any;
}

export function fetchWorkloadResource({
  resourceId,
  namespace,
  clusterId,
  kind,

  clusterVersion
}: IProps) {
  return Request.get<any, any>(`${apiPath(clusterVersion, kind)}/namespaces/${namespace}/${kind}s/${resourceId}`, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
}

export function updateWorkloadResource({ resourceId, namespace, clusterId, kind, data, clusterVersion }: IProps) {
  return Request.patch(`${apiPath(clusterVersion, kind)}/namespaces/${namespace}/${kind}s/${resourceId}`, data, {
    headers: {
      'Content-Type': 'application/strategic-merge-patch+json',
      'X-TKE-ClusterName': clusterId
    }
  });
}
