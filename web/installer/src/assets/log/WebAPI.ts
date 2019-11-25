import { collectionPaging } from '@tencent/qcloud-lib';
import { RecordSet, LogFilter, Log } from '../common/models';
import { includes, orderBy, operationResult, formatRequestRequest } from '../common/utils';
import { QueryState } from '@tencent/qcloud-redux-query';
import axios from 'axios';

export async function fetchLogList(query?: QueryState<LogFilter>) {
  let { search, filter, paging } = query,
    { appId, ownerUin, uin, action, status, returnCode, startTime, endTime, body, result } = filter;

  let queryStr = 'logType: * AND referer:*tke*',
    gte = new Date(startTime).getTime(),
    lte = new Date(endTime).getTime();

  for (let key in { appId, ownerUin, uin, action, status }) {
    if (filter[key]) {
      queryStr += ' AND ' + key + ':' + filter[key];
    }
  }

  for (let key in { body, result }) {
    if (filter[key]) {
      queryStr += ' AND ' + key + ':*' + filter[key] + '*';
    }
  }

  if (returnCode) {
    queryStr += ' AND returnCode:(' + returnCode + ')';
  }

  let size = 500,
    from = 0;
  if (paging) {
    size = paging.pageSize;
    from = (paging.pageIndex - 1) * paging.pageSize;
  }

  let params = {
    size: size,
    from: from,
    sort: [{ '@timestamp': { order: 'desc', unmapped_type: 'boolean' } }],
    query: {
      bool: {
        must: [
          {
            query_string: {
              analyze_wildcard: true,
              query: queryStr
            }
          },
          {
            range: {
              '@timestamp': {
                gte,
                lte,
                format: 'epoch_millis'
              }
            }
          }
        ]
      }
    }
  };

  let rsp = await axios.post('/oss/oss_api?i=log/getConsoleLog', {
    region: 1,
    serviceType: 'log',
    action: 'getConsoleLog',
    data: {
      search: JSON.stringify(params)
    }
  });

  let rspData = formatRequestRequest(rsp);
  let { dataList, length, isAuthorized, isLoginedSec, message, redirect } = rspData;

  const re: RecordSet<Log> = {
    recordCount: length,
    records: dataList,
    auth: {
      isAuthorized,
      isLoginedSec,
      message,
      redirect
    }
  };

  return re;
}
