import { Identifiable } from '@tencent/ff-redux';

export interface User extends Identifiable {
    metadata?: {
        /** 用户的资源id */
        name: string;
    };

    /** 用户名（唯一） */
    spec: {
        /** 用户名 */
        name: string;

        /** 展示名 */
        displayName: string;

        /** 邮箱 */
        email?: string;

        /** 手机号 */
        phoneNumber?: string;

        /** 密码 */
        hashedPassword: string;

        /** 额外属性 */
        extra?: {
            /** 是否是管理员 */
            platformadmin?: boolean;
            [props: string]: any;
        };
    };

    status?: any;
}

export interface UserFilter {
    /** 业务Id */
    projectId?: string;

    /** 用户名(唯一) */
    name?: string;

    /** 展示名 */
    displayName?: string;

    /** 相关参数 */
    search?: string;

    ifAll?: boolean;

    /** 是否只拉取策策略需要绑定的用户 */
    isPolicyUser?: boolean;
}

/** 原有User在localidentities和user概念之间混用了，anyway，用于关联角色等 */
export interface UserPlain extends Identifiable {
    /** 名称 */
    name?: string;
    /** 展示名 */
    displayName?: string;
}

export interface Member extends Identifiable {
    projectId: string;
    users: any;
    policies: string[];
}
