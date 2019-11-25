import { Validation } from '../../common/models';
import { Identifiable } from '@tencent/qcloud-lib';

export interface PortMap extends Identifiable {
  /**
   * 协议（TCP或UDP）
   */
  protocol?: string;
  v_protocol?: Validation;

  /**
   * 容器监听的端口
   */
  targetPort?: string;
  v_targetPort?: Validation;

  /**
   * 服务端口的端口
   */
  port?: string;
  v_port?: Validation;

  /**
   * 主机的端口
   */
  nodePort?: string;
  v_nodePort?: Validation;
}
