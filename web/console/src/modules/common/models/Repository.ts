import { Identifiable } from '@tencent/ff-redux';

export interface Repository extends Identifiable {
  /**镜像名称 */
  reponame?: string;

  /**镜像类型 */
  repotype?: string;

  /**tag数量 */
  tagCount?: number;

  /**是否公有 */
  public?: number;

  /**是否被用户收藏 */
  isUserFavor?: boolean;

  /**是否为Qcloud官方镜像 */
  isQcloudOfficial?: boolean;

  /**收藏数 */
  favorCount?: number;

  /**下载数 */
  pullCount?: number;

  /**描述*/
  description?: string;

  /**仓库地址 */
  address?: string;

  /**创建时间 */
  creationTime?: string;

  /**logo地址 */
  logo?: string;

  /**镜像描述（仅官方镜像有）*/
  simpleDesc?: string;

  /**地域 */
  regionId?: string | number;
}

export interface RepositoryFilter {
  /**镜像类别 */
  repotype?: string;

  /**镜像名称 */
  reponame?: string;

  /**是否公有 */
  public?: number;

  /**命名空间 */
  namespace?: string;

  /**地域 */
  regionId?: string | number;
}
