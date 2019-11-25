import { uuid } from '@tencent/qcloud-lib';
import { initValidation } from '../../common/models';

export const initMachine = {
  id: uuid(),
  status: 'editing',
  host: '',
  v_host: initValidation,
  port: '',
  v_port: initValidation,
  authWay: 'password',
  user: '',
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

export const initEdit = {
  //基本设置
  username: '',
  v_username: initValidation,
  password: '',
  v_password: initValidation,
  confirmPassword: '',
  v_confirmPassword: initValidation,

  //高可用设置
  haType: 'tke',
  haTkeVip: '',
  v_haTkeVip: initValidation,
  haThirdVip: '',
  v_haThirdVip: initValidation,

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
  issueURL: '',
  v_issueURL: initValidation,
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

  //业务模块设置
  openBusiness: true,

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
