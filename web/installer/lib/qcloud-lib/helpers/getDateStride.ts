import { getDateDelta } from './getDateDelta';

function rangeLength(from: string, to: string) {
  return Math.abs(getDateDelta(from, to)) + 1;
}

export function getDateStride(dateFrom: string, dateTo: string) {
  let length = rangeLength(dateFrom, dateTo);

  /*
        小于等于1天，5分钟粒度数据。
        大于1天小于等于7天，1小时粒度数据。
        大于7天，1天粒度数据。
    */
  if (length <= 1) {
    return 5;
  } else if (length <= 7) {
    return 60;
  } else {
    return 24 * 60;
  }
}
