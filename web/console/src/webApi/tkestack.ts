import Request from './request';

export const getTkeStackVersion = async () => {
  const rsp = await Request.get<any, { data?: { tkeVersion: string } }>(
    '/api/v1/namespaces/kube-public/configmaps/cluster-info',
    {
      headers: {
        'X-TKE-ClusterName': 'global'
      }
    }
  );
  return rsp?.data?.tkeVersion ?? '';
};
