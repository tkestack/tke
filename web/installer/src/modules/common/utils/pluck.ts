export const pluck = (arr: Array<any>, ckey: string) => {
  let res = [];
  if (Array.isArray(arr) && !!ckey) {
    arr.forEach(item => {
      if (item[ckey] !== undefined) {
        res.push(item[ckey]);
      }
    });
  }

  return res;
};
