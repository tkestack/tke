import { Computer } from './Computer';
import { Identifiable, RecordSet } from '@tencent/qcloud-lib';
import { FetcherState, FetchState } from '@tencent/qcloud-redux-fetcher';
import { QueryState } from '@tencent/qcloud-redux-query';
import { Validation } from '../../common/models';
import { VolumeItem, Resource, ResourceFilter, ConfigItems, ContainerItem } from '../models';

export interface WorkloadEdit extends Identifiable {
  /** workload name */
  workloadName?: string;
  v_workloadName?: Validation;

  /** description */
  description?: string;
  v_description?: Validation;

  /** labels */
  workloadLabels?: WorkloadLabel[];

  /** annotataions */
  workloadAnnotations?: WorkloadLabel[];

  /** namespace */
  namespace?: string;
  v_namespace?: Validation;

  /** workload的类型 */
  workloadType?: string;

  /** cron的执行策略 */
  cronSchedule?: string;
  v_cronSchedule?: Validation;

  /** Job的重复执行次数 */
  completion?: string;
  v_completion?: Validation;

  /** 并行执行次数 */
  parallelism?: string;
  v_parallelism?: Validation;

  /** 失败重启策略 */
  restartPolicy?: string;

  /**tapp的节点异常迁移策略 */
  nodeAbnormalMigratePolicy?: string;

  /** 数据卷 */
  volumes?: VolumeItem[];

  /** 是否所有的数据卷都已经被挂载 */
  isAllVolumeIsMounted?: boolean;

  /** 是否展示 云硬盘的选择dialog */
  isShowCbsDialog?: boolean;

  /** 是否展示 配置项的选择dialog */
  isShowConfigDialog?: boolean;

  /** 是否展示 新建pvc选项的dialog */
  isShowPvcDialog?: boolean;

  /** 是否展示 主机路径选项的dialog */
  isShowHostPathDialog?: boolean;

  /** 当前正在编辑的 volume */
  currentEditingVolumeId?: string;

  /** pvc列表的query */
  pvcQuery?: QueryState<ResourceFilter>;

  /** pvc列表的 */
  pvcList?: FetcherState<RecordSet<Resource>>;

  /** 配置项 和 secret 当中的相关信息 */
  configEdit?: ConfigEdit;

  /** 是否能够新增容器实例 */
  canAddContainer?: boolean;

  /** 运行容器 */
  containers?: ContainerItem[];

  /** 实例的更新类型 */
  scaleType?: string;

  /** 容器的数量 */
  containerNum?: string;

  /** 是否需要展示容器的数量 */
  isNeedContainerNum?: boolean;

  /** 自动调节的最小实例数 */
  minReplicas?: string;
  v_minReplicas?: Validation;

  /** 自动调节的最大实例数 */
  maxReplicas?: string;
  v_maxReplicas?: Validation;

  /** metrics */
  metrics?: HpaMetrics[];

  /** cronhpa metrics */
  cronMetrics?: CronMetrics[];

  /** 是否同时创建Service */
  isCreateService?: boolean;

  /** imagePullSecret列表 */
  imagePullSecrets?: ImagePullSecrets[];

  /** 亲和性调度指定节点 */
  nodeSelection?: Computer[];

  /**节点校验 */
  v_nodeSelection?: Validation;

  /**亲和性调度方式："node" 指定节点调度 "rule" 自定义规则 "unset" 不使用*/
  nodeAffinityType?: string;

  /**亲和性调度自定义规则 */
  nodeAffinityRule?: AffinityRule;

  /** 是否支持 gpu 白名单 和 集群的版本 > 1.8 */
  isCanUseGpu?: boolean;

  isCanUseTapp?: boolean;

  /** 是否使用gpumanager模式*/
  isCanUseGpuManager?: boolean;

  /** 网络模式 */
  networkType?: string;

  /** 浮动ip回收机制 */
  floatingIPReleasePolicy?: string;

  /**超售比 */
  oversoldRatio?: { [props: string]: string };

  /** ===================== start 下面是实例更新的相关 ====================== */

  hpaQuery?: QueryState<ResourceFilter>;

  hpaList?: FetcherState<RecordSet<Resource>>;

  /** ===================== end 实例更新的相关 ====================== */

  /** ===================== start 下面是滚动镜像更新的相关 ====================== */

  /** 资源的更新策略 */
  resourceUpdateType?: string;

  /** 资源更新间隔 */
  minReadySeconds?: string;
  v_minReadySeconds?: Validation;

  /** 滚动更新的策略 createPod | destroyPod */
  rollingUpdateStrategy?: string;

