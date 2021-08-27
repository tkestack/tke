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
/** 公有云 | 私有云切换 */
export const SelectMode = 'SelectMode';
export const IsI18n = 'IsI18n';

/** ============================== start 集群创建 相关 =============================== */
export const ModifyClusterNameWorkflow = 'ModifyClusterNameWorkflow';
export const ChangeMode = 'ChangeMode';
export const IsShowTips = 'IsShowTips';
export const isI18n = 'isI18n';
/** ============================== end 集群创建 相关 =============================== */

/** ============================== start 集群列表 相关 =============================== */
export const UpdateDialogState = 'UpdateDialogState';
export const DeleteCluster = 'DeleteCluster';
export const CreateCluster = 'CreateCluster';
export const UpdateclusterCreationState = 'UpdateclusterCreationState';
export const ClearClusterCreation = 'ClearClusterCreation';
/** ============================== end 集群列表 相关 =============================== */

/** ============================== start namespace 相关 =============================== */
export const FetchNamespaceList = 'FetchNamespaceList';
export const QueryNamespaceList = 'QueryNamespaceList';
export const SelectNamespace = 'SelectNamespace';
export const InitNamespaceInfo = 'InitNamespaceInfo';
/** ============================== end namespace 相关 =============================== */

/** ============================== start 集群详情 相关 =============================== */
export const ClusterVersion = 'ClusterVersion';
export const QueryClusterInfo = 'QueryClusterInfo';
export const FetchClusterInfo = 'FetchClusterInfo';
export const ClearDetail = 'ClearDetail';
export const ClearResourceDetail = 'ClearResourceDetail';
export const FetchClustercredential = 'FetchClustercredential';
export const UpdateClusterAllocationRatioEdition = 'UpdateClusterAllocationRatioEdition';
export const UpdateClusterAllocationRatio = 'UpdateClusterAllocationRatio';
export const UpdateClusterToken = 'UpdateClusterToken';
/** ============================== end 集群详情 相关 =============================== */

/** ============================== start 节点 相关 =============================== */
export const DeleteComputer = 'DeleteComputer';
export const BatchUnScheduleComputer = 'BatchUnScheduleComputer';
export const BatchTurnOnSchedulingComputer = 'BatchTurnOnSchedulingComputer';
export const DrainComputer = 'DrainComputer';
export const UpdateNodeLabel = 'UpdateNodeLabel';
export const UpdateLabelEdition = 'UpdateLabelEdition';
export const UpdateNodeTaint = 'UpdateNodeTaint';
export const UpdateTaintEdition = 'UpdateTaintEdition';
export const FetchComputerPodList = 'FetchComputerPodList';
export const QueryComputerPodList = 'QueryComputerPodList';
export const IsShowMachine = 'IsShowMachine';
export const FetchDeleteMachineResouceIns = 'FetchDeleteMachineResouceIns';
export const ClearComputer = 'ClearComputer';
/** ============================== end 节点 相关 =============================== */

/** ============================== start Resource 相关 =============================== */
export const ClearResource = 'ClearResource';
export const SelectMultipleResource = 'SelectMultipleResource';
export const SelectDeleteResource = 'SelectDeleteResource';
export const InitResourceName = 'InitResourceName';
export const InitResourceInfo = 'InitResourceInfo';
export const InitDetailResourceInfo = 'InitDetailResourceInfo';
export const InitDetailResourceName = 'InitDetailResourceName';
export const SelectDetailResourceSelection = 'SelectDetailResourceSelection';
export const InitDetailResourceList = 'InitDetailResourceList';
export const SelectDetailDeleteResourceSelection = 'SelectDetailDeleteResourceSelection';
export const ApplyResource = 'ApplyResource';
export const ModifyResource = 'ModifyResource';
export const ModifyMultiResource = 'ModifyMultiResource';

export const DeleteResource = 'DeleteResource';
export const ModifyResourcePodNum = 'ModifyResourcePodNum';
export const UpdateResourcePart = 'UpdateResourcePart';
export const UpdateMultiResource = 'UpdateMultiResource';

