/**
 * 查找集合中满足给定条件的记录
 */
export function findByCondition<T>(collection: T[], condition: (item: T) => boolean): T {
  for (let i = 0; i < collection.length; i++) {
    if (condition(collection[i])) {
      return collection[i];
    }
  }
  return null;
}
