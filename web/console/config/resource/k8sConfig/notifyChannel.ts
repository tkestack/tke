import { generateResourceInfo } from '../common';

/** notify的相关配置 */
export const notifyChannel = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'channel',
    requestType: {
      list: 'channels'
    }
  });
};

export const notifyTemplate = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'template',
    isRelevantToNamespace: true,
    requestType: {
      list: 'templates'
    }
  });
};

export const notifyMessage = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'message',
    requestType: {
      list: 'messages'
    }
  });
};
export const notifyReceiver = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'receiver',
    requestType: {
      list: 'receivers'
    }
  });
};
export const notifyReceiverGroup = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'receiverGroup',
    requestType: {
      list: 'receivergroups'
    }
  });
};
