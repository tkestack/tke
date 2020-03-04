import { Identifiable } from '@tencent/ff-redux';

export interface User extends Identifiable {
  metadata?: {
    /** 用户的资源id */
    name: string;
  };

  /** 用户名（唯一） */
  spec: {
    /** 用户名 */
    username: string;

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
