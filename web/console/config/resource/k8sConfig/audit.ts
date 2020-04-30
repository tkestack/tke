import { generateResourceInfo } from '../common';


export const audit = (k8sVersion: string) => {
    return generateResourceInfo({
        k8sVersion,
        resourceName: 'audit',
        requestType: {
            list: 'events'
        }
    });
};
//
// export const auditFilterValues = (k8sVersion: string) => {
//     return generateResourceInfo({
//         k8sVersion,
//         resourceName: 'event',
//         requestType: {
//             list: 'events'
//         }
//     });
// };
