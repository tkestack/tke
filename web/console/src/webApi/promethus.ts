import Request from './request';

export interface EnablePromethusParams {
  clusterName: string;
  resources: {
    limits: {
      cpu: number;
      memory: number;
    };
    requests: {
      cpu: number;
      memory: number;
    };
  };
  runOnMaster: boolean;
  alertRepeatInterval: number;
  notifyWebhook: string;
}

export interface EnablePromethusResponse {
  metadata: {
    name: string;
  };
}

export const enablePromethus = (params: EnablePromethusParams): Promise<EnablePromethusResponse> => {
  return Request.post('/apis/monitor.tkestack.io/v1/prometheuses', {
    apiVersion: 'monitor.tkestack.io/v1',
    kind: 'Prometheus',
    metadata: {
      generateName: 'prometheus'
    },
    spec: {
      version: 'v1.0.0',
      withNPD: false,
      ...{
        ...params,
        alertRepeatInterval: params.alertRepeatInterval + 'm'
      }
    }
  });
};

export const closePromethus = (promethusId: string) =>
  Request.delete(`/apis/monitor.tkestack.io/v1/prometheuses/${promethusId}`);