  /** 批量大小，用于设定更新策略当中，前两者设置maxSurge 或者 maxUnavailable */
  batchSize?: string;
  v_batchSize?: Validation;

  /** 最大更新数量 */
  maxSurge?: string;
  v_maxSurge?: Validation;

  /** 最大停止数量 */
  maxUnavailable?: string;
  v_maxUnavailable?: Validation;

  /** statefulset 的 partition */
  partition?: string;
  v_partition?: Validation;

  /** ===================== end 滚动镜像更新的相关 ====================== */
}

export interface ImagePullSecrets extends Identifiable {
  /** 名称 */
  secretName: string;
  v_secretName: Validation;
}

/** HpaMetrics */
export interface HpaMetrics extends Identifiable {
  /** 指标的名称 */
  type?: string;
  v_type?: Validation;

  /** 指标的值 */
  value?: string;
  v_value?: Validation;
}

/** CronMetrics */
export interface CronMetrics extends Identifiable {
  /** crontab */
  crontab: string;
  v_crontab: Validation;

  /** targetReplicas */
  targetReplicas: string;
  v_targetReplicas: Validation;
}

/** workload的label的相关配置 */
export interface WorkloadLabel extends Identifiable {
  /** label的key */
  labelKey?: string;
  v_labelKey?: Validation;

  /** label的value */
  labelValue?: string;
  v_labelValue?: Validation;
}

/** 配置项 和 secret 弹窗的相关编辑工作 */
interface ConfigEdit {
  /** configMap的查询 */
  configQuery?: QueryState<ResourceFilter>;

  /** configmap的列表 */
  configList?: FetcherState<RecordSet<Resource>>;

  /** secret的查询 */
  secretQuery?: QueryState<ResourceFilter>;

  /** secret的列表 */
  secretList?: FetcherState<RecordSet<Resource>>;

  /** 选项： allKey | specificKey */
  configItems?: ConfigItems[];

  /** 当前选择的configMap */
  configSelection?: Resource[];

  /** 当前的configMap的 key的选择类型 */
  keyType?: string;

  /** 当前configMapSelection下的 key */
  configKeys?: string[];
}

/** 创建workload的时候，提交的jsonSchema */
export interface WorkloadEditJSONYaml {
  /** 资源的类型 */
  kind?: string;

  /** api的版本 */
  apiVersion?: string;

  /** metadata */
  metadata?: WorkloadMetadata;

  /** spec */
  spec: WorkloadSpec;
}

/** metadata的配置，非全部选项 */
interface WorkloadMetadata {
  /** 插件能力 */
  annotations?: {
    [props: string]: string;
  };

  /** 集群名称 */
  clusterName?: string;

  /** labels */
  labels?: {
    [props: string]: string;
  };

  /** service的名称 */
  name?: string;

  /** service的命名空间 */
  namespace?: string;
}

interface WorkloadSpec {
  /** template的相关配置 */
  template?: SpecTemplate;

  [props: string]: any;
}

interface SpecTemplate {
  /** metadata */
  metadata?: WorkloadMetadata;

  /** spec */
  spec?: PodSpec;
}

interface PodSpec {
  /** 容器相关的 */
  containers?: any;

  /** volumes */
  volumes?: any;

  [props: string]: any;
}

/** hpa的 jsonSchema */
export interface HpaEditJSONYaml {
  /** 资源的类型 */
  kind?: string;

  /** api的版本 */
  apiVersion?: string;

  /** metadata */
  metadata?: WorkloadMetadata;

  /** spec */
  spec: HpaSpec;
}

interface HpaSpec {
  /** 最大实例数量 */
  maxReplicas: number;

  /** 最小实例数量 */
  minReplicas: number;

  /** metrics */
  metrics: MetricOption[];

  scaleTargetRef?: {
    apiVersion: string;
    kind: string;
    name: string;
  };
}

export interface MetricOption {
  /** pods */
  pods?: {
    metricName: string;
    targetAverageValue: number | string;
  };

  resource?: {
    name: string;
    targetAverageUtilization?: number;
    targetAverageValue?: number | string;
  };

  type: 'Object' | 'Pods' | 'Resource';
}

/**亲和性类型 */
export interface AffinityRule {
  requiredExecution: NodeSelectorTerms[];
  preferredExecution: PreferredSchedulingTerm[];
}

interface PreferredSchedulingTerm {
  preference: { matchExpressions: MatchExpressions[] };
  weight: number;
}

interface NodeSelectorTerms {
  matchExpressions: MatchExpressions[];
}

export interface MatchExpressions extends Identifiable {
  key: string;
  operator: string;
  values: string;
  v_key: Validation;
  v_values: Validation;
}
