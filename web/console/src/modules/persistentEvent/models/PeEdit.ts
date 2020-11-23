import { Validation } from '../../common/models/Validation';

export interface PeEdit {
  /** 是否开启事件持久化 */
  isOpen?: boolean;

  /** es的地址 */
  esAddress?: string;
  v_esAddress?: Validation;

  /** 索引名称 */
  indexName?: string;
  v_indexName?: Validation;

  /** ES 认证用户名 */
  esUsername?: string;

  /** ES 认证密码 */
  esPassword?: string;
}

/** 编辑事件持久化的时候，提交的jsonSchema */
export interface PeEditJSONYaml {
  /** 资源的类型 */
  kind?: string;

  /** api的版本 */
  apiVersion?: string;

  /** metadata */
  metadata?: PeMetadata;

  /** spec */
  spec?: PeSpec;
}

interface PeMetadata {
  /** 插件能力 */
  annotations?: {
    [props: string]: string;
  };

  /** labels */
  labels?: {
    [props: string]: string;
  };

  /** service的名称 */
  name?: string;

  /** generateName */
  generateName?: string;

  [props: string]: any;
}

interface PeSpec {
  /** 集群的名称 */
  clusterName?: string;

  /** persistentBackEnd */
  persistentBackEnd: PersistentBackEnd;
}

export interface PersistentBackEnd {
  /** es的配置 */
  es: EsInfo;
}

export interface EsInfo {
  ip: string;
  port: number;
  scheme: string;
  indexName: string;
  user: string;
  password: string;
}
