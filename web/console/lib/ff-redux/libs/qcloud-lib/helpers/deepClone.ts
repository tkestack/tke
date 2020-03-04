/**
 * 深度拷贝值
 */
export function deepClone<T>(source: T): T {
  if (source instanceof Array) {
    return deepCloneArray(source as any) as any;
  }

  if (source && typeof source === 'object') {
    return deepCloneObject(source);
  }

  return source;
}

function deepCloneObject(source: any): any {
  const target: any = {};

  Object.keys(source).forEach(key => {
    target[key] = deepClone(source[key]);
  });

  return target;
}

function deepCloneArray(source: any[]) {
  return source.map(x => deepClone(x));
}
