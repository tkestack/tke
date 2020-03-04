import { hashString } from './hashString';

export function hashObject(obj: any) {
  if (obj.prototype) {
    throw new RangeError('can only hash a plain object');
  }
  return hashString(JSON.stringify(obj));
}
