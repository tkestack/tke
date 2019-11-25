import { isEmpty } from './isEmpty';

/** 判断数组是否存在某个值 */
export const includes = (arr: Array<any> | string, value: any): boolean => {
  if (isEmpty(arr)) {
    /**如果数组为空 */
    return false;
  } else if (arr instanceof Array) {
    //判断数组是否包含元素
    return arr.findIndex(item => item === value) > -1;
  } else {
    //判断字符串是否包含子串
    return arr.indexOf(value) > -1;
  }
};
