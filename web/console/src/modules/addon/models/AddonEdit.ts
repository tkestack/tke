import { PeEdit } from './PeEdit';

export interface AddonEdit {
  /** 当前选择的组件 */
  addonName?: string;

  /** 事件持久化的编辑项 */
  peEdit?: PeEdit;
}

interface AddonEditBasicJsonYaml {
  /** 资源的类型 */
  kind: string;

  /** api的版本 */
  apiVersion: string;

  /** metadata */
  metadata?: any;

  /** spec */
  spec?: any;
}

/** ====================== Helm、GameApp 创建相关的yaml ====================== */
export interface AddonEditUniversalJsonYaml extends AddonEditBasicJsonYaml {
  metadata: {
    generateName: string;
  };

  spec: {
    clusterName: string;
  };
}
/** ====================== Helm、GameApp 创建相关的yaml ====================== */

/** ====================== persistentEvent创建相关的yaml ===================== */
export interface AddonEditPeJsonYaml extends AddonEditBasicJsonYaml {
  metadata: {
    generateName: string;
  };

  spec: {
    clusterName: string;
    persistentBackEnd: PersistentBackEnd;
  };
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
/** ====================== persistentEvent创建相关的yaml ===================== */
