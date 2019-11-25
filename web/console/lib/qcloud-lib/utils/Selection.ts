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
export function selectionInsert<T extends Identifiable>(records: T[], selectionParams: T[], item: T): T[] {
  let selection = selectionParams ? selectionParams.slice() : [];
  const index = selection.indexOf(item);

  if (index === -1) {
    const insertIndex = findInsertPosition(records, selection, item);
    selection.splice(insertIndex, 0, item);
  }

  return selection;
}

/**
 * 从选取中删除
 */
export function selectionRemove<T extends Identifiable>(records: T[], selectionParams: T[], item: T): T[] {
  let selection = selectionParams.slice();
  const index = selection.indexOf(item);

  if (index > -1) {
    selection.splice(index, 1);
  }

  return selection;
}
