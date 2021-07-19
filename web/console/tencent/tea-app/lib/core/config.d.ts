import { EntryModuleFactory } from "./entry";
/**
 * 在原有的模块导出的基础上，附加菜单和样式的配置导出
 */
export declare const config: Config;
export interface Config {
    (factory: EntryModuleFactory): EntryModuleFactory;
}
