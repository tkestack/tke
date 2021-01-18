import Request from './request';

export const getTkeStackVersion = async () => {
  const rsp = await Request.get<any, { items: Array<{ data?: { tkeVersion?: string } }> }>(
    '/api/v1/namespaces/kube-public/configmaps',
    {
      headers: {
        'X-TKE-ClusterName': 'global'
      }
    }
  );
  return rsp?.items?.[0]?.data?.tkeVersion ?? '';
};
