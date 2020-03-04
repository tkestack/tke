import { Identifiable } from '@tencent/ff-redux';

export interface Tag extends Identifiable {
  /**Tag名称 */
  tagName?: string;

  /**Parent */
  parent?: string;

  /**durationDays */
  durationDays?: string;

  /**architecture */
  architecture?: string;

  /**dockerVersion */
  dockerVersion?: string;

  /**OS */
  OS?: string;

  /**author */
  author?: string;

  /**创建时间 */
  creationTime: string;
}
