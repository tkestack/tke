import * as React from "react";
/**
 * 为模块生成渲染和销毁方法
 */
export declare const factory: EntryModuleFactory;
export interface EntryModuleFactory {
    (controller: string, action: string, entry: AppEntry): () => EntryModule;
}
export interface EntryModule {
    render(): void;
    destroy(): void;
    title?: string;
}
/**
 * 单个路由入口定义，包含渲染方法、菜单和 CSS 定义
 */
export declare type AppEntry = React.ComponentType<any> | AppEntryDetail | 404;
export declare type AppEntryDetail = {
    /**
     * 路由渲染组件
     */
    component?: React.ComponentType<any>;
    /**
     * 路由渲染方法，需返回一个 React.ReactNode
     */
    render?: () => JSX.Element;
    /**
     * 该路由下所有的文档标题
     *
     * 如果希望动态设置标题，可以使用 useDocumentTitle() / setDocumentTitle()
     */
    title?: string;
};
