import { satisfyClusterVersion } from '../../../helpers';
import { BasicEntryType, GroupEnumType, K8SKindNameEnum, VersionEnumType } from './apiKind';
import { apiServerVersion } from './apiServerVersion';

export type ApiVersionKeyName = keyof ApiVersion;
export interface ApiVersion {
  limitRange?: ResourceApiInfo;
  resourceQuota?: ResourceApiInfo;
  deployment?: ResourceApiInfo;
  statefulset?: ResourceApiInfo;
  daemonset?: ResourceApiInfo;
  job?: ResourceApiInfo;
  cronjob?: ResourceApiInfo;
  pods?: ResourceApiInfo;
  rc?: ResourceApiInfo;
  rs?: ResourceApiInfo;
  svc?: ResourceApiInfo;
  ingress?: ResourceApiInfo;
  configmap?: ResourceApiInfo;
  secret?: ResourceApiInfo;
  pv?: ResourceApiInfo;
  pvc?: ResourceApiInfo;
  sc?: ResourceApiInfo;
  hpa?: ResourceApiInfo;
  event?: ResourceApiInfo;
  node?: ResourceApiInfo;
  masteretcd?: ResourceApiInfo;
  cluster?: ResourceApiInfo;
  ns?: ResourceApiInfo;
  np?: ResourceApiInfo;
  logcs?: ResourceApiInfo;
  eniipamds?: ResourceApiInfo;
  loglistener?: ResourceApiInfo;
  hpc?: ResourceApiInfo;
  vpa?: ResourceApiInfo;
  endpoint?: ResourceApiInfo;
  injectpod?: ResourceApiInfo;

  // 集群備份相關
  backup?: ResourceApiInfo;
  backupSchedule?: ResourceApiInfo;
  restore?: ResourceApiInfo;

  /** rbac的相关资源 */
  clusterRole?: ResourceApiInfo;
  clusterRoleBinding?: ResourceApiInfo;
  role?: ResourceApiInfo;
  roleBinding?: ResourceApiInfo;

  /** 服务网格资源 */
  serviceForMesh?: ResourceApiInfo;
  gateway?: ResourceApiInfo;
  requestAuthentication?: ResourceApiInfo;
  peerAuthentication?: ResourceApiInfo;
  authentication?: ResourceApiInfo;
  authorization?: ResourceApiInfo;
  serviceaccount?: ResourceApiInfo;
  virtualservice?: ResourceApiInfo;
  destinationrule?: ResourceApiInfo;
  serviceentry?: ResourceApiInfo;
  controlPlane?: ResourceApiInfo;
  istioResource?: ResourceApiInfo;
  sidecar?: ResourceApiInfo;
  workloadEntry?: ResourceApiInfo;
  envoyFilter?: ResourceApiInfo;
  workloadGroup?: ResourceApiInfo;
  authorizationPolicy?: ResourceApiInfo;
  nginx_ingress_crd?: ResourceApiInfo;

  loadBalancerResource?: ResourceApiInfo;

  /** ============= 以下是ai environment的相关配置 ============== */
  ai_environment?: ResourceApiInfo;
  /** ============= 以下是ai environment的相关配置 ============== */

  /** ============= 以下是addon的相关配置 ============== */
  addon?: ResourceApiInfo;
  addonApps?: ResourceApiInfo;
  clusterAddon?: ResourceApiInfo;
  addon_helm?: ResourceApiInfo;
  addon_gameapp?: ResourceApiInfo;
  addon_gpumanager?: ResourceApiInfo;
  addon_logcollector?: ResourceApiInfo;
  addon_persistentevent?: ResourceApiInfo;
  addon_lbcf?: ResourceApiInfo;
  addon_cfs?: ResourceApiInfo;
  addon_cos?: ResourceApiInfo;
  addon_p2p?: ResourceApiInfo;
  addon_npd?: ResourceApiInfo;
  addon_dns_autoscaler?: ResourceApiInfo;
  addon_dns_localcache?: ResourceApiInfo;
  addon_tcr?: ResourceApiInfo;
  addon_oom_guard?: ResourceApiInfo;
  addon_nginx_ingress_controller_operator?: ResourceApiInfo;
  // addon_vertical_pod_autoscaler?: ResourceApiInfo;
  addon_scheduler?: ResourceApiInfo;
  addon_network_policy?: ResourceApiInfo;
  addon_hpc?: ResourceApiInfo;
  addon_cbs?: ResourceApiInfo;
  addon_de_scheduler?: ResourceApiInfo;
  addon_olm?: ResourceApiInfo;
  addon_qgpu?: ResourceApiInfo;
  addon_craned?: ResourceApiInfo;
  addon_crane_scheduler?: ResourceApiInfo;
  /** ============= 以上是addon的相关配置 ============== */

