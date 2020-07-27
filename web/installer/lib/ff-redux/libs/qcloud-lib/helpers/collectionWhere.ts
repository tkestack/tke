/**
 * 集合筛选
 * */
export function collectionWhere<T>(collection: T[], condition: (x: T) => boolean) {
  return collection.reduce((found: T[], current: T) => (condition(current) ? found.concat(current) : found), []);
}
