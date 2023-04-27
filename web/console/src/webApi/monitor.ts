import { Request } from './request';

export const fetchMetrics = async ({ startTime, endTime, table, groupBy, fields, conditions }) => {
  const { jsonResult } = await Request.post<any, any>('/apis/monitor.tkestack.io/v1/metrics', {
    apiVersion: 'monitor.tkestack.io/v1',
    kind: 'Metric',
    query: {
      startTime,
      endTime,
      table,
      orderBy: 'timestamp',
      order: 'asc',
      offset: 0,
      limit: 65535,
      groupBy,
      fields,
      conditions
    }
  });

  return JSON.parse(jsonResult);
};
