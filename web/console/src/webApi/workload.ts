import { Request } from './request';

interface IProps {
  resourceId: string;
  namespace: string;
  clusterId: string;
  kind: string;
  data?: any;
}

export function fetchWorkloadResource({ resourceId, namespace, clusterId, kind }: IProps) {
  return Request.get<any, any>(`/apis/apps/v1/namespaces/${namespace}/${kind}s/${resourceId}`, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
}

export function updateWorkloadResource({ resourceId, namespace, clusterId, kind, data }: IProps) {
  return Request.patch(`/apis/apps/v1/namespaces/${namespace}/${kind}s/${resourceId}`, data, {
    headers: {
      'Content-Type': 'application/strategic-merge-patch+json',
      'X-TKE-ClusterName': clusterId
    }
  });
}
