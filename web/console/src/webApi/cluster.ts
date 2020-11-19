import Request from './request';

export const getK8sValidVersions = () => {
  return Request.get('/v1/namespaces/kube-public/configmaps/cluster-info');
};
