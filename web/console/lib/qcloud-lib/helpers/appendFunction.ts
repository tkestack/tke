/**
 * 在指定函数后面追加一个函数
 */
export function appendFunction<T extends Function, R>(origin: (...args) => R, append: (r: R, ...args) => R) {
  return function appended<T>(...args) {
    let result = origin.apply(this, args);
    return append.apply(this, [result].concat(args));
  };
}