export const ApplyDifferentInterfaceResource = ' ApplyDifferentInterfaceResource';
/** ============================== end Resource 相关 =============================== */

/** ============================== start subRouter 相关 =============================== */
export const QuerySubRouterList = 'QuerySubRouterList';
export const FetchSubRouterList = 'FetchSubRouterList';
/** ============================== end subRouter 相关 =============================== */

/** ============================== start subRoot 相关 =============================== */
export const ClearSubRoot = 'ClearSubRoot';
export const IsNeedFetchNamespace = 'IsNeedFetchNamespace';
export const IsNeedExistedLb = 'IsNeedExistedLb';
export const FetchClusterAddons = 'FetchClusterAddons';
/** ============================== end subRoot 相关 =============================== */

/** ============================== start resourcrDetail 相关 =============================== */
export const FetchYaml = 'FetchYaml';
export const IsAutoPollingEvent = 'IsAutoPollingEvent';
export const QueryRsList = 'QueryRsList';
export const FetchRsList = 'FetchRsList';
export const RollBackResource = 'RollBackResource';
export const RsSelection = 'RsSelection';
export const QueryPodList = 'QueryPodList';
export const FetchPodList = 'FetchPodList';
export const PodFilterInNode = 'PodFilterInNode';
export const FetchContainerList = 'FetchContainerList';
export const PodSelection = 'PodSelection';
export const DeletePod = 'DeletePod';
export const IsShowLoginDialog = 'IsShowLoginDialog';
export const QueryLogList = 'QueryLogList';
export const FetchLogList = 'FetchLogList';
export const PodLogAgent = 'PodLogAgent';
export const PodLogHierarchy = 'PodLogHierarchy';
export const PodLogContent = 'PodLogContent';
export const PodName = 'PodName';
export const ContainerName = 'ContainerName';
export const LogFile = 'LogFile';
export const TailLines = 'TailLines';
export const IsAutoRenewPodLog = 'IsAutoRenewPodLog';
export const QuerySecretList = 'QuerySecretList';
export const FetchSecretList = 'FetchSecretList';
export const SecretSelection = 'SecretSelection';
export const ModifyNamespaceSecret = 'ModifyNamespaceSecret';
export const RemoveTappPod = 'RemoveTappPod';
export const UpdateGrayTapp = 'UpdateGrayTapp';
export const W_TappGrayUpdate = 'W_TappGrayUpdate';
/** ============================== end resourcrDetail 相关 =============================== */

/** ============================== start 创建服务 Service 相关 =============================== */
export const ClearServiceEdit = 'ClearServiceEdit';
export const S_ServiceName = 'S_ServiceName';
export const SV_ServiceName = 'SV_ServiceName';
export const S_Description = 'S_Description';
export const SV_Description = 'SV_Description';
export const S_Namespace = 'S_Namespace';
export const SV_Namespace = 'SV_Namespace';
export const S_CommunicationType = 'S_CommunicationType';
export const S_UpdatePortsMap = 'S_UpdatePortsMap';
export const S_Selector = 'S_Selector';
export const S_IsOpenHeadless = 'S_IsOpenHeadless';
export const S_IsShowWorkloadDialog = 'S_IsShowWorkloadDialog';
export const S_WorkloadType = 'S_WorkloadType';
export const S_QueryWorkloadList = 'S_QueryWorkloadList';
export const S_FetchWorkloadList = 'S_FetchWorkloadList';
export const S_WorkloadSelection = 'S_WorkloadSelection';
export const S_ChooseExternalTrafficPolicy = 'S_ChooseExternalTrafficPolicy';
export const S_ChoosesessionAffinity = 'S_ChoosesessionAffinity';
export const S_InputsessionAffinityTimeout = 'S_InputsessionAffinityTimeout';
export const SV_sessionAffinityTimeout = 'SV_sessionAffinityTimeout';

/** ============================== end 创建服务 Service 相关 =============================== */

/** ============================== start 创建Namespace相关 =============================== */
export const ClearNamespaceEdit = 'ClearNamespaceEdit';
export const N_Name = 'N_Name';
export const NV_Name = 'NV_Name';
export const N_Description = 'N_Description';
export const NV_Description = 'NV_Description';
/** ============================== end 创建Namespace相关 =============================== */

