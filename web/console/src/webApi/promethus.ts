/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

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
  Request.delete(`/apis/monitor.tkestack.io/v1/prometheuses/${promethusId}`);
