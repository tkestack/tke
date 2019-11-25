import { RecordSet } from '@tencent/qcloud-lib';
import { QueryState } from '@tencent/qcloud-redux-query';
import { Cluster, ClusterFilter, RequestParams } from '../models';
import { reduceNetworkRequest } from '../../../../helpers';

/**
 * 拉取集群列表
 * @param {QueryState<ClusterFilter>} query 抓取使用的查询
 * @returns
 */
export async function fetchClusterList(query: QueryState<ClusterFilter>) {
  let { search, filter, paging, sort } = query;

  let options = {};
  if (sort && sort.by) {
    options['orderField'] = sort.by;
    options['orderType'] = sort.desc ? 'desc' : 'asc';
  }

  if (paging) {
    let { pageIndex, pageSize } = paging;
    options['offset'] = (pageIndex - 1) * pageSize;
    options['limit'] = pageSize;
  }

  if (search) {
    options['clusterName'] = search;
  }

  if (filter.uniqVpcId) {
    options['uniqVpcId'] = filter.uniqVpcId;
  }

  if (filter.uniqSubnetId) {
    options['uniqSubnetId'] = filter.uniqSubnetId;
  }

  /** 构建参数 */
  let params: RequestParams = {
    apiParams: {
      module: 'ccs',
      interfaceName: 'DescribeCluster',
      regionId: +filter.regionId,
      restParams: options
    }
  };

  let response = await reduceNetworkRequest(params);

  let clusters: Cluster[] = [];
  let total = 0;
  if (response.code === 0) {
    clusters = response.data.clusters.map(item => Object.assign({}, item, { id: item.clusterId }));
    total = response.data.totalCount;
    if (filter && filter.status) {
      clusters = clusters.filter(x => x.status === filter.status);
    }
  }

  const result: RecordSet<Cluster> = {
    recordCount: total,
    records: clusters
  };

  return result;
}
