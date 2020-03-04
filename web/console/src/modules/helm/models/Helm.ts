import { Identifiable } from '@tencent/ff-redux';

export interface ClusterHelmStatus {
  code?: string;
  reason?: string;
}

interface HelmStatus {
  code?: string;
  resources?: string;
  notes?: string;
}
interface HelmInfo {
  status?: HelmStatus;
  first_deployed?: string;
  last_deployed?: string;
  Description?: string;
}

interface HelmMaintainer {
  name?: string;
  email?: string;
}
interface HelmChartMetadata {
  //chart name
  name?: string;
  //resource
  repo?: string;
  //namespace
  chart_ns?: string;
  //version
  version?: string;

  home?: string;
  sources?: string[];
  description?: string;
  keywords?: string[];
  maintainers?: HelmMaintainer[];
  engine?: string;
  icon?: string;
  appVersion?: string;
}

export interface HelmResource extends Identifiable {
  name?: string;
  kind?: string;
  yaml?: string;
}

export interface HelmConfig extends Identifiable {
  key?: string;
  value?: string;
}
export interface Helm extends Identifiable {
  name?: string;
  info?: HelmInfo;
  chart_metadata?: HelmChartMetadata;
  config?: HelmConfig[];
  configYaml?: string;
  version?: number;
  namespace?: string;
  resources?: HelmResource[];
  valueYaml?: string;
}

interface HelmTemplate {
  name?: string;
  data?: string;
}

interface HelmFile {
  type_url?: string;
  value?: string;
}
interface HelmChart {
  metadata?: HelmChartMetadata;
  templates?: HelmTemplate[];
  values?: Object;
  files?: HelmFile[];
}
export interface HelmDetail extends Identifiable {
  name?: string;
  info?: HelmInfo;
  // chart?: HelmChart;
  // config?: Object;
  // manifest?: string;
  // version?: number;
  namespace?: string;
}

export interface HelmFilter {
  searchName?: string;
  status?: string;
  /**地域Id */
  regionId?: string | number;
  clusterId?: string;
}

export interface HelmHistoryFilter {
  helmName?: string;
  /**地域Id */
  regionId?: number;
  clusterId?: string;
}

export interface HelmHistory extends Identifiable {
  name?: string;
  info?: HelmInfo;
  chart?: HelmChart;
  config?: Object;
  manifest?: string;
  version?: number;
  namespace?: string;
}

export interface InstallingHelm extends Identifiable {
  name?: string;
  status?: number;
}

export interface InstallingHelmDetail {
  code?: number;
  message?: string;
}
