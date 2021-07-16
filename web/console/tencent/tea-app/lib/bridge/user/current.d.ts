export interface AppUserData {
    /** 是否为主账号 */
    isOwner: boolean;
    /** 当前用户登录的 UIN */
    loginUin: number;
    /** 当前用户登录的主账号 UIN */
    ownerUin: number;
    /** 当前用户登录的主账号 APPID */
    appId: number;
    /** 当前用户的实名认证信息，如果未实名认证，此字段为 null */
    identity: AppUserIdentityInfo | null;
    /** 用户昵称 */
    nickName?: string;
    /** 用户标识名称，含开发商信息 */
    displayName?: string;
}
export interface AppUserIdentityInfo {
    /**
     * 认证主体类型
     *   - `0`: 个人
     *   - `1`: 企业
     */
    subjectType: number;
    /**
     * 认证渠道
     *   - `0`: 未知
     *   - `1`: 有效证件（个人：身份证/护照，企业：营业执照）
     *   - `2`: 财付通（个人）
     *   - `3`: 银行卡（企业）
     *   - `4`: 微信（个人）
     *   - `5`: 手Q（个人）
     *   - `6`: 公众平台（企业）
     *   - `7`: 线下认证（很少）
     *   - `8`: 国际信用卡（个人/企业）
     *   - `9`: 企业线下打款
     *   - `10`: 线上申请，线下审核（企业修改实名认证流程）
     *   - `11`: 米大师认证
     *   - `12`: 个人人脸核身认证
     *   - `20`: 代理商
     */
    authType: number;
    /**
     * 认证地区
     *   - `-1`: 未知
     *   - `0`: 大陆
     *   - `1`: 港澳
     *   - `2`: 台湾
     *   - `3`: 外籍
     */
    authArea: number;
}
export declare const current: () => Promise<AppUserData>;
