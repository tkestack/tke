import { OperationResult } from '@tencent/ff-redux';

// 返回标准操作结果
export const operationResult = <T>(target: T[] | T, error?: any): OperationResult<T>[] => {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error }));
  }
  return [{ success: !error, target: target as T, error }];
};
