import { RecordSet, uuid } from '@tencent/qcloud-lib';
import { QueryState } from '@tencent/qcloud-redux-query';
import { SubRouterFilter, SubRouter, CLB } from './models';
import { RequestParams } from '../common/models';
import { subRouterConfig } from '../../../config';
import { reduceNetworkRequest } from '../../../helpers';

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

/**
 * 获取CLB的列表
 * @param regionId: number
 */
export async function fetchCLBList(regionId: number) {
  let clbList = [];

  /** 构建参数 */
  let params: RequestParams = {
    apiParams: {
      interfaceName: 'DescribeLoadBalancers',
      regionId,
      module: 'lb',
      restParams: {
        anycast: -1, // 1 为anycast，不传则不查anycast，-1则为全部
        forward: 0, // 0 为传统型，1为应用型,
        offset: 0,
        limit: 100,
        loadBalancerTypes: ['2', '3'] // 2 为公网，3为内网
      }
    }
  };

  try {
    let response = await reduceNetworkRequest(params);

    if (response.code === 0) {
      let clbData: CLB[] = response.loadBalancerSet;
      clbList = clbData.map(clb => {
        return Object.assign({}, clb, { id: clb.loadBalancerId });
      });
    }
    return clbList;
  } catch (error) {
    return clbList;
  }
}
