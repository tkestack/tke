import { generateResourceInfo } from '../common';

/** lbcf的相关配置 */
export const alarmPolicy = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'alarmPolicy',
    requestType: {
      list: 'alarmpolicies'
    }
  });
};
