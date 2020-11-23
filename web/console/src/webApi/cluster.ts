import Request from './request';

export const getK8sValidVersions = () => {
  return Request.get<any, { k8sValidVersions: Array<string> }>('/v1/namespaces/kube-public/configmaps/cluster-info');
};
