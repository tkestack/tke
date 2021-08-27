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

export interface Cls extends Identifiable {
  /** 创建时间 */
  create_time?: string;

  /** 日志集id */
  logset_id?: string;

  /** 日志集名称 */
  logset_name?: string;

  /** 日志集保存时间: unit（天） */
  period?: number;
}

export interface ClsFilter {
  /** 当前的地域 */
  regionId?: number;

  /** 是否能够拉取cls列表 */
  isCanFetchClsList?: boolean;
}

export interface ClsTopic extends Identifiable {
  /**创建时间 */
  create_time?: string;

  /**collection */
  collection?: boolean;

  /**额外的规则 */
  extract_rule?: any;

  index?: boolean;

  log_type?: string;

  /**相对应的日志集的id */
  logset_id?: string;

  /**机器组 */
  machine_group?: any;

  /**日志路径 */
  path?: string;

  /** topicId */
  topic_id?: string;

  /** topicName */
  topic_name?: string;
}

export interface ClsTopicFilter {
  /** 当前的logsetId */
  logsetId?: string;

  /** 地域的Id */
  regionId?: number;

  /**是否能够拉取ClsTopic */
  isCanFetchClsTopic?: boolean;
}
