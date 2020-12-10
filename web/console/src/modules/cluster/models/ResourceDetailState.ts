import { Validation, LogAgent } from 'src/modules/common';

import { FetcherState, FFListModel, QueryState, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { CreateResource, Event, Pod, Replicaset, ResourceFilter } from './';
import { PodFilterInNode } from './Pod';
import { Resource } from './ResourceOption';

type ResourceModifyWorkflow = WorkflowState<CreateResource, number>;

export interface ResourceDetailState {
  /** 拉取资源的详情信息 */
  resourceDetailInfo?: FFListModel<Resource, ResourceFilter>;

  /** yaml 的数据列表 */
  yamlList?: FetcherState<RecordSet<string>>;

  /** event的 FFRedux 列表 */
  event?: FFListModel<Event, ResourceFilter>;

  /** rs修订版本 */
  rsQuery?: QueryState<ResourceFilter>;

  /** rs列表 */
  rsList?: FetcherState<RecordSet<Replicaset>>;

  /** rsSelection */
  rsSelection?: Replicaset[];

  /** 回滚操作的工作流 */
  rollbackResourceFlow?: ResourceModifyWorkflow;

  /** pod的查询 */
  podQuery?: QueryState<ResourceFilter>;

  /** pod的列表 */
  podList?: FetcherState<RecordSet<Pod>>;

  /** node详情页内的pod列表的过滤业务 */
  podFilterInNode?: PodFilterInNode;

  /** container 列表 */
  containerList?: any[];

  /** podSelection */
  podSelection?: Pod[];

  /** 删除pod操作流 */
  deletePodFlow?: ResourceModifyWorkflow;

  /** 删除tapp pod 操作流 */
  removeTappPodFlow?: ResourceModifyWorkflow;

  /**tapp 灰度升级操作流 */
  updateGrayTappFlow?: ResourceModifyWorkflow;

  /**tapp 灰度升级编辑项 */
  editTappGrayUpdate?: TappGrayUpdateEditItem;

  /** 是否展示 登录弹框 */
  isShowLoginDialog?: boolean;

  /** log的查询 */
  logQuery?: QueryState<PodLogFilter>;

  /** log的列表 */
  logList?: FetcherState<RecordSet<string>>;

  /** logOption 用于日志的选择过滤条件 */
  logOption?: LogOption;

  logAgent?: LogAgent;

  logHierarchy?: string[];

  logContent?: string;
}

export interface LogOption {
  /** podName */
  podName?: string;

  /** containerName */
  containerName?: string;

  /** logFile */
  logFile?: string;

  /** tailLines */
  tailLines?: string;

  /** 是否开启自动刷新 */
  isAutoRenew?: boolean;
}

/**
 * 查询日志目录结构的 payload 定义
 */
export interface LogHierarchyQuery {
  agentName?: string;

  namespace?: string;

  clusterId?: string;

  pod?: string;

  container?: string;
}

/**
 * 查询日志内容的 payload 定义
 */
export interface LogContentQuery {
  agentName?: string;

  namespace?: string;

  clusterId?: string;

  pod?: string;

  container?: string;

  start?: number;

  length?: number;

  filepath?: string;
}

/**
 * 下载日志文件的 payload 定义
 */
export interface DownloadLogQuery {
  agentName?: string;

  namespace?: string;

  clusterId?: string;

  pod?: string;

  container?: string;

  filepath?: string;
}

export interface PodLogFilter extends ResourceFilter {
  /** container的名称 */
  container?: string;

  /** 日志文件路径 */
  logFile?: string;

  /** 显示日志的条数 */
  tailLines?: string;
}

export interface RsEditJSONYaml {
  /** 资源的类型 */
  kind: string;

  /** api的版本 */
  apiVersion: string;

  /** name: deployment的名字 */
  name?: string;

  /** 回滚到哪个版本 */
  rollbackTo?: {
    revision: number;
  };

  /** updatedAnnotations */
  updatedAnnotations?: any;
}
export interface TappGrayUpdateEditItem {
  /** 实例名称 */
  // name: string;

  // generateName: string;
  /** 容器 */
  containers: {
    /**容器名称 */
    // name: string;
    /**容器镜像名称 */
    imageName: string;
    /**容器镜像版本 */
    imageTag: string;
    /**
     * 校验选项
     */
    v_imageName: Validation;

    [props: string]: any;
  }[];
}
