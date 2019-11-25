export const isEmpty = (value: any): boolean => {
  if (Array.isArray(value)) {
    //value为数组
    return !value.length;
  } else if (typeof value === 'object') {
    //value为对象
    if (value === null) {
      //value为null
      return true;
    } else {
      //value是否没有key
      return !Object.keys(value).length;
    }
  } else if (typeof value === 'undefined') {
    //value为undefinded
    return true;
  } else if (Number.isFinite(value)) {
    //value为数值
    return false;
  } else {
    //value为默认值
    return !value;
  }
};
