import { ApiVersionKeyName } from './apiVersion';

/** CRD资源的 Kind 全称 */
export enum K8SKindNameEnum {
  // K8S原生资源
  Deployment = 'Deployment',
  StatefulSet = 'StatefulSet',
  DaemonSet = 'DaemonSet',
  Job = 'Job',
  CronJob = 'CronJob',
  Pod = 'Pod',
  Namespace = 'Namespace',
  ReplicationController = 'ReplicationController',
  ReplicaSet = 'ReplicaSet',
  Service = 'Service',
  ServiceAccount = 'ServiceAccount',
  Ingress = 'Ingress',
  ConfigMap = 'ConfigMap',
  Secret = 'Secret',
  PersistentVolume = 'PersistentVolume',
  PersistentVolumeClaim = 'PersistentVolumeClaim',
  StorageClass = 'StorageClass',
  HorizontalPodAutoscaler = 'HorizontalPodAutoscaler',
  HorizontalPodCronscaler = 'HorizontalPodCronscaler',
  VerticalPodAutoscaler = 'VerticalPodAutoscaler',
  Node = 'Node',

  // RBAC相关
  ClusterRole = 'ClusterRole',
  ClusterRoleBinding = 'ClusterRoleBinding',
  Role = 'Role',
  RoleBinding = 'RoleBinding',

  // CRD
  Cluster = 'Cluster',
  PersistentEvent = 'PersistentEvent',
  LogCollector = 'LogCollector',
  LogListener = 'loglistener',
  Eniipamds = 'Eniipamds',
  Gateway = 'Gateway',
  RequestAuthentication = 'RequestAuthentication',
  PeerAuthentication = 'PeerAuthentication',
  Authorization = 'Authorization',
  VirtualService = 'Virtual Service',
  DestinationRule = 'Destination Rule',
  None = 'None',
  ControlPlane = 'Control Plane',
  Sidecar = 'Sidecar',
  WorkloadEntry = 'Workload Entry',
  EnvoyFilter = 'Envoy Filter',
  WorkloadGroup = 'Workload Group',
  AuthorizationPolicy = 'Authorization Policy',
  IstioResource = 'Istio Resource',
  // AI 相关
  Environment = 'Environment',

  // Addon相关
  Addon = 'Addon',
  ClusterAddon = 'ClusterAddon',
  AddonLogCollector = 'AddonLogCollector',
  Helm = 'Helm',
  GameApp = 'GameApp',
  GPUManager = 'GPUManager',
  LBCF = 'LBCF',
  CFS = 'CFS',
  COS = 'COS',
  ImageP2P = 'ImageP2P',
  NodeProblemDetector = 'NodeProblemDetector',
  DNSAutoscaler = 'DNSAutoscaler',
  NodeLocalDNSCache = 'NodeLocalDNSCache',
  Tcr = 'Tcr',
  OOMGuard = 'OOMGuard',

  NginxIngressOperator = 'NginxIngress',
  NginxIngressCRD = 'NginxIngressCRD',
  ResourceQuota = 'ResourceQuota',
  LimitRange = 'LimitRange',
  Scheduler = 'DynamicScheduler',
  NetworkPolicy = 'NetworkPolicy',
  HPC = 'HPC',
  CBS = 'CBS',
  DeScheduler = 'DeScheduler',
  OLM = 'OLM',

  // 分布式云相关
  Localization = 'Localization',
  Globalization = 'Globalization',
  Description = 'Description',
  Manifest = 'Manifest',
  Subscription = 'Subscription',
  HelmChart = 'HelmChart',
  HelmRelease = 'HelmRelease',

  RecommendationValue = 'recommendationValue',
  QGPU = 'QGPU',
  housekeeperPolicy = 'housekeeperPolicy',
  Endpoints = 'Endpoints',
  InjectPod = 'InjectPod',
  'Backup' = 'Backup',
  BackupSchedule = 'Schedule',
  'Restore' = 'Restore'
}

/** 资源所归属的集合 */
export enum GroupEnumType {
  /** apps */
  Apps = 'apps',

  /** extensions */
  Extensions = 'extensions',

  /** batch */
  Batch = 'batch',

  /** autoscaling */
  AutoScaling = 'autoscaling',

  AutoScalingK8sIo = 'autoscaling.k8s.io',