/** ============================== start 创建Workload相关 =============================== */
export const ClearWorkloadEdit = 'ClearWorkloadEdit';
export const W_WorkloadName = 'W_WorkloadName';
export const WV_WorkloadName = 'WV_WorkloadName';
export const W_Description = 'W_Description';
export const WV_Description = 'WV_Description';
export const W_Namespace = 'W_Namespace';
export const WV_Namespace = 'WV_Namespace';
export const W_WorkloadType = 'W_WorkloadType';
export const W_UpdateVolumes = 'W_UpdateVolumes';
export const W_IsShowCbsDialog = 'W_IsShowCbsDialog';
export const W_IsShowConfigDialog = 'W_IsShowConfigDialog';
export const W_IsShowPvcDialog = 'W_IsShowPvcDialog';
export const W_IsShowHostPathDialog = 'W_IsShowHostPathDialog';
export const W_CurrentEditingVolumeId = 'W_CurrentEditingVolumeId';
export const W_QueryPvcList = 'W_QueryPvcList';
export const W_FetchPvcList = 'W_FetchPvcList';
export const W_QueryConfig = 'W_QueryConfig';
export const W_FetchConfig = 'W_FetchConfig';
export const W_QuerySecret = 'W_QuerySecret';
export const W_FetchSecret = 'W_FetchSecret';
export const W_UpdateConfigItems = 'W_UpdateConfigItems';
export const ClearWorkloadConfig = 'ClearWorkloadConfig';
export const W_ConfigSelect = 'W_ConfigSelect';
export const W_ChangeKeyType = 'W_ChangeKeyType';
export const W_UpdateConfigKeys = 'W_UpdateConfigKeys';
export const W_CanAddContainer = 'W_CanAddContainer';
export const W_IsCanUseGpu = 'W_IsCanUseGpu';
export const W_IsCanUseTapp = 'W_IsCanUseTapp';
export const W_NetworkType = 'W_NetworkType';
export const W_FloatingIPReleasePolicy = 'W_FloatingIPReleasePolicy';
export const W_UpdateContainers = 'W_UpdateContainers';
export const W_ChangeScaleType = 'W_ChangeScaleType';
export const W_IsOpenCronHpa = 'W_IsOpenCronHpa';
export const W_ContainerNum = 'W_ContainerNum';
export const W_IsAllVolumeIsMounted = 'W_IsAllVolumeIsMounted';
export const W_IsNeedContainerNum = 'W_IsNeedContainerNum';
export const W_IsCreateService = 'W_IsCreateService';
export const ImagePullSecrets = 'ImagePullSecrets';
export const W_CronSchedule = 'W_CronSchedule';
export const WV_CronSchedule = 'WV_CronSchedule';
export const W_Completion = 'W_Completion';
export const WV_Completion = 'WV_Completion';
export const W_Parallelism = 'W_Parallelism';
export const WV_Parallelism = 'WV_Parallelism';
export const W_RestartPolicy = 'W_RestartPolicy';
export const W_WorkloadLabels = 'W_WorkloadLabels';
export const W_WorkloadAnnotations = 'W_WorkloadAnnotations';
export const W_ResourceUpdateType = 'W_ResourceUpdateType';
export const W_MinReadySeconds = 'W_MinReadySeconds';
export const WV_MinReadySeconds = 'WV_MinReadySeconds';
export const W_RollingUpdateStrategy = 'W_RollingUpdateStrategy';
export const W_BatchSize = 'W_BatchSize';
export const WV_BatchSize = 'WV_BatchSize';
export const W_MaxSurge = 'W_MaxSurge';
export const WV_MaxSurge = 'WV_MaxSurge';
export const W_MaxUnavailable = 'W_MaxUnavailable';
export const WV_MaxUnavailable = 'WV_MaxUnavailable';
export const W_Partition = 'W_Partition';
export const WV_Partition = 'WV_Partition';
export const W_MinReplicas = 'W_MinReplicas';
export const WV_MinReplicas = 'WV_MinReplicas';
export const W_MaxReplicas = 'W_MaxReplicas';
export const WV_MaxReplicas = 'WV_MaxReplicas';
export const W_UpdateMetrics = 'W_UpdateMetrics';
export const W_UpdateCronMetrics = 'W_UpdateCronMetrics';
export const W_QueryHpaList = 'W_QueryHpaList';
export const W_FetchHpaList = 'W_FetchHpaList';
export const W_NodeAbnormalMigratePolicy = 'W_NodeAbnormalMigratePolicy';
export const W_UpdateOversoldRatio = 'W_UpdateOversoldRatio';

