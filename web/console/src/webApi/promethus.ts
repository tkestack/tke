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
  return Request.post('monitor.tkestack.io/v1/prometheuses', {
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
        alertRepeatInterval: params.alertRepeatInterval + 'm',
        resources: {
          limits: {
            cpu: params.resources.limits.cpu,
            memory: params.resources.limits.memory + 'Mi'
          },
          requests: {
            cpu: params.resources.requests.cpu,
            memory: params.resources.requests.memory + 'Mi'
          }
        }
      }
    }
  });
};

export const closePromethus = (promethusId: string) =>
  Request.delete(`monitor.tkestack.io/v1/prometheuses/${promethusId}`);
