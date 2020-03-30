import { string } from 'prop-types';

import { extend, Identifiable, RecordSet, FFListModel } from '@tencent/ff-redux';

import { KeyValue, Validation } from '../../common/models';
import { Selector, Resource, ResourceFilter } from '../models';

export interface LbcfEdit extends Identifiable {
  name?: string;
  v_name?: Validation;

  namespace?: string;

  /** namespace */
  v_namespace?: Validation;

  config?: KeyValue[];

  v_config?: Validation;

  args?: KeyValue[];

  v_args?: Validation;

  driver?: FFListModel<Resource, ResourceFilter>;

  /** vpcSelection */
  // vpcSelection?: string;

  // clbList?: CLB[];

  // clbSelection?: string;
  // v_clbSelection?: Validation;

  // createLbWay?: string;

  /**lbcfBackGroup */

  lbcfBackGroupEditions?: GameBackgroupEdition[];

  // gameAppList?: Resource[];

  // gameAppSelection?: string;

  // isShowGameAppDialog?: boolean;
}

/** 创建负载均衡的时候，提交的jsonSchema */
export interface LbcfLBJSONYaml {
  /** 资源的类型 */
  kind?: string;

  /** api的版本 */
  apiVersion?: string;

  /** metadata */
  metadata?: LbcfLBMetadata;

  /** spec */
  spec?: LbcfLBSpec;

  /** status */
  status?: LbcfLBStatus;
}

/** metadata的配置，非全部选项 */
interface LbcfLBMetadata {
  annotations?: {
    [props: string]: string;
  };

  /** ingress的名称 */
  name?: string;

  /** ingress的命名空间 */
  namespace?: string;
}

/** spec的配置项，非全部选项 */
interface LbcfLBSpec {
  lbDriver: string;
  lbSpec?: LbSpec;
  attributes?: LbAttribute;
}

interface LbSpec {
  loadBalancerID?: string;
  vpcID?: string;
  loadBalancerType?: string;
  subnetID?: string;
  listenerPort?: string;
  listenerProtocol?: string;
  domain?: string;
  url?: string;
}

interface LbAttribute {
  listenerCertID?: string;
}

/** status的配置项，非全部选项 */
interface LbcfLBStatus {}

interface Port {
  id?: string;
  portNumber: string;
  protocol: string;
  v_portNumber?: Validation;
}

interface StringArray {
  id?: string;
  value: string;
  v_value?: Validation;
}

/** 创建backGroup的时候，提交的jsonSchema */
export interface LbcfBGJSONYaml {
  /** 资源的类型 */
  kind?: string;

  /** api的版本 */
  apiVersion?: string;

  /** metadata */
  metadata?: LbcfBGMetadata;

  /** spec */
  spec?: LbcfBGSpec;

  /** status */
  status?: LbcfBGStatus;
}

/** metadata的配置，非全部选项 */
interface LbcfBGMetadata {
  annotations?: {
    [props: string]: string;
  };

  /** ingress的名称 */
  name?: string;

  /** ingress的命名空间 */
  namespace?: string;
}

/** spec的配置项，非全部选项 */
interface LbcfBGSpec {
  lbName: string;
  pods?: {
    port: {
      portNumber: number;
      protocol: string;
    };
    byLabel: {
      selector: {
        [props: string]: string;
      };
    };
    byName: string[];
  };
  service?: {
    name: string;
    port: {
      portNumber: number;
      protocol: string;
    };
    nodeSelector: {
      [props: string]: string;
    };
  };
  static?: string[];
}
/** status的配置项，非全部选项 */
interface LbcfBGStatus {}

export interface GameBackgroupEdition extends Identifiable {
  onEdit: boolean;

  name?: string;

  v_name?: Validation;

  backgroupType?: string;

  //static
  staticAddress?: StringArray[];

  //pod
  byName?: StringArray[];

  //service
  serviceName?: string;

  v_serviceName?: Validation;

  //commom
  labels?: Selector[];

  ports?: Port[];
}