/**亲和性调度相关 */
export const WV_NodeSelector = 'WV_NodeSelector';
export const W_UpdateNodeAffinityRule = 'W_UpdateNodeAffinityRule';
export const W_SelectNodeAffinityType = 'W_SelectNodeAffinityType';
/** ============================== end 创建Workload相关 =============================== */

/** ============================== start 创建configMap相关 =============================== */
export const ClearConfigMapEdit = 'ClearConfigMapEdit';
export const CM_Name = 'CM_Name';
export const V_CM_Name = 'V_CM_Name';
export const CM_Namespace = 'CM_Namespace';
export const CM_AddVariable = 'CM_AddVariable';
export const CM_EditVariable = 'CM_EditVariable';
export const CM_DeleteVariable = 'CM_DeleteVariable';
export const CM_ValidateVariable = 'CM_ValidateVariable';
/** ============================== end 创建configMap相关 =============================== */

/** ============================== start 创建 服务日志相关声明 相关 =============================== */
export const ClearResourceLog = 'ClearResourceLog';
export const L_WorkloadType = 'L_WorkloadType';
export const L_NamespaceSelection = 'L_NamespaceSelection';
export const L_WorkloadSelection = 'L_WorkloadSelection';
export const L_QueryWorkloadList = 'L_QueryWorkloadList';
export const L_FetchWorkloadList = 'L_FetchWorkloadList';
export const L_QueryPodList = 'L_QueryPodList';
export const L_FetchPodList = 'L_FetchPodList';
export const L_PodSelection = 'L_PodSelection';
export const L_ContainerSelection = 'L_ContainerSelection';
export const L_QueryLogList = 'L_QueryLogList';
export const L_FetchLogList = 'L_FetchLogList';
export const L_TailLines = 'L_TailLines';
export const L_IsAutoRenew = 'L_IsAutoRenew';
/** ============================== end 创建 服务日志相关声明 相关 =============================== */

/** ============================== start 创建 服务事件相关声明 =============================== */
export const ClearResourceEvent = 'ClearResourceEvent';
export const E_WorkloadType = 'E_WorkloadType';
export const E_NamespaceSelection = 'E_NamespaceSelection';
export const E_QueryWorkloadList = 'E_QueryWorkloadList';
export const E_FetchWorkloadList = 'E_FetchWorkloadList';
export const E_WorkloadSelection = 'E_WorkloadSelection';
export const E_QueryEventList = 'E_QueryEventList';
export const E_FetchEventList = 'E_FetchEventList';
export const E_IsAutoRenew = 'E_IsAutoRenew';
/** ============================== end 创建 服务事件 相关声明 =============================== */

/** ============================== start 创建 Secret创建 相关 =============================== */
export const ClearSecretEdit = 'ClearSecretEdit';
export const Sec_Name = 'Sec_Name';
export const SecV_Name = 'SecV_Name';
export const Sec_FetchNsList = 'Sec_FetchNsList';
export const Sec_QueryNsList = 'Sec_QueryNsList';
export const Sec_SecretType = 'Sec_SecretType';
export const Sec_UpdateData = 'Sec_UpdateData';
export const Sec_NsType = 'Sec_NsType';
export const Sec_NamespaceSelection = 'Sec_NamespaceSelection';
export const Sec_Domain = 'Sec_Domain';
export const SecV_Domain = 'SecV_Domain';
export const Sec_Username = 'Sec_Username';
export const SecV_Username = 'SecV_Username';
export const Sec_Password = 'Sec_Password';
export const SecV_Password = 'SecV_Password';
/** ============================== end 创建 Secret创建 相关 =============================== */

