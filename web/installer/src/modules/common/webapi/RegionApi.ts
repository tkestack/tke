import { RecordSet } from '@tencent/qcloud-lib';
import { Region, RegionFilter, RequestParams } from '../models';
import { regionMapList } from '../../../../config/region';
import { reduceNetworkRequest } from '../../../../helpers';
import { QueryState } from '@tencent/qcloud-redux-query';

/**
 * 地域的查询
 * @param query 地域查询的一些过滤条件
 * @param isNeedFilter: boolean 判断是否需要过滤 SUITABLE_TKE 的地域出来
 */
export async function fetchRegionList(query: QueryState<RegionFilter>, isNeedFilter: boolean = false) {
  let regionList = [],
    afterRegionList = [];

  /** 构建参数 */
  let params: RequestParams = {
    apiParams: {
      interfaceName: 'DescribeCCSRegion',
      regionId: 1,
      module: 'ccs'
    }
  };

  let response = await reduceNetworkRequest(params);

  if (response.code === 0) {
    regionList = response.data.regions
      .sort((prev, next) => prev.Id - next.Id)
      .map(item => ({ Remark: item.Remark, regionId: item.RegionId }));

    regionList.forEach(item => {
      if (regionMapList[item.regionId]) {
        afterRegionList.push(Object.assign({}, regionMapList[item.regionId], { Remark: item.Remark }));
      }
    });

    const result: RecordSet<Region> = {
      recordCount: afterRegionList.length,
      records: afterRegionList
    };

    return result;
  }
}
