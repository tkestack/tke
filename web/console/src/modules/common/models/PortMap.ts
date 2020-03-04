import { Identifiable } from '@tencent/ff-redux';

import { BaseType, Validation } from './';

export interface PortMapItem extends Identifiable {
  /**协议（TCP或UDP）*/
  protocol?: string;
  v_protocol?: Validation;

  /**容器监听的端口 */
  containerPort?: string;
  v_containerPort?: Validation;

  /**外网lb的端口 */
  lbPort?: string;
  v_lbPort?: Validation;

  /**主机端口生成方式 */
  generateType?: string;

  /**主机的端口 */
  nodePort?: string;
  v_nodePort?: Validation;
}

/**协议类型 */
export interface Protocol extends BaseType {}
