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