  /** ============= 以上是addon的相关配置 ============== */
  subscription?: ResourceApiInfo;
  manifest?: ResourceApiInfo;
  localization?: ResourceApiInfo;
  globalization?: ResourceApiInfo;
  description?: ResourceApiInfo;
  helmchart?: ResourceApiInfo;
  helmrelease?: ResourceApiInfo;
  /** ============= 以上是addon的相关配置 ============== */
  recommendationValue?: ResourceApiInfo;
  housekeeperPolicy?: ResourceApiInfo;
}

interface ResourceApiInfo {
  group: GroupEnumType;
  version: VersionEnumType;
  basicEntry: BasicEntryType;
  headTitle?: K8SKindNameEnum;
}

/** 这里的apiVersion的配置，只列出与1.8不同的，其余都保持和1.8相同的配置 */
const k8sApiVersionFor17: ApiVersion = {
  deployment: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Deployment
  },
  statefulset: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.StatefulSet
  },
  daemonset: {
    group: GroupEnumType.Extensions,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.DaemonSet
  },
  rs: {
    group: GroupEnumType.Extensions,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ReplicaSet
  },
  cronjob: {
    group: GroupEnumType.Batch,
    version: VersionEnumType.V2Alpha1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.CronJob
  },
  hpa: {
    group: GroupEnumType.AutoScaling,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.HorizontalPodAutoscaler
  },
  gateway: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Gateway
  },
  authentication: {
    group: GroupEnumType.SecurityIstioIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.PeerAuthentication
  },
  virtualservice: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.VirtualService
  },
  destinationrule: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.DestinationRule
  },
  serviceentry: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.None
  }
};

/**
 * 这里的apiVersion的配置，只列出与1.8不同的，其余都保持和1.8相同的配置
 * apps/v1beta1 及 apps/v1beta2 ，请使用apps/v1
 * extensions/v1beta1下的daemonsets , deployments , replicasets ，请使用apps/v1
 * extensions/v1beta1下的networkpolicies，请使用networking.k8s.io/v1
 * extensions/v1beta1下的podsecuritypolicies ，请使用policy/v1beta1
 */
const k8sApiVersionFor118: ApiVersion = {
  deployment: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Deployment
  },
  statefulset: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.StatefulSet
  },
  daemonset: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.DaemonSet
  },
  rs: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ReplicaSet
  },
  controlPlane: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ControlPlane
  }
};

/**
 * 这里的apiVersion的配置，只列出与1.8不同的，其余都保持和1.8相同的配置
 * 1.12及以上 extensions/v1beta1 => networking.k8s.io/v1beta1
 */
const k8sApiVersionForOver112: ApiVersion = {
  ingress: {
    group: GroupEnumType.NetworkingK8sIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Ingress
  }
};

const k8sApiVersionForOver122: ApiVersion = {
  clusterRole: {
    group: GroupEnumType.RbacAuthK8sIo,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ClusterRole
  },
  clusterRoleBinding: {
    group: GroupEnumType.RbacAuthK8sIo,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ClusterRoleBinding
  },
  role: {
    group: GroupEnumType.RbacAuthK8sIo,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Role
  },
  roleBinding: {
    group: GroupEnumType.RbacAuthK8sIo,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.RoleBinding
  },
  ingress: {
    group: GroupEnumType.NetworkingK8sIo,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Ingress
  }
};

const k8sApiVersionForOver126: ApiVersion = {
  hpa: {
    group: GroupEnumType.AutoScaling,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.HorizontalPodAutoscaler
  },
  cronjob: {
    group: GroupEnumType.Batch,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.CronJob
  }
};

