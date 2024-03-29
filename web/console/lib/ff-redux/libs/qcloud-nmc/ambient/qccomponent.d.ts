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
  export interface BeeInstance {
    $set(props: any): void;
    $set(path: string, props: any): void;
    $get(path: string): void;
    $destroy(): void;
  }

  export interface Bee {
    mount(id: string | HTMLElement, props?: any): BeeInstance;
    utils: {
      ie: number;
    };
  }
  export namespace Bee {}
}