  /** networking.istio.io */
  NetworkingIstioIo = 'networking.istio.io',

  /** security.istio.io */
  SecurityIstioIo = 'security.istio.io',

  /** networking.k8s.io */
  NetworkingK8sIo = 'networking.k8s.io',

  /** storage.k8s.io */
  StorageK8sIo = 'storage.k8s.io',

  /** ccs.cloud.tencent.com */
  CcsQcloudTencentCom = 'ccs.cloud.tencent.com',

  /** cls.cloud.tencent.com */
  ClsCloudTencentCom = 'cls.cloud.tencent.com',

  /** empty */
  Empty = '',

  /** rbac.authorization.k8s.io */
  RbacAuthK8sIo = 'rbac.authorization.k8s.io',

  /** HPC */
  HPCTencentCom = 'autoscaling.cloud.tencent.com',

  /** 分布式云 */
  Clusternet = 'apps.clusternet.io',

  /** 推薦 */
  RecommendationTkeIo = 'recommendation.tke.io',

  CloudTencent = 'cloud.tencent.com',

  /** 新版Addon */
  ApplicationTkeStackIo = 'application.tkestack.io',

  NetworkingTkeCloudTencent = 'networking.tke.cloud.tencent.com',
  HousekeeperPolicy = 'scheduling.crane.io',

  EKSCloudTencentCom = 'eks.tke.cloud.tencent.com',

  BackupTKE = 'backup.tke.cloud.tencent.com'
}

/** 资源所归属的版本 */
export enum VersionEnumType {
  /** v1 */
  V1 = 'v1',

  /** v1beta1 */
  V1Beta1 = 'v1beta1',

  /** v1beta2 */
  V1Beta2 = 'v1beta2',

  /** v2beta1 */
  V2Beta1 = 'v2beta1',

  /** v2alpha1 */
  V2Alpha1 = 'v2alpha1',

  /** v1alpha3 */
  V1Alpha3 = 'v1alpha3',

  V1Aplha1 = 'v1alpha1',

  V2 = 'v2'
}

/** 资源所归属的路径 */
export enum BasicEntryType {
  /** api */
  Api = 'api',

  /** apis */
  Apis = 'apis'
}

