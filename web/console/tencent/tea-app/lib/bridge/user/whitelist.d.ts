/**
 * 批量检查用户白名单
 */
export declare const checkWhitelistBatch: (keys: string[]) => Promise<{
    [key: string]: number;
}>;
/**
 * 检查白名单
 */
export declare const checkWhitelist: (key: string) => Promise<number>;
