import { Request, generateQueryString } from './request';

export function fetchAlarmList(
  { clusterId },
  { limit = null, continueToken = null, query = null, alertStatus = null }
) {
  return Request.get<any, any>(
    `/apis/notify.tkestack.io/v1/messages?${generateQueryString({
      limit,
      continue: continueToken,
      fieldSelector: generateQueryString(
        {
          'spec.clusterID': clusterId,
          'spec.alarmPolicyName': query,
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
