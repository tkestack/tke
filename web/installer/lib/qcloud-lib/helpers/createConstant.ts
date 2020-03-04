/** 从枚举类型创建 key-value 一样的常量 */
export function createConstant<T>(enumType: T): T {
  let constant = {};
  for (let key in enumType) {
    if (enumType.hasOwnProperty(key) && isNaN(+key)) {
      constant[key as any] = key;
    }
  }
  return constant as T;
}
