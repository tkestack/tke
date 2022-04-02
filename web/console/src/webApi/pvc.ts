import { Request } from './request';

export function fetchPVCInfo({ namespace, name, clusterId }) {
  return Request.get<any, any>(`/api/v1/namespaces/${namespace}/persistentvolumeclaims/${name}`, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
}
