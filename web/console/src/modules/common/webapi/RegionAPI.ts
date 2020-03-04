import { QueryState, RecordSet, uuid } from '@tencent/ff-redux';

import { Region, RegionFilter } from '../models';

/**获取地域列表 */
export async function fetchRegionList(query?: QueryState<RegionFilter>) {
  // 目前是hardcode，后面换成接口获取
  let regionList = [
    {
      id: uuid(),
      Remark: 'SUITABLE_TKE',
      area: '华南地区',
      name: '广州',
      value: 1
    },
    {
      id: uuid(),
      Remark: 'SUITABLE_TKE',
      area: '华东地区',
      name: '上海',
      value: 4
    }
  ];

  const result: RecordSet<Region> = {
    recordCount: regionList.length,
    records: regionList
  };

  return result;
}
