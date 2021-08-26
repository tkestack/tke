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

export const NestedMetricsResourceMap = {
  cpu: {
    targetAverageUtilization: { key: 'cpuUtilization', meaning: 'CPU利用率', unit: '%' },
    targetAverageValue: { key: 'cpuAverage', meaning: 'CPU使用量', unit: '核' }
  },
  cpuAverage: {
    targetAverageUtilization: { key: 'cpuUtilization', meaning: 'CPU利用率', unit: '%' },
    targetAverageValue: { key: 'cpuAverage', meaning: 'CPU使用量', unit: '核' }
  },
  cpuUtilization: {
    targetAverageUtilization: { key: 'cpuUtilization', meaning: 'CPU利用率', unit: '%' },
    targetAverageValue: { key: 'cpuAverage', meaning: 'CPU使用量', unit: '核' }
  },
  memory: {
    targetAverageUtilization: { key: 'memoryUtilization', meaning: '内存利用率', unit: '%' },
    targetAverageValue: { key: 'memoryAverage', meaning: '内存使用量', unit: 'Mib' }
  },
  memoryAverage: {
    targetAverageUtilization: { key: 'memoryUtilization', meaning: '内存利用率', unit: '%' },
    targetAverageValue: { key: 'memoryAverage', meaning: '内存使用量', unit: 'Mib' }
  },
  memoryUtilization: {
    targetAverageUtilization: { key: 'memoryUtilization', meaning: '内存利用率', unit: '%' },
    targetAverageValue: { key: 'memoryAverage', meaning: '内存使用量', unit: 'Mib' }
  }
};

export const MetricsResourceMap = {
  cpuUtilization: { key: 'cpuUtilization', meaning: 'CPU利用率', unit: '%' },
  cpuAverage: { key: 'cpuAverage', meaning: 'CPU使用量', unit: '核' },
  memoryUtilization: { key: 'memoryUtilization', meaning: '内存利用率', unit: '%' },
  memoryAverage: { key: 'memoryAverage', meaning: '内存使用量', unit: 'Mib' }
};
