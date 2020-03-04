declare namespace nmc {
  interface Tips {
    /**
     * 显示成功提醒
     */
    success(message: string, duration?: number): void;

    /**
     * 显示错误提醒
     */
    error(message: string, duration?: number): void;

    /**
     * 显示加载指示器
     */
    showLoading(loadingText?: string): void;

    /**
     * 停止加载指示器
     */
    stopLoading(): void;
  }
}