/** K8sKindName to ApiVersionKeyName */
export const KindNameToApiKeyName: { [name in Exclude<K8SKindNameEnum, 'None'>]: ApiVersionKeyName } = {
  [K8SKindNameEnum.LimitRange]: 'limitRange',
  [K8SKindNameEnum.Endpoints]: 'endpoint',
  [K8SKindNameEnum.InjectPod]: 'injectpod',
  [K8SKindNameEnum.ResourceQuota]: 'resourceQuota',
  [K8SKindNameEnum.Deployment]: 'deployment',
  [K8SKindNameEnum.StatefulSet]: 'statefulset',
  [K8SKindNameEnum.DaemonSet]: 'daemonset',
  [K8SKindNameEnum.Job]: 'job',
  [K8SKindNameEnum.CronJob]: 'cronjob',
  [K8SKindNameEnum.Pod]: 'pods',
  [K8SKindNameEnum.Namespace]: 'ns',
  [K8SKindNameEnum.ReplicationController]: 'rc',
  [K8SKindNameEnum.ReplicaSet]: 'rs',
  [K8SKindNameEnum.Service]: 'svc',
  [K8SKindNameEnum.ServiceAccount]: 'serviceaccount',
  [K8SKindNameEnum.Ingress]: 'ingress',
  [K8SKindNameEnum.ConfigMap]: 'configmap',
  [K8SKindNameEnum.Secret]: 'secret',
  [K8SKindNameEnum.PersistentVolume]: 'pv',
  [K8SKindNameEnum.PersistentVolumeClaim]: 'pvc',
  [K8SKindNameEnum.StorageClass]: 'sc',
  [K8SKindNameEnum.HorizontalPodAutoscaler]: 'hpa',
  [K8SKindNameEnum.HorizontalPodCronscaler]: 'hpc',
  [K8SKindNameEnum.VerticalPodAutoscaler]: 'vpa',
  [K8SKindNameEnum.Node]: 'node',
  [K8SKindNameEnum.ClusterRole]: 'clusterRole',
  [K8SKindNameEnum.ClusterRoleBinding]: 'clusterRoleBinding',
  [K8SKindNameEnum.Role]: 'role',
  [K8SKindNameEnum.RoleBinding]: 'roleBinding',
  [K8SKindNameEnum.Cluster]: 'cluster',
  [K8SKindNameEnum.PersistentEvent]: 'addon_persistentevent',
  [K8SKindNameEnum.LogCollector]: 'logcs',
  [K8SKindNameEnum.LogListener]: 'loglistener',
  [K8SKindNameEnum.Eniipamds]: 'eniipamds',
  [K8SKindNameEnum.Gateway]: 'gateway',
  [K8SKindNameEnum.RequestAuthentication]: 'requestAuthentication',
  [K8SKindNameEnum.PeerAuthentication]: 'peerAuthentication',
  [K8SKindNameEnum.Authorization]: 'authorization',
  [K8SKindNameEnum.VirtualService]: 'virtualservice',
  [K8SKindNameEnum.DestinationRule]: 'destinationrule',
  [K8SKindNameEnum.ControlPlane]: 'controlPlane',
  [K8SKindNameEnum.Sidecar]: 'sidecar',
  [K8SKindNameEnum.WorkloadEntry]: 'workloadEntry',
  [K8SKindNameEnum.EnvoyFilter]: 'envoyFilter',
  [K8SKindNameEnum.WorkloadGroup]: 'workloadGroup',
  [K8SKindNameEnum.AuthorizationPolicy]: 'authorizationPolicy',
  [K8SKindNameEnum.IstioResource]: 'istioResource',
  [K8SKindNameEnum.Addon]: 'addon',
  [K8SKindNameEnum.ClusterAddon]: 'clusterAddon',
  [K8SKindNameEnum.AddonLogCollector]: 'addon_logcollector',
  [K8SKindNameEnum.Helm]: 'addon_helm',
  [K8SKindNameEnum.GameApp]: 'addon_gameapp',
  [K8SKindNameEnum.GPUManager]: 'addon_gpumanager',
  [K8SKindNameEnum.LBCF]: 'addon_lbcf',
  [K8SKindNameEnum.CFS]: 'addon_cfs',
  [K8SKindNameEnum.COS]: 'addon_cos',
  [K8SKindNameEnum.ImageP2P]: 'addon_p2p',
  [K8SKindNameEnum.NodeProblemDetector]: 'addon_npd',
  [K8SKindNameEnum.DNSAutoscaler]: 'addon_dns_autoscaler',
  [K8SKindNameEnum.NodeLocalDNSCache]: 'addon_dns_localcache',
  [K8SKindNameEnum.Tcr]: 'addon_tcr',
  [K8SKindNameEnum.OOMGuard]: 'addon_oom_guard',
  [K8SKindNameEnum.NginxIngressOperator]: 'addon_nginx_ingress_controller_operator',
  [K8SKindNameEnum.NginxIngressCRD]: 'nginx_ingress_crd',
  [K8SKindNameEnum.Scheduler]: 'addon_scheduler',
  [K8SKindNameEnum.NetworkPolicy]: 'addon_network_policy',
  [K8SKindNameEnum.HPC]: 'addon_hpc',
  [K8SKindNameEnum.CBS]: 'addon_cbs',
  [K8SKindNameEnum.DeScheduler]: 'addon_de_scheduler',
  [K8SKindNameEnum.OLM]: 'addon_olm',
  [K8SKindNameEnum.Environment]: 'ai_environment',

  // 集群備份相關
  [K8SKindNameEnum.Backup]: 'backup',
  [K8SKindNameEnum.BackupSchedule]: 'backupSchedule',
  [K8SKindNameEnum.Restore]: 'restore',

  [K8SKindNameEnum.Subscription]: 'subscription',
  [K8SKindNameEnum.Manifest]: 'manifest',
  [K8SKindNameEnum.Localization]: 'localization',
  [K8SKindNameEnum.Globalization]: 'globalization',
  [K8SKindNameEnum.HelmChart]: 'helmchart',
  [K8SKindNameEnum.HelmRelease]: 'helmrelease',
  [K8SKindNameEnum.Description]: 'description',
  [K8SKindNameEnum.RecommendationValue]: 'recommendationValue',
  [K8SKindNameEnum.QGPU]: 'addon_qgpu',
  [K8SKindNameEnum.housekeeperPolicy]: 'housekeeperPolicy'
};