/** 业务侧应用页面 */
export const InitClusterList = 'InitClusterList';
export const ClusterSelection = 'ClusterSelection';
export const InitProjectList = 'InitProjectList';
export const ProjectSelection = 'ProjectSelection';
export const FetchProjectNamespace = 'FetchProjectNamespace';
export const QueryProjectNamespace = 'QueryProjectNamespace';

/** ============================== start 创建独立集群相关 =============================== */
// IC = IndependentCluster
export const IC_Clear = 'IC_Clear';
export const IC_Name = 'IC_Name';
export const v_IC_Name = 'v_IC_Name';
export const v_IC_NetworkDevice = 'v_IC_NetworkDevice';
export const IC_NetworkDevice = 'IC_NetworkDevice';
export const IC_VipAddress = 'IC_VipAddress';
export const v_IC_VipAddress = 'v_IC_VipAddress';
export const IC_VipPort = 'IC_VipPort';
export const v_IC_VipPort = 'v_IC_VipPort';
export const v_IC_Vip = 'v_IC_Vip';
export const v_IC_Gpu = 'v_IC_Gpu';
export const v_IC_GpuType = 'v_IC_GpuType';
export const v_IC_Mertics_server = 'v_IC_Mertics_server';
export const v_IC_Cilium = 'v_IC_Cilium';
export const v_IC_NetworkMode = 'v_IC_NetworkMode';
export const v_IC_AS = 'v_IC_AS';
export const IC_AS = 'IC_AS';
export const v_IC_SwitchIp = 'v_IC_SwitchIp';
export const IC_SwitchIp = 'IC_SwitchIp';

export const IC_K8SVersion = 'IC_K8SVersion';
export const IC_Cidr = 'IC_Cidr';
export const IC_ComputerList = 'IC_ComputerList';
export const IC_ComputerEdit = 'IC_ComputerEdit';
export const CreateIC = 'CreateIC';
export const IC_FetchK8SVersion = 'IC_FetchK8SVersion';
export const IC_MaxClusterServiceNum = 'IC_MaxClusterServiceNum';
export const IC_MaxNodePodNum = 'IC_MaxNodePodNum';
export const IC_EnableContainerRuntime = 'IC_EnableContainerRuntime';
/** ============================== end 创建独立集群相关 =============================== */

/** ============================== start 新增节点相关 =============================== */
export const CreateComputer = 'CreateComputer';
/** ============================== end 新增节点相关  =============================== */

/** ============================== start 创建 Lbcf创建 相关 =============================== */
export const Gate_Name = 'Gate_Name';
export const V_Gate_Name = 'V_Gate_Name';
export const Gate_Namespace = 'Gate_Namespace';
export const V_Gate_Namespace = 'V_GLB_Namespace';

export const Lbcf_Config = 'Lbcf_Config';
export const Lbcf_Args = 'Lbcf_Args';
export const V_Lbcf_Config = 'V_Lbcf_Config';
export const V_Lbcf_Args = 'V_Lbcf_Args';
export const V_Lbcf_Driver = 'V_Lbcf_Driver';

export const GLB_VpcSelection = 'GLB_VpcSelection';
export const GLB_FecthClb = 'GLB_FecthClb';
export const GLB_SelectClb = 'GLB_SelectClb';
export const V_GLB_SelectClb = 'V_GLB_SelectClb';
export const GLB_SwitchCreateLbWay = 'GLB_SwitchCreateLbWay';

export const GBG_UpdateLbcfBackGroup = 'GBG_UpdateLbcfBackGroup';
export const GBG_FetchGameApp = 'GBG_FetchGameApp';
export const GBG_SelectGameApp = 'GBG_SelectGameApp';
export const GBG_ShowGameAppDialog = 'GBG_ShowGameAppDialog';
export const ClearLbcfEdit = 'ClearGBGEdit';
/** ============================== end 创建 Lbcf创建 相关 =============================== */
