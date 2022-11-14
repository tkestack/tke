import { Request, generateQueryString } from './request';

export function fetchAlarmList({ clusterId }, { limit = null, continueToken = null, query = {}, alertStatus = null }) {
  return Request.get<any, any>(
    `/apis/notify.tkestack.io/v1/messages?${generateQueryString({
      limit,
      continue: continueToken,
      fieldSelector: generateQueryString(
        {
          ...query,
          'spec.clusterID': clusterId,
          'status.alertStatus': alertStatus
        },
        ','
      )
    })}`,
    {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    }
  );
}
