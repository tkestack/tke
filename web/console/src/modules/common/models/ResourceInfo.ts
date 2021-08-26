/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

/** resourceConfig当中所需要定义的字段 */
export interface ResourceInfo {
  k8sVersion?: string;

  /** resource资源的名称 */
  headTitle?: string;

  /** 基本的api入口，api还是apis */
  basicEntry: string;

  /** 资源属于哪个group */
  group: string;

  /** api的版本 */
  version: string;

  /** 资源是否与 namespace 有关 */
  namespaces: string;

  /** 列表展示的数据 metadata.name 等 */
  displayField?: DisplayField;

  /** 请求的数据类型，list的话 是 deployments，单个的话又不同 */
  requestType: RequestType;

  /** 操作栏能够进行的操作 */
  actionField?: ActionField;

  /** 详情页面的相关配置 */
  detailField?: DetailField;
}

export interface DetailField {
  /** tabList */
  tabList?: TabItem[];

  /** detailInfo */
  detailInfo?: DetailInfo;
}

export interface DetailInfo {
  /** info，基本信息的展示区域 */
  info?: {
    [props: string]: DetailInfoProps;
  };

  /** 数据卷的展示区域 */
  volume?: {
    [props: string]: DetailInfoProps;
  };

  /** 容器的展示区域 */
  container?: {
    [props: string]: DetailInfoProps;
  };

  /** 转发规则的展示区域 */
  rules?: {
    [props: string]: DetailInfoProps;
  };

  /** 高级设置展示区域 */
  advancedInfo?: {
    [props: string]: DetailInfoProps;
  };

  /** LBCF */
  backGroup?: {
    [props: string]: DetailInfoProps;
  };
}

export interface DetailInfoProps {
  /** 需要展示的spec当中的资源 */
  dataField: string[];

  /** displaycField */
  displayField?: DetailDisplayField;
}

export interface DetailDisplayFieldProps {
  /** 需要展示的resourceIno的资源字段 */
  dataField: string[];

  /** 资源展示的类型，如 text, labels */
  dataFormat?: string;

  /** 该listitem的标题头 */
  label?: string;

  /** Listitem 的提示信息 */
  tips?: string;

  /** 是否为链接形式 */
  isLink?: boolean;

  /** 如果不存在该字段，使用默认值展示 */
  noExsitedValue?: string;

  /** subDisplay，内部是否还需要展示其他信息 */
  subDisplayField?: DetailDisplayField;

  /** extraInfo，用于展示数据后面是否增加 */
  extraInfo?: string;

  /** 映射值的展示 */
  mapTextConfig?: any;

  /** 展示的顺序，从0，5，10，15，20，……决定展示的顺序 */
  order?: string;
}

interface DetailDisplayField {
  [props: string]: DetailDisplayFieldProps;
}

interface TabItem {
  /** tab的id */
  id: string;

  /** tab 的名称 */
  label: string;
}

export interface ActionField {
  /** 创建按钮 */
  create: ActionItemField;

  /** 搜索框的配置 */
  search?: ActionItemField;

  /** 手动刷新 */
  manualRenew?: ActionItemField;

  /** 自动刷新 */
  autoRenew?: ActionItemField;

  /** 下载按钮 */
  download?: ActionItemField;
}

export interface ActionItemField {
  /** 是否提供该按钮 */
  isAvailable: boolean;

  /** 额外的一些属性 */
  attributes?: any[];
}

export interface DisplayField {
  [props: string]: DisplayFiledProps;
}

/**
 * 展示每个数据的具体定义
 */
export interface DisplayFiledProps {
  /** 需要展示的resourceIno的资源字段 */
  dataField: string[];

  /** 资源展示的类型，如 text, labels */
  dataFormat: string;

  /** 该项所占表格的宽度 */
  width: string;

  /** 该单元格的headeTitle */
  headTitle: string;

  /** 该单元格的headCell */
  headCell?: string[];

  /** 单元格的提示信息 */
  tips?: any;

  /** 如果不存在该字段，使用默认值展示 */
  noExsitedValue?: string;

  /** 是否需要跳转 */
  isLink?: boolean;

  /** 是否提复制值功能 */
  isClip?: boolean;

  /** 映射值的展示 */
  mapTextConfig?: any;

  /** 操作列表 */
  operatorList?: OperatorProps[];
}

/** resourceList当中操作列表 */
export interface OperatorProps {
  /** 操作的名称 */
  name: string;

  /** actionType: modify | delete | ……  */
  actionType: string;

  /** 是否放在更多当中 */
  isInMoreOp: boolean;
}

/**
 * 请求的类型 list | info 等
 */
export interface RequestType {
  /** 拉取列表数据 */
  list?: string;

  /** 是否为 addon资源 */
  addon?: boolean;

  /** 是否使用detailResourceInfo */
  useDetailInfo?: boolean;

  /** 必须提供info的 */
  detailInfoList?: {
    [props: string]: any[];
  };
}
