import { Validation } from '../../common/models';
import { Identifiable } from '@tencent/qcloud-lib';

export interface RuleMap extends Identifiable {
  /** host 域名 */
  host?: string;
  v_host?: Validation;

  /** url路径 */
  path?: string;
  v_path?: Validation;

  /** 后端服务名称 */
  serviceName?: string;
  v_serviceName?: Validation;

  /** 后端服务端口 */
  servicePort?: string | number;
  v_servicePort?: Validation;

  /** protocol */
  protocol?: string;
  v_protocol?: Validation;
}
