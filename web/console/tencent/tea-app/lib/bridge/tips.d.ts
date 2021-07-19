export interface AppTipsBridge {
    /**
     * 向用户提示成功信息
     * @param message 提示内容
     * @param duration 持续时间
     */
    success: typeof success;
    /**
     * 向用户提示错误信息
     * @param message 提示内容
     * @param duration 持续时间
     */
    error: typeof error;
    /**
     * 想用户提示加载信息
     * @param message 提示内容
     * @param duration 持续时间
     */
    loading: typeof loading;
}
/**
 * 提供全局用户提示
 */
export declare const tips: AppTipsBridge;
declare function success(options: {
    message: string;
    duration?: number;
}): void;
declare function success(message: string, duration?: number): void;
declare function error(options: {
    message: string;
    duration?: number;
}): void;
declare function error(message: string, duration?: number): void;
declare function loading(options: {
    message?: string;
    duration?: number;
}): {
    stop: () => void;
};
declare function loading(message?: string, duration?: number): {
    stop: () => void;
};
/**
 * @deprecated
 * 使用 `app.tips` 代替
 */
export declare const tip: {
    /**
     * @deprecated
     * 向用户提示成功信息
     * @param message 提示内容
     * @param duration 持续时间
     */
    success(message: string, duration?: number): void;
    /**
     * @deprecated
     * 向用户提示错误信息
     * @param message 提示内容
     * @param duration 持续时间
     */
    error(message: string, duration?: number): void;
    /**
     * @deprecated
     * 想用户提示加载信息
     * @param message 提示内容
     * @param duration 持续时间
     */
    loading(message: string, duration?: number): void;
};
export {};
