import { RecordSet, uuid } from '@tencent/qcloud-lib';
import { QueryState } from '@tencent/qcloud-redux-query';
import { Version, VersionFilter } from '../models';
import { sendCapiRequest } from '../../../../helpers';

/**
 * 拉去所有的 Version 数据
 * @param {QueryState<VersionFilter>} query 抓取使用的查询
 * @returns
 */
export async function fetchVersion(query: QueryState<VersionFilter>) {
  const { date, search, filter, paging, sort } = query;
  let rid = filter.regionId ? filter.regionId : 1;
  let options = {
    configId: filter.configId
  };

  if (paging) {
    const { pageIndex, pageSize } = paging;
    options['offset'] = (pageIndex - 1) * pageSize;
    options['limit'] = pageSize;
  }

  const list = await sendCapiRequest('ccs', 'DescribeConfigVersion', options, rid as number);

  let versions: Version[] = list.data
    ? list.data.versionInfos.map(item => {
        let id = uuid();
        return Object.assign({}, item, { id });
      })
    : [];

  const total = list.data.totalCount;

  const result: RecordSet<Version> = {
    recordCount: total,
    records: versions
  };

  return result;
}
