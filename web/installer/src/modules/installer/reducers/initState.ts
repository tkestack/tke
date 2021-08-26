/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { uuid } from '@tencent/ff-redux';

import { initValidation } from '../../common/models';
import { EditState } from '../models';

export const initMachine = {
  id: uuid(),
  status: 'editing',
  host: '',
  v_host: initValidation,
  port: '',
  v_port: initValidation,
  authWay: 'password',
  user: 'root',
  v_user: initValidation,
  password: '',
  v_password: initValidation,
  cert: '',
  v_cert: initValidation
};

export const initArg = {
  id: uuid(),
  key: '',
  v_key: initValidation,
  value: '',
  v_value: initValidation
};

export const initEdit: EditState = {
  id: '',
  //基本设置
  username: '',
  v_username: initValidation,
  password: '',
  v_password: initValidation,
  confirmPassword: '',
  v_confirmPassword: initValidation,

  //高可用设置
  haType: 'none',
  haTkeVip: '',
  v_haTkeVip: initValidation,
  haThirdVip: '',
  v_haThirdVip: initValidation,
  haThirdVipPort: '6443',
  v_haThirdVipPort: initValidation,

  //集群设置
  networkDevice: 'eth0',
  v_networkDevice: initValidation,
  gpuType: 'none',
  machines: [Object.assign({}, initMachine, { id: uuid() })],
  cidr: '192.168.0.0/16',
  podNumLimit: 256,
  serviceNumLimit: 256,

  //自定义集群设置
  dockerExtraArgs: [Object.assign({}, initArg, { id: uuid() })],
  kubeletExtraArgs: [Object.assign({}, initArg, { id: uuid() })],
  apiServerExtraArgs: [Object.assign({}, initArg, { id: uuid() })],
  controllerManagerExtraArgs: [Object.assign({}, initArg, { id: uuid() })],
  schedulerExtraArgs: [Object.assign({}, initArg, { id: uuid() })],

  //认证模块设置
  authType: 'tke',
  tenantID: '',
  v_tenantID: initValidation,
  issuerURL: '',
  v_issuerURL: initValidation,
  clientID: '',
  v_clientID: initValidation,
  caCert: '',
  v_caCert: initValidation,

  //镜像仓库设置
  repoType: 'tke',
  repoTenantID: '',
  v_repoTenantID: initValidation,
  repoSuffix: 'registry.tke.com',
  v_repoSuffix: initValidation,
  repoAddress: '',
  v_repoAddress: initValidation,
  repoUser: '',
  v_repoUser: initValidation,
  repoPassword: '',
  v_repoPassword: initValidation,
  repoNamespace: '',
  v_repoNamespace: initValidation,
  application: true,

  //业务模块设置
  openBusiness: true,
  openAudit: false,
  auditEsUrl: '',
  auditEsReserveDays: 7,
  v_auditEsUrl: initValidation,
  auditEsUsername: '',
  v_auditEsUsername: initValidation,
  auditEsPassword: '',
  v_auditEsPassword: initValidation,

  //监控模块设置
  monitorType: 'tke-influxdb',
  esUrl: '',
  v_esUrl: initValidation,
  esUsername: '',
  v_esUsername: initValidation,
  esPassword: '',
  v_esPassword: initValidation,
  influxDBUrl: '',
  v_influxDBUrl: initValidation,
  influxDBUsername: '',
  v_influxDBUsername: initValidation,
  influxDBPassword: '',
  v_influxDBPassword: initValidation,

  // 控制台设置
  openConsole: true,
  consoleDomain: 'console.tke.com',
  v_consoleDomain: initValidation,
  certType: 'selfSigned',
  certificate: '',
  v_certificate: initValidation,
  privateKey: '',
  v_privateKey: initValidation
};
