import { Request } from './request';

export function fetchStorageClassList(clusterId) {
  return Request.get<any, { items: any[] }>('/apis/storage.k8s.io/v1/storageclasses', {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
}
