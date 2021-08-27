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
declare namespace nmc {
  interface Tips {
    /**
     * 显示成功提醒
     */
    success(message: string, duration?: number): void;

    /**
     * 显示错误提醒
     */
    error(message: string, duration?: number): void;

    /**
     * 显示加载指示器
     */
    showLoading(loadingText?: string): void;

    /**
     * 停止加载指示器
     */
    stopLoading(): void;
  }
}
