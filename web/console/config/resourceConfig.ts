import {
  deployment,
  statefulset,
  daemonset,
  jobs,
  cronjobs,
  tapps,
  pods,
  rc,
  rs,
  svc,
  ingress,
  np,
  configmap,
  secret,
  pv,
  pvc,
  sc,
  event,
  hpa,
  node,
  pe,
  cluster,
  namespace,
  moduleConfig,
  info,
  projects,
  portal,
  namespaces,
  localidentity,
  machines,
  policy,
  category,
  logcs,
  addon,
  addon_gpumanager,
  addon_helm,
  addon_logcollector,
  addon_tappcontroller,
  addon_csioperator,
  addon_lbcf,
  addon_cronhpa,
  addon_coredns,
  addon_galaxy,
  addon_prometheus,
  addon_volumedecorator,
  addon_ipam,
  alarmPolicy,
  helm,
  notifyChannel,
  notifyTemplate,
  notifyMessage,
  notifyReceiver,
  notifyReceiverGroup,
  platforms,
  lbcf,
  lbcf_bg,
  lbcf_br,
  clustercredential,
  serviceForMesh,
  gateway,
  controlPlane,
  virtualService,
  destinationRule,
  prometheus,
  logoutConfig,
  apiKey,
  user,
  cronhpa,
  role,
  localgroup,
  group
} from './resource/k8sConfig';
import { serviceEntry } from './resource/k8sConfig/serviceEntry';
import { ResourceInfo } from '../src/modules/common/models';
import { ApiVersion } from './resource/common';

/** 获取正确的集群版本号 */
export const ResourceConfigVersionMap = (k8sVersion: string) => {
  let finalK8sVersion: string;

  let [marjor, minor] = k8sVersion.split('.');
  let minorVersion = +minor;

  if (minorVersion >= 8 && minorVersion <= 12) {
    finalK8sVersion = '1.8';
  } else if (minorVersion > 12) {
    finalK8sVersion = '1.14';
  } else {
    finalK8sVersion = '1.7';
  }

  return finalK8sVersion;
};

const getResourceConfig = (resourceFunc: (k8sVersion: string) => ResourceInfo, k8sVersion: string) => {
  let result = resourceFunc(k8sVersion);
  return Object.assign({}, result, { k8sVersion });
};

/** ResourceConfig的返回定义 */
export type ResourceConfigKey = { [key in keyof ApiVersion]: ResourceInfo };

