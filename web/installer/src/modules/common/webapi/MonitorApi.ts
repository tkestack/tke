import axios from 'axios';
export async function requestMonitorData(params) {
  let rsp = await axios.post('/oss/oss_api?i=monitor/getMonitorData', {
    region: +params.region,
    serviceType: 'monitor',
    action: 'getMonitorData',
    data: {
      Appid: params.appId,
      NamespaceName: params.table,
      StartTime: params.startTime,
      EndTime: params.endTime,
      Fields: params.fields,
      Conditions: params.conditions.map(c => JSON.stringify(c)),
      OrderBy: params.orderBy,
      GroupBys: params.groupBy,
      Order: params.order,
      Limit: params.limit
    }
  });
  return {
    columns: rsp.data.data.Response.Columns,
    data: JSON.parse(rsp.data.data.Response.Data)
  };
}
