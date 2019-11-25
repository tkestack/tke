import { Identifiable } from '@tencent/qcloud-lib';
export interface User extends Identifiable {
  /** 用户名（唯一） */
  name: string;
  Spec: {
    /** 密码 */
    hashedPassword: string;

    /** 额外属性 */
    extra: {
      /** 展示名 */
      displayName: string;

      /** 邮箱 */
      email?: string;

      /** 手机号 */
      phoneNumber?: string;

      /** 是否是管理员 */
      platformadmin?: boolean;
      [props: string]: any;
    };
  };
  [props: string]: any;
}
