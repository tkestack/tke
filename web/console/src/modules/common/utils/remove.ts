/**
 * @returns 删除数据后的新数组，改变原数组
 */
export const remove = (arr: Array<any>, func: (value: any, index?: number, array?: any[]) => any) => {
  let targets = arr.filter(func);

  targets.forEach(elem => {
    arr.forEach((value, index) => {
      if (JSON.stringify(value) === JSON.stringify(elem)) {
        arr.splice(index, 1);
      }
    });
  });

  return arr;
};
