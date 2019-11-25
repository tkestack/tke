import { Price } from '../models';
import { sendCapiRequest, localCache } from '../../../../helpers';
import { isInWhitelist } from '../../../../helpers/Whitelist';

/**
 * 判断是否为带宽上移的用户
 */
export async function checkIsNetworkShift() {
  let isNetworkShiftUser = localCache.isNetworkShiftUser;

  if (!isNetworkShiftUser) {
    let response = await isInWhitelist('network_shift_up');
    localCache.isNetworkShiftUser = !!response;
    isNetworkShiftUser = !!response;
  }

  return isNetworkShiftUser;
}

/**
 * 查询带宽的费用
 * @param devPaymode:string     LB的付费模式 monthly: 包年包月，hourly：按量付费
 * @param loadBalancerType:number   2: open, 3: internal
 * @param lbChargeType:string   PREPAID: 预付费，POSTPAID_BY_HOUR： 后付费； 目前仅支持 后付费
 * @param bandwidthType:string  PayByHour: 按带宽付费；PayByTraffic：按流量付费
 * @param bandwidth:number  带宽的大小
 * @param regionId:number   地域
 */
export async function queryBandWidthPrice(
  devPaymode: string,
  bandwidthType: string,
  bandwidth: number,
  regionId: number,
  loadBalancerType: number = 2
): Promise<Price> {
  let params = {
    loadBalancerType,
    lbChargeType: devPaymode === 'hourly' ? 'POSTPAID_BY_HOUR' : 'PREPAID',
    internetAccessible: {
      internetChargeType: bandwidthType === 'PayByHour' ? 'BANDWIDTH_POSTPAID_BY_HOUR' : 'TRAFFIC_POSTPAID_BY_HOUR',
      internetMaxBandwidthOut: bandwidth
    }
  };

  let re = await sendCapiRequest('lb', 'InquiryLBPriceAll', params, regionId);

  let result = re.price;

  if (devPaymode === 'monthly') {
    let monthlyPrice = {};

    if (bandwidthType === 'PayByHour') {
      monthlyPrice = {
        price: result.networkPrice.discountPrice + result.lbIdPrice.discountPrice,
        originalPrice: result.networkPrice.originalPrice + result.lbIdPrice.originalPrice
      };
    } else {
      let lb = {
        price: result.lbIdPrice.discountPrice,
        originalPrice: result.lbIdPrice.originalPrice
      };
      let bandwidth = {
        price: result.networkPrice.unitPrice,
        price_unit: result.networkPrice.chargeUnit
      };
      monthlyPrice = {
        lb,
        bandwidth
      };
    }

    return { monthlyPrice };
  } else {
    let lb = {
        price: result.lbIdPrice.unitPrice,
        price_unit: result.lbIdPrice.chargeUnit
      },
      bandwidth = {
        price: result.networkPrice.unitPrice,
        price_unit: result.networkPrice.chargeUnit
      },
      hourlyPrice = {
        lb,
        bandwidth
      };
    return { hourlyPrice };
  }
}
