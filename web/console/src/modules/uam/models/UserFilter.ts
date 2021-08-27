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
export interface UserFilter {
  /** 用户名(唯一) */
  name?: string;

  /** 展示名 */
  displayName?: string;

  /** 相关参数 */
  search?: string;

  ifAll?: boolean;

  /** 是否只拉取策策略需要绑定的用户 */
  isPolicyUser?: boolean;
}
