import { Identifiable } from '@tencent/qcloud-lib';
import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { Record, Validation } from '../../common/models';
import { WorkflowState } from '@tencent/qcloud-redux-workflow';

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
  v_key?: string;
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
  issueURL?: string;
  v_issueURL?: Validation;
  clientID?: string;
  v_clientID?: Validation;
  caCert?: string;
  v_caCert?: Validation;

  //镜像仓库设置
  repoType?: 'tke' | 'thirdParty';
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

  //业务模块设置
  openBusiness?: boolean;

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
