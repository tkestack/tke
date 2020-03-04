import { RecordSet, uuid } from '@tencent/qcloud-lib';
import { SubRouterFilter, SubRouter, CLB } from './models';
import { subRouterConfig } from '../../../config';
import { QueryState } from '@tencent/ff-redux';

/** 将各种资源的接口导出 */
export * from './WebAPI/index';

/**
 * subRouter列表的拉取
 * @param query subRouter列表的查询
 */
export async function fetchSubRouterList(query: QueryState<SubRouterFilter>) {
  let response = subRouterConfig(query.filter.module);

  let subRouterList = response.map(item => {
    return Object.assign({}, item, { id: uuid() });
  });

  const result: RecordSet<SubRouter> = {
    recordCount: subRouterList.length,
    records: subRouterList
  };

  return result;
}