export const resourceConfig = (k8sVersion: string = '1.16'): ResourceConfigKey => {
  let finalK8sVersion = ResourceConfigVersionMap(k8sVersion) || '1.16';

  return {
    deployment: getResourceConfig(deployment, finalK8sVersion),
    statefulset: getResourceConfig(statefulset, finalK8sVersion),
    daemonset: getResourceConfig(daemonset, finalK8sVersion),
    job: getResourceConfig(jobs, finalK8sVersion),
    cronjob: getResourceConfig(cronjobs, finalK8sVersion),
    tapp: getResourceConfig(tapps, finalK8sVersion),
    pods: getResourceConfig(pods, finalK8sVersion),
    rc: getResourceConfig(rc, finalK8sVersion),
    rs: getResourceConfig(rs, finalK8sVersion),
    svc: getResourceConfig(svc, finalK8sVersion),
    ingress: getResourceConfig(ingress, finalK8sVersion),
    np: getResourceConfig(np, finalK8sVersion),
    event: getResourceConfig(event, finalK8sVersion),
    configmap: getResourceConfig(configmap, finalK8sVersion),
    secret: getResourceConfig(secret, finalK8sVersion),
    pv: getResourceConfig(pv, finalK8sVersion),
    pvc: getResourceConfig(pvc, finalK8sVersion),
    sc: getResourceConfig(sc, finalK8sVersion),
    hpa: getResourceConfig(hpa, finalK8sVersion),
    cronhpa: getResourceConfig(cronhpa, finalK8sVersion),
    node: getResourceConfig(node, finalK8sVersion),
    masteretcd: getResourceConfig(node, finalK8sVersion),
    pe: getResourceConfig(pe, finalK8sVersion),
    cluster: getResourceConfig(cluster, finalK8sVersion),
    ns: getResourceConfig(namespace, finalK8sVersion),
    module: getResourceConfig(moduleConfig, finalK8sVersion),
    info: getResourceConfig(info, finalK8sVersion),
    machines: getResourceConfig(machines, finalK8sVersion),
    helm: getResourceConfig(helm, finalK8sVersion),
    logcs: getResourceConfig(logcs, finalK8sVersion),
    logout: getResourceConfig(logoutConfig, finalK8sVersion),
    clustercredential: getResourceConfig(clustercredential, finalK8sVersion),

    lbcf: getResourceConfig(lbcf, finalK8sVersion),
    lbcf_bg: getResourceConfig(lbcf_bg, finalK8sVersion),
    lbcf_br: getResourceConfig(lbcf_br, finalK8sVersion),

    /** =============== 这里是业务相关的 =============== */
    projects: getResourceConfig(projects, finalK8sVersion),
    portal: getResourceConfig(portal, finalK8sVersion),
    platforms: getResourceConfig(platforms, finalK8sVersion),
    namespaces: getResourceConfig(namespaces, finalK8sVersion),
    /** =============== 这里是业务相关的 =============== */

    /** =============== 这里是addon相关的 =============== */
    addon: getResourceConfig(addon, finalK8sVersion),
    addon_persistentevent: getResourceConfig(pe, finalK8sVersion),
    addon_logcollector: getResourceConfig(addon_logcollector, finalK8sVersion),
    addon_helm: getResourceConfig(addon_helm, finalK8sVersion),
    addon_gpumanager: getResourceConfig(addon_gpumanager, finalK8sVersion),
    addon_tappcontroller: getResourceConfig(addon_tappcontroller, finalK8sVersion),
    addon_csioperator: getResourceConfig(addon_csioperator, finalK8sVersion),
    addon_lbcf: getResourceConfig(addon_lbcf, finalK8sVersion),
    addon_cronhpa: getResourceConfig(addon_cronhpa, finalK8sVersion),
    addon_coredns: getResourceConfig(addon_coredns, finalK8sVersion),
    addon_galaxy: getResourceConfig(addon_galaxy, finalK8sVersion),
    addon_prometheus: getResourceConfig(addon_prometheus, finalK8sVersion),
    addon_volumedecorator: getResourceConfig(addon_volumedecorator, finalK8sVersion),
    addon_ipam: getResourceConfig(addon_ipam, finalK8sVersion),
    /** =============== 这里是addon相关的 =============== */

    /** =============== 这里是权限相关的 =============== */
    localidentity: getResourceConfig(localidentity, finalK8sVersion),
    policy: getResourceConfig(policy, finalK8sVersion),
    category: getResourceConfig(category, finalK8sVersion),
    apiKey: getResourceConfig(apiKey, finalK8sVersion),
    user: getResourceConfig(user, finalK8sVersion),
    role: getResourceConfig(role, finalK8sVersion),
    localgroup: getResourceConfig(localgroup, finalK8sVersion),
    group: getResourceConfig(group, finalK8sVersion),
    /** =============== 这里是权限相关的 =============== */

    /** 告警配置 */
    prometheus: getResourceConfig(prometheus, finalK8sVersion),
    alarmPolicy: getResourceConfig(alarmPolicy, finalK8sVersion),

    /** 告警通知 */
    channel: getResourceConfig(notifyChannel, finalK8sVersion),
    template: getResourceConfig(notifyTemplate, finalK8sVersion),
    message: getResourceConfig(notifyMessage, finalK8sVersion),
    receiver: getResourceConfig(notifyReceiver, finalK8sVersion),
    receiverGroup: getResourceConfig(notifyReceiverGroup, finalK8sVersion),

    serviceForMesh: getResourceConfig(serviceForMesh, finalK8sVersion),
    gateway: getResourceConfig(gateway, finalK8sVersion),
    controlPlane: getResourceConfig(controlPlane, finalK8sVersion),
    virtualservice: getResourceConfig(virtualService, finalK8sVersion),
    destinationrule: getResourceConfig(destinationRule, finalK8sVersion),
    serviceentry: getResourceConfig(serviceEntry, finalK8sVersion)
  };
};
