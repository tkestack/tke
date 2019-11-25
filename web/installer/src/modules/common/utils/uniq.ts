/**
 * 数组去重函数
 * @param Array<any> arr
 * @param string by 数组去重字段，非必需
 * @returns 返回去重后的数组，原数组保持不变
 */
export const uniq = (arr: Array<any>, by?: string) => {
  return arr.filter((elem, pos, arr) => {
    if (by) {
      return pos === arr.findIndex(item => JSON.stringify(item[by]) === JSON.stringify(elem[by]));
    } else {
      return pos === arr.findIndex(item => JSON.stringify(item) === JSON.stringify(elem));
    }
  });
};
