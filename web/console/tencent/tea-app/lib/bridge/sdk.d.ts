/**
 * 提供 SDK 注册和加载接口
 */
export declare const sdk: {
    /**
     * 注册 SDK
     *
     * @param sdkName SDK 名称
     * @param sdkFactory SDK 工厂方法，应该返回 SDK 提供的 API
     */
    register: (sdkName: string, sdkFactory: () => any) => void;
    /**
     * 加载并使用指定的 SDK
     *
     * @param sdkName SDK 名称
     */
    use: <T = any>(sdkName: string) => Promise<T>;
};