/** 以1.8的为基准，后续有新增再继续更改 */
const k8sApiVersionFor18: ApiVersion = {
  deployment: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1Beta2,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Deployment
  },
  statefulset: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1Beta2,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.StatefulSet
  },
  daemonset: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1Beta2,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.DaemonSet
  },
  job: {
    group: GroupEnumType.Batch,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Job
  },
  cronjob: {
    group: GroupEnumType.Batch,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.CronJob
  },
  pods: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.Pod
  },
  rc: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.ReplicationController
  },
  rs: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1Beta2,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ReplicaSet
  },
  svc: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.Service
  },
  serviceaccount: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.ServiceAccount
  },
  ingress: {
    group: GroupEnumType.Extensions,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Ingress
  },
  configmap: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.ConfigMap
  },
  secret: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.Secret
  },
  pv: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.PersistentVolume
  },
  pvc: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.PersistentVolumeClaim
  },
  sc: {
    group: GroupEnumType.StorageK8sIo,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.StorageClass
  },
  hpa: {
    group: GroupEnumType.AutoScaling,
    version: VersionEnumType.V2Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.HorizontalPodAutoscaler
  },
  hpc: {
    group: GroupEnumType.HPCTencentCom,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.HorizontalPodCronscaler
  },
  vpa: {
    group: GroupEnumType.AutoScalingK8sIo,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.VerticalPodAutoscaler
  },
  event: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api
  },
  node: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.Node
  },
  masteretcd: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.Node
  },
  cluster: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.Cluster
  },
  ns: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.Namespace
  },
  np: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.Namespace
  },
  logcs: {
    group: GroupEnumType.CcsQcloudTencentCom,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.LogCollector
  },
  loglistener: {
    group: GroupEnumType.ClsCloudTencentCom,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.LogListener
  },
  eniipamds: {
    group: apiServerVersion.group,
    basicEntry: apiServerVersion.basicUrl,
    version: apiServerVersion.version,
    headTitle: K8SKindNameEnum.Eniipamds
  },
  gateway: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Gateway
  },
  requestAuthentication: {
    group: GroupEnumType.SecurityIstioIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.RequestAuthentication
  },
  peerAuthentication: {
    group: GroupEnumType.SecurityIstioIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.PeerAuthentication
  },
  authorization: {
    group: GroupEnumType.SecurityIstioIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Authorization
  },
  virtualservice: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.VirtualService
  },
  destinationrule: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.DestinationRule
  },
  serviceentry: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.None
  },
  sidecar: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Sidecar
  },
  workloadEntry: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.WorkloadEntry
  },
  envoyFilter: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.EnvoyFilter
  },
  workloadGroup: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.WorkloadGroup
  },
  authorizationPolicy: {
    group: GroupEnumType.SecurityIstioIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.AuthorizationPolicy
  },
  istioResource: {
    group: GroupEnumType.NetworkingIstioIo,
    version: VersionEnumType.V1Alpha3,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.IstioResource
  },
  controlPlane: {
    group: GroupEnumType.Apps,
    version: VersionEnumType.V1Beta2,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ControlPlane
  },
  clusterRole: {
    group: GroupEnumType.RbacAuthK8sIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ClusterRole
  },
  clusterRoleBinding: {
    group: GroupEnumType.RbacAuthK8sIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ClusterRoleBinding
  },
  role: {
    group: GroupEnumType.RbacAuthK8sIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Role
  },
  roleBinding: {
    group: GroupEnumType.RbacAuthK8sIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.RoleBinding
  },
  resourceQuota: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.ResourceQuota
  },
  limitRange: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.LimitRange
  },
  recommendationValue: {
    group: GroupEnumType.RecommendationTkeIo,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.RecommendationValue
  },
  ai_environment: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.Environment
  },
  loadBalancerResource: {
    group: GroupEnumType.NetworkingTkeCloudTencent,
    version: VersionEnumType.V1Aplha1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.RecommendationValue
  },
  backup: {
    group: GroupEnumType.BackupTKE,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Backup
  },
  backupSchedule: {
    group: GroupEnumType.BackupTKE,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.BackupSchedule
  },
  restore: {
    group: GroupEnumType.BackupTKE,
    version: VersionEnumType.V1Beta1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Restore
  },
  endpoint: {
    group: GroupEnumType.Empty,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Api,
    headTitle: K8SKindNameEnum.Endpoints
  },
  injectpod: {
    group: GroupEnumType.EKSCloudTencentCom,
    basicEntry: BasicEntryType.Apis,
    version: VersionEnumType.V1,
    headTitle: K8SKindNameEnum.Eniipamds
  }
};

