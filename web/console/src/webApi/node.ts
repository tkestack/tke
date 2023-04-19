import { Request } from './request';

export function fetchNodeList({ clusterId }) {
  return Request.get<any, any>(`/api/v1/nodes`, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
}
