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
import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';
import { ContainerEnv } from './ContainerEnv';

export interface ContainerItem extends Identifiable {
  /** 状态 */
  status?: string;

  /** 容器名称 */
  name?: string;
  v_name?: Validation;

  /** 镜像 */
  registry?: string;
  v_registry?: Validation;

  /** 镜像版本 */
  tag?: string;

  /** 挂载点 */
  mounts?: MountItem[];

  /** 内存推荐值 */
  memLimit: LimitItem[];

  /** 推荐值 */
  cpuLimit?: LimitItem[];

  /** 环境变量 */
  envItems?: ContainerEnv.ItemWithId[];

  /** 是否开启高级设置 */
  isOpenAdvancedSetting?: boolean;

  /** 高级设置是否校验错误 */
  isAdvancedError?: boolean;

  /** gpu的个数限制 */
  gpu?: number;

  /**gpuManager */
  gpuCore?: string;
  v_gpuCore?: Validation;

  gpuMem?: string;
  v_gpuMem?: Validation;

  /** 工作目录 */
  workingDir?: string;
  v_workingDir?: Validation;

  /** 日志目录 */
  logDir?: string;
  logPath?: string;
  v_logDir?: Validation;

  /** cmd */
  cmd?: string;
  v_cmd?: Validation;

  /** 运行参数 */
  arg?: string;
  v_arg?: Validation;

  /** 健康检查 */
  healthCheck?: HealthCheck;

  /** 是否特权级容器 */
  privileged?: boolean;

  /** 增加权限集 */
  addCapabilities?: string[];

  /** 删除权限集 */
  dropCapabilities?: string[];

  /** 镜像更新策略 */
  imagePullPolicy?: string;
}

/** cpu、mem 的一些限定值 */
export interface LimitItem extends Identifiable {
  /**类型 */
  type?: string;

  /**限制值 */
  value?: string;
  v_value?: Validation;
}

/** Mounitem */
export interface MountItem extends Identifiable {
  /** 数据卷 */
  volume?: string;
  v_volume?: Validation;

  /** 目标路径 */
  mountPath?: string;
  v_mountPath?: Validation;

  /** 子路径 */
  mountSubPath?: string;
  v_mountSubPath?: Validation;

  /** 权限 */
  mode?: string;
  v_mode?: Validation;
}

/** 健康检查项 */
export interface HealthCheck {
  /**是否开启存活检查 */
  isOpenLiveCheck?: boolean;

  /**存活检查参数 */
  liveCheck?: HealthCheckItem;

  /**是否开启就绪检查 */
  isOpenReadyCheck?: boolean;

  /**就绪检查参数*/
  readyCheck?: HealthCheckItem;
}

export interface HealthCheckItem {
  /**检查方法 */
  checkMethod?: string;

  /**检查端口 */
  port?: string;
  v_port?: Validation;

  /**检查协议 */
  protocol?: string;

  /**检查路径 */
  path?: string;
  v_path?: Validation;

  /**命令 */
  cmd?: string;
  v_cmd?: Validation;

  /**启动延时 */
  delayTime?: number;
  v_delayTime?: Validation;

  /**响应超时 */
  timeOut?: number;
  v_timeOut?: Validation;

  /**间隔时间 */
  intervalTime?: number;
  v_intervalTime?: Validation;

  /**健康阈值 */
  healthThreshold?: number;
  v_healthThreshold?: Validation;

  /**不健康阈值 */
  unhealthThreshold?: number;
  v_unhealthThreshold?: Validation;
}
