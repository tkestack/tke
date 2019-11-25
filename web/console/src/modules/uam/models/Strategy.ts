import { Identifiable } from '@tencent/qcloud-lib';
export interface Strategy extends Identifiable {
  /** 策略名字 */
  name: string;

  /** 描述 */
  description?: string;

  statement: {
    /** 操作 */
    action: Array<string>;

    /** 资源 */
    resource: string;

    /** 效果 */
    effect: string;
  };
  [props: string]: any;
}
