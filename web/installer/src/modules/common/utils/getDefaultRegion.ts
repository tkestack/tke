import { Region } from '../models';
import { getRegionId, setRegionId } from '../../../../helpers';

interface RegionModel extends Region {
  regionId?: number | string;
}

export const getDefaultRegion = (urlRegion?: number | string, regionList?: Array<RegionModel>) => {
  let rid = urlRegion || getRegionId();
  if (regionList.length && regionList.some(r => r.regionId === rid)) {
    return rid;
  } else {
    return regionList[0].regionId;
  }
};
