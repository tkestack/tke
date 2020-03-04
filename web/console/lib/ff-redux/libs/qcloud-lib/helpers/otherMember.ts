/**
 * 查找一个集合中不再排除集合里的部分
 *
 * @example
 *
 * ```ts
 *     console.log(otherMember([1, 2, 3, 4, 5], [2, 4])); // [1, 3, 5]
 * ```
 * */
export function otherMember<T>(allMembers: T[], excludeMembers: T[]): T[] {
  return allMembers.reduce((collected, member) => {
    return excludeMembers.indexOf(member) > -1 ? collected : collected.concat(member);
  }, []);
}
