/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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
