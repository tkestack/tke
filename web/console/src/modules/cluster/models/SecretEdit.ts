import { Identifiable, RecordSet } from '@tencent/qcloud-lib';
import { Validation } from '../../common/models';
import { Namespace } from './Namespace';
import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { QueryState } from '@tencent/qcloud-redux-query';
import { ResourceFilter } from './ResourceOption';

export interface SecretEdit extends Identifiable {
  /** secret名称 */
  name?: string;
  v_name?: Validation;

  /** namespace列表 */
  nsList?: FetcherState<RecordSet<Namespace>>;

  /** namepsace列表的查询 */
  nsQuery?: QueryState<ResourceFilter>;

  /** secret类型 */
  secretType?: string;

  /** secret的数据 */
  data?: SecretData[];

  /** ns的类型，是全部命名空间 还是 指定命名空间 */
  nsType?: string;

  /** 添加第三方镜像仓库的命名空间 */
  nsListSelection?: Namespace[];

  /** 当前填写的第三方镜像仓库的域名 */
  domain?: string;
  v_domain?: Validation;

  /** 第三方镜像仓库的用户名 */
  username?: string;
  v_username?: Validation;

  /** 第三方镜像仓库的密码 */
  password?: string;
  v_password?: Validation;
}

/** secret的数据类型 */
export interface SecretData extends Identifiable {
  /** key名称 */
  keyName?: string;
  v_keyName?: Validation;

  /** value名称 */
  value?: string;
  v_value?: Validation;
}

export interface SecretEditJSONYaml {
  /** 资源的类型 */
  kind?: string;

  /** api的版本 */
  apiVersion?: string;

  /** metadata */
  metadata?: SecretMetadata;

  /** data */
  data: {
    [props: string]: string;
  };

  type?: string;
}

interface SecretMetadata {
  /** 插件能力 */
  annotations?: {
    [props: string]: string;
  };

  /** pvc的名称 */
  name: string;

  /** pvc的命名空间 */
  namespace?: string;

  /** labels */
  labels?: {
    [props: string]: string;
  };
}