const k8sApiVersionForAddon: ApiVersion = {
  addon: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.Addon
  },
  clusterAddon: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.ClusterAddon
  },
  addonApps: {
    group: GroupEnumType.ApplicationTkeStackIo,
    version: VersionEnumType.V1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.ClusterAddon
  },
  addon_helm: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.Helm
  },
  addon_gameapp: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.GameApp
  },
  addon_gpumanager: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.GPUManager
  },
  addon_logcollector: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.AddonLogCollector
  },
  addon_persistentevent: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.PersistentEvent
  },
  addon_lbcf: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.LBCF
  },
  addon_cfs: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.CFS
  },
  addon_cos: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.COS
  },
  addon_p2p: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.ImageP2P
  },
  addon_npd: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.NodeProblemDetector
  },
  addon_dns_autoscaler: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.DNSAutoscaler
  },
  addon_dns_localcache: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.NodeLocalDNSCache
  },
  addon_tcr: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.Tcr
  },
  addon_oom_guard: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.OOMGuard
  },
  addon_nginx_ingress_controller_operator: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.NginxIngressOperator
  },
  nginx_ingress_crd: {
    group: apiServerVersion.group,
    version: VersionEnumType.V1Aplha1,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.NginxIngressCRD
  },
  addon_scheduler: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.Scheduler
  },
  addon_network_policy: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.NetworkPolicy
  },
  addon_hpc: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.HPC
  },
  addon_cbs: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.CBS
  },
  addon_de_scheduler: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.DeScheduler
  },
  addon_olm: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.OLM
  },
  addon_qgpu: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.QGPU
  },
  addon_craned: {
    group: GroupEnumType.HousekeeperPolicy,
    version: VersionEnumType.V1Aplha1,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.housekeeperPolicy
  },
  addon_crane_scheduler: {
    group: GroupEnumType.HousekeeperPolicy,
    version: VersionEnumType.V1Aplha1,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.housekeeperPolicy
  },
  housekeeperPolicy: {
    group: GroupEnumType.HousekeeperPolicy,
    version: VersionEnumType.V1Aplha1,
    basicEntry: apiServerVersion.basicUrl,
    headTitle: K8SKindNameEnum.housekeeperPolicy
  }
};

const k8sApiVersionForClusternet: ApiVersion = {
  subscription: {
    group: GroupEnumType.Clusternet,
    version: VersionEnumType.V1Aplha1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Subscription
  },
  manifest: {
    group: GroupEnumType.Clusternet,
    version: VersionEnumType.V1Aplha1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Manifest
  },
  localization: {
    group: GroupEnumType.Clusternet,
    version: VersionEnumType.V1Aplha1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Localization
  },
  globalization: {
    group: GroupEnumType.Clusternet,
    version: VersionEnumType.V1Aplha1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Globalization
  },

  helmchart: {
    group: GroupEnumType.Clusternet,
    version: VersionEnumType.V1Aplha1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.HelmChart
  },
  helmrelease: {
    group: GroupEnumType.Clusternet,
    version: VersionEnumType.V1Aplha1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.HelmRelease
  },
  description: {
    group: GroupEnumType.Clusternet,
    version: VersionEnumType.V1Aplha1,
    basicEntry: BasicEntryType.Apis,
    headTitle: K8SKindNameEnum.Description
  }
};
/**
 * 这里配置的是k8s各版本的的k8s资源的 group 和 version 的最新版本
 * @pre 只需要保证每个k8s大版本当中的 group 和 version使用最新的即可，会向下兼容
 */
export function apiVersion(clusterVersion = '1.8'): ApiVersion {
  let finalApiVersion: ApiVersion = { ...k8sApiVersionForClusternet, ...k8sApiVersionForAddon, ...k8sApiVersionFor18 };

  // 1.7当中的group变化
  if (satisfyClusterVersion(clusterVersion, '1.7', 'eq')) {
    finalApiVersion = { ...finalApiVersion, ...k8sApiVersionFor17 };
  }

  // 1.12及其以上的版本，ingress的group 为 extensions/v1beta1 => networking.k8s.io/v1beta1
  if (satisfyClusterVersion(clusterVersion, '1.12', 'gt')) {
    finalApiVersion = { ...finalApiVersion, ...k8sApiVersionForOver112 };
  }

  // 1.18版本的group变化
  if (satisfyClusterVersion(clusterVersion, '1.18', 'ge')) {
    finalApiVersion = { ...finalApiVersion, ...k8sApiVersionFor118 };
  }

  if (satisfyClusterVersion(clusterVersion, '1.22', 'ge')) {
    finalApiVersion = { ...finalApiVersion, ...k8sApiVersionForOver122 };
  }

  if (satisfyClusterVersion(clusterVersion, '1.26', 'ge')) {
    finalApiVersion = { ...finalApiVersion, ...k8sApiVersionForOver126 };
  }

  return finalApiVersion;
}

//TODO 对照修复
