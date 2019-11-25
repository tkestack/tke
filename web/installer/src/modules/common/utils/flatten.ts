const flat = (arr: Array<any>) => {
  return arr.reduce((prev, cur) => prev.concat(cur), []);
};

export const flattenDeep = (arr: Array<any>) => {
  let flatten = flat(arr);

  if (!flatten.filter(x => Array.isArray(x)).length) {
    return flatten;
  }
  return flattenDeep(flatten);
};

export const flatten = (arr: Array<any>, isDeep?: boolean) => {
  if (isDeep) {
    return flattenDeep(arr);
  }
  return flat(arr);
};
