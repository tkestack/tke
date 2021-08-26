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

export interface Validation {
  /**验证状态 0: 初始状态；1：校验通过；2：校验不通过；*/
  status?: number;

  /**结果描述 */
  message?: string | React.ReactNode;

  /**
   * 返回的校验列表
   * 目前仅 CIDR 有使用
   */
  list?: any[];
}

export const initValidator = {
  /**验证状态 0: 初始状态；1：校验通过；2：校验不通过；*/
  status: 0,

  /**结果描述 */
  message: ''
};
