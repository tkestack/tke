import { Identifiable } from '@tencent/qcloud-lib';

export interface Addon extends Identifiable {
  /** 最新的版本 */
  latestVersion?: string;

  /** 组件的级别 */
  level?: string;

  /** metadata */
  metadata: AddonMetadata;

  /** 类型 */
  type?: string;

  /** 组件的相关描述 */
  description: string;
}

interface AddonMetadata {
  /** 创建的时间 */
  creationTimestamp?: string;

  /** addon的名称 */
  name?: string;
}
