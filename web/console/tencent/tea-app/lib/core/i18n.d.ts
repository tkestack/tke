/// <reference types="react" />
import i18next from "i18next";
/**
 * 国际化包括的信息以及所需的工具
 */
export declare const i18n: {
    /**
     * 当前用户的国际化语言，已知语言：
     *  - `zh` 中文
     *  - `en` 英文
     *  - `jp` 日语
     *  - `ko` 韩语
     */
    lng: string;
    /**
     * 当前用户所在站点
     *  - `1` 表示国内站；
     *  - `2` 表示国际站；
     *
     * @type {1 | 2}
     */
    site: number;
    /**
     * 注册国家
     */
    country: {
        name: string;
        code: string;
    };
    /**
     * 初始化当前语言的国际化配置
     */
    init: ({ translation }: I18NInitOptions) => void;
    /**
     * 标记翻译句子
     * 详细的标记说明，请参考 http://tapd.oa.com/QCloud_2015/markdown_wikis/view/#1010103951008390841
     */
    t: (sentence: string, options?: I18NTranslationOptions) => string;
    /**
     * 标记翻译组件
     * 详细的标记说明，请参考 http://tapd.oa.com/QCloud_2015/markdown_wikis/view/#1010103951008390841
     */
    Trans: import("react").ComponentClass<import("react-i18next").TransProps, any>;
};
export declare const getI18NInstance: () => i18next.i18n;
/**
 * @internal 国际化容器，内部使用
 */
export declare const I18NProvider: import("react").ComponentClass<import("react-i18next").I18nextProviderProps, any>;
export interface I18NInitOptions {
    translation: I18NTranslation;
}
export interface I18NTranslation {
    [key: string]: string;
}
export interface I18NTranslationOptions {
    /** 用于确定单复数的数量值 */
    count?: number;
    /** 用于确定上下文的说明文本，只能使用字符串常量，否则无法扫描 */
    context?: string;
    [key: string]: any;
}
