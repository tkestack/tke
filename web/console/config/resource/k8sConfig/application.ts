import { generateResourceInfo } from '../common';

/**
 * app
 * @param k8sVersion
 */
export const app = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'app',
    requestType: {
      list: 'apps'
    },
    isRelevantToNamespace: true
  });
};
