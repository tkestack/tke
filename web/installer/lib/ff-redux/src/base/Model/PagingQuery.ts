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
/** 表示一次分页请求 */
export interface PagingQuery {
  /** 请求的页码，从 1 开始索引 */
  pageIndex?: number;

  /** 请求的每页记录数 */
  pageSize?: number;

  append?: boolean;

  /** 翻页清空pages数组 */
  clear?: boolean;
}
