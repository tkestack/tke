
/**
 * 从 source 抓取 target 中定义的值
 *
 * @example
 *
   ```js
   let foo = { a: 1, b: 2, c: null };
   let bar = { a: 3, b: 4, c: 5, d: 6 };

   fetch(foo, bar); // { a: 3, b: 4, c: 5 }
   fetch(foo, bar, true); // { a: 1, b: 2, c: 5 }
   ```
 */
export function fetch<T extends Object, S extends Object>(target: T, source: S, fetchNullOnly: boolean = false): T {
    for (let prop in target) {
        if (target.hasOwnProperty(prop)) {
            if (fetchNullOnly && target[prop] !== null) break;
            if (source.hasOwnProperty(prop)) {
                target[prop] = source[prop as any];
            }
        }
    }
    return target;
}
