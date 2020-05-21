import { generateResourceInfo } from '../common';

/** lbcf的相关配置 */
export const alarmRecord = (k8sVersion: string) => {
    return generateResourceInfo({
        k8sVersion,
        resourceName: 'alarmRecord',
        requestType: {
            list: 'messages'
        }
    });
};
