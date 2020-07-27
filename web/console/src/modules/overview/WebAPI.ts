import { RequestParams } from './../common/models/requestParams';
import { Method, reduceNetworkRequest } from '@helper/reduceNetwork';
import { OperationResult, RecordSet, uuid } from '@tencent/ff-redux';

export async function fetchClusteroverviews(query) {
  let url = 'apis/monitor.tkestack.io/v1/clusteroverviews';
  let params: RequestParams = {
    method: Method.post,
    url,
    data: {
      apiVersion: 'monitor.tkestack.io/v1',
      kind: 'ClusterOverview'
    }
  };

  let response = await reduceNetworkRequest(params);
  if (response.code === 0) {
    console.log(response);
    return response.data.result;
  }
}
