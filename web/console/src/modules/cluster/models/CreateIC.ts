import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface LabelsKeyValue {
  key?: string;
  value?: string;
}
export interface ICComponter {
  ipList?: string;
  ssh?: string;
  cidr?: string;
  role?: string;
  labels?: LabelsKeyValue[];
  authType?: string;
  username?: string;
  password?: string;
  privateKey?: string;
  passPhrase?: string;
  isEditing?: boolean;

  //添加节点时候复用了
  isGpu?: boolean;
}
export interface CreateIC extends Identifiable {
  /**集群名称 */
  name?: string;
  v_name?: Validation;

  k8sVersion?: string;

  networkDevice?: string;
  v_networkDevice?: Validation;

  cidr?: string;

  maxClusterServiceNum?: string;

  maxNodePodNum?: string;

  k8sVersionList?: any[];

  computerList?: ICComponter[];
  computerEdit?: ICComponter;

  vipAddress?: string;
  v_vipAddress?: Validation;

  vipPort?: string;
  v_vipPort?: Validation;

  vipType?: string;

  gpu?: boolean;

  gpuType?: string;

  merticsServer?: boolean;

  cilium?: string;

  networkMode?: string;
}
