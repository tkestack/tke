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

import { Identifiable } from '@tencent/ff-redux';

export interface Repository extends Identifiable {
  /**镜像名称 */
  reponame?: string;

  /**镜像类型 */
  repotype?: string;

  /**tag数量 */
  tagCount?: number;

  /**是否公有 */
  public?: number;

  /**是否被用户收藏 */
  isUserFavor?: boolean;

  /**是否为Qcloud官方镜像 */
  isQcloudOfficial?: boolean;

  /**收藏数 */
  favorCount?: number;

  /**下载数 */
  pullCount?: number;

  /**描述*/
  description?: string;

  /**仓库地址 */
  address?: string;

  /**创建时间 */
  creationTime?: string;

  /**logo地址 */
  logo?: string;

  /**镜像描述（仅官方镜像有）*/
  simpleDesc?: string;

  /**地域 */
  regionId?: string | number;
}

export interface RepositoryFilter {
  /**镜像类别 */
  repotype?: string;

  /**镜像名称 */
  reponame?: string;

  /**是否公有 */
  public?: number;

  /**命名空间 */
  namespace?: string;

  /**地域 */
  regionId?: string | number;
}
