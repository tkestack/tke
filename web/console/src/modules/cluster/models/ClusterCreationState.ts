import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface ClusterCreationState extends Identifiable {
  /**链接集群名字 */
  name?: string;
  v_name?: Validation;

  /**apiServer地址 */
  apiServer?: string;
  v_apiServer?: Validation;

  /**证书 */
  certFile?: string;
  v_certFile?: Validation;

  token?: string;
  v_token?: Validation;

  jsonData?: any;

  currentStep?: number;

  clientCertificate?: string;
  clientKey?: string;
}
