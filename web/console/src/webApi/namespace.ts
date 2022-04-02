import { Request } from './request';

export const fetchNamespaceList = (clusterId: string) => {
  return Request.get<any, { items: any }>('/api/v1/namespaces', {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
};
