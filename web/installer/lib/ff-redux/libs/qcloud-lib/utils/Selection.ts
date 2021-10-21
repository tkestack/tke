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
 * @fileOverview
 *
 * 有序维护选区
 */

import { Identifiable } from '../types/Identifiable';

/**
 * 保证 selection 保持与 records 有相同的偏序关系，查找插入位置
 */
function findInsertPosition<T extends Identifiable>(records: T[], selection: T[], item: T): number {
  let recordIndex = records.indexOf(item);
  let ri: number = 0;
  let si: number = 0;
  while (records[ri] && selection[si]) {
    if (ri === recordIndex) break;
    if (records[ri] === selection[si]) {
      si++;
    }
    ri++;
  }
  return si;
}

/**
 * 插入到选区中，保持与 records 相同的偏序顺序
 */
export function selectionInsert<T extends Identifiable>(records: T[], selection: T[], item: T): T[] {
  let newSelection = selection ? selection.slice() : [];
  const index = newSelection.indexOf(item);

  if (index === -1) {
    const insertIndex = findInsertPosition(records, newSelection, item);
    newSelection.splice(insertIndex, 0, item);
  }

  return newSelection;
}

/**
 * 从选取中删除
 */
export function selectionRemove<T extends Identifiable>(records: T[], selection: T[], item: T): T[] {
  let newSelection = selection.slice();
  const index = newSelection.indexOf(item);

  if (index > -1) {
    newSelection.splice(index, 1);
  }

  return newSelection;
}
