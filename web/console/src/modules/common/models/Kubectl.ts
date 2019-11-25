import { Identifiable } from '@tencent/qcloud-lib';

export interface Kubectl extends Identifiable {
  /**用户名 */
  userName?: string;

  /**密码 */
  password?: string;

  /**凭证 */
  certificationAuthority?: string;

  /**外网访问地址 */
  clusterExternalEndpoint?: string;
}

export interface KubectlFilter {
  /**集群Id */
  clusterId?: string;

  /**地域Id */
  regionId?: number;
}
