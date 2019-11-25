
export function extend<T, S>(target: T, source: S): T & S;
export function extend<T, S1, S2>(target: T, source1: S1, source2: S2): T & S1 & S2;
export function extend<T, S1, S2>(target: T, source1: S1, source2: S2): T & S1 & S2;
export function extend<T, S1, S2, S3>(target: T, source1: S1, source2: S2, source3: S3): T & S1 & S2 & S3;
export function extend<T, S1, S2, S3, S4>(target: T, source1: S1, source2: S2, source3: S3, source4: S4): T & S1 & S2 & S3 & S4;

export function extend(target: any, ...sources: any[]): any {
    if (!target) return null;
    let hasOwnProperty = Object.prototype.hasOwnProperty;

    sources.forEach((source) => {
        if (!source) return;
        for (let prop in source) {
            if (hasOwnProperty.call(source, prop)) {
                target[prop] = source[prop];
            }
        }
    });

    return target;

}