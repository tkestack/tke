import { RecordSet } from '@tencent/qcloud-lib';
import { QueryState } from '@tencent/qcloud-redux-query';
import { Config, ConfigFilter } from '../models';
import { sendCapiRequest } from '../../../../helpers';

/**
 * 拉取所有的 Config 数据
 * @param {QueryState<ConfigFilter>} query 抓取使用的查询
 * @returns
 */
export async function fetchConfig(query: QueryState<ConfigFilter>) {
  const { date, search, filter, paging, sort } = query;
  let options = {};
  if (sort && sort.by) {
    options['orderField'] = sort.by;
    options['orderType'] = sort.desc ? 'desc' : 'asc';
  }

  if (paging) {
    const { pageIndex, pageSize } = paging;
    options['offset'] = (pageIndex - 1) * pageSize;
    options['limit'] = pageSize;
  }

  if (search) {
    options['searchValue'] = search;
  }

  const list = await sendCapiRequest('ccs', 'DescribeConfig', options);

  let configs: Config[] = list.data
    ? list.data.configInfos.map(item => Object.assign({}, item, { id: item.configId }))
    : [];

  const total = list.data.totalCount;

  const result: RecordSet<Config> = {
    recordCount: total,
    records: configs
  };

  return result;
}
