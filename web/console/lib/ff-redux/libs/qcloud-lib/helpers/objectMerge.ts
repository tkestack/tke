const deepExtend = require('deep-extend');

export function merge<T, S>(target: T, source: S): T & S;
export function merge<T, S1, S2>(target: T, source1: S1, source2: S2): T & S1 & S2;
export function merge<T, S1, S2>(target: T, source1: S1, source2: S2): T & S1 & S2;
export function merge<T, S1, S2, S3>(target: T, source1: S1, source2: S2, source3: S3): T & S1 & S2 & S3;
export function merge<T, S1, S2, S3, S4>(
  target: T,
  source1: S1,
  source2: S2,
  source3: S3,
  source4,
  S4
): T & S1 & S2 & S3 & S4;

export function merge(target: any, ...sources: any[]): any {
  return deepExtend(target, ...sources);
}
