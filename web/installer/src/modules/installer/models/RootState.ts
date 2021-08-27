/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import { FetcherState, Identifiable, WorkflowState } from '@tencent/ff-redux';

import { Record, Validation } from '../../common/models';

type ClusterEditWorkflow = WorkflowState<EditState, void>;

export interface RootState {
  step?: string;

  cluster?: FetcherState<Record<any>>;

  isVerified?: number;

  licenseConfig?: any;

  clusterProgress?: FetcherState<Record<any>>;

  editState?: EditState;

  createCluster?: ClusterEditWorkflow;
}

export interface Arg extends Identifiable {
  key: string;
  v_key?: Validation;
  value: string;
  v_value?: Validation;
}

export interface EditState extends Identifiable {
  //基本设置
  username?: string;
  v_username?: Validation;
  password?: string;
  v_password?: Validation;
  confirmPassword?: string;
  v_confirmPassword?: Validation;

  //高可用设置
  haType?: string;
  haTkeVip?: string;
  v_haTkeVip?: Validation;
  haThirdVip?: string;
  v_haThirdVip?: Validation;
  haThirdVipPort?: string;
  v_haThirdVipPort?: Validation;

  //集群设置
  networkDevice?: string;
  v_networkDevice?: Validation;
  gpuType?: string;
  machines?: Array<Machine>;
  cidr?: string;
  podNumLimit?: number;
  serviceNumLimit?: number;

  //自定义集群设置
  dockerExtraArgs?: Array<Arg>;
  kubeletExtraArgs?: Array<Arg>;
  apiServerExtraArgs?: Array<Arg>;
  controllerManagerExtraArgs?: Array<Arg>;
  schedulerExtraArgs?: Array<Arg>;

  //认证模块设置
  authType?: string;
  tenantID?: string;
  v_tenantID?: Validation;
  issuerURL?: string;
  v_issuerURL?: Validation;
  clientID?: string;
  v_clientID?: Validation;
  caCert?: string;
  v_caCert?: Validation;

  //镜像仓库设置
  repoType?: string;
  repoTenantID?: string;
  v_repoTenantID?: Validation;
  repoSuffix?: string;
  v_repoSuffix?: Validation;
  repoAddress?: string;
  v_repoAddress: Validation;
  repoUser?: string;
  v_repoUser?: Validation;
  repoPassword?: string;
  v_repoPassword?: Validation;
  repoNamespace?: string;
  v_repoNamespace?: Validation;
  application: boolean;

  //业务模块设置
  openBusiness?: boolean;
  openAudit?: boolean;
  auditEsUrl?: string;
  auditEsReserveDays?: number;
  v_auditEsUrl?: Validation;
  auditEsUsername?: string;
  v_auditEsUsername?: Validation;
  auditEsPassword?: string;
  v_auditEsPassword?: Validation;

  //监控模块设置
  monitorType?: string;
  esUrl?: string;
  v_esUrl?: Validation;
  esUsername?: string;
  v_esUsername?: Validation;
  esPassword?: string;
  v_esPassword?: Validation;
  influxDBUrl?: string;
  v_influxDBUrl?: Validation;
  influxDBUsername?: string;
  v_influxDBUsername?: Validation;
  influxDBPassword?: string;
  v_influxDBPassword?: Validation;

  // 控制台设置
  openConsole?: boolean;
  consoleDomain?: string;
  v_consoleDomain?: Validation;
  certType?: string;
  certificate?: string;
  v_certificate?: Validation;
  privateKey?: string;
  v_privateKey?: Validation;
}

export interface Machine extends Identifiable {
  status?: 'editing' | 'edited';
  host?: string;
  v_host?: Validation;
  port?: string;
  v_port?: Validation;
  authWay?: 'password' | 'cert';
  user?: string;
  v_user?: Validation;
  password?: string;
  v_password?: Validation;
  cert?: string;
  v_cert?: Validation;
}
