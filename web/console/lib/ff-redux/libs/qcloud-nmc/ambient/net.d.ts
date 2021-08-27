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
/// <reference path="manager.d.ts" />
/// <reference path="constants.d.ts" />
declare namespace nmc {
  interface NetResponse {
    code: number;
    data: any;
  }

  interface NetSendConfig {
    /**
     * GET / POST
     * */
    method?: string;

    /**
     * URL
     * */
    url?: string;
  }

  interface NetSendOption {
    /**
     * Data to send
     * */
    data?: {
      [key: string]: any;
    };

    /**
     * Callback to receive data
     * */
    cb?: (response: NetResponse) => any;

    /**
     * if true, a loading will be display on top
     * */
    global?: boolean;
  }

  interface Net {
    send(config: NetSendConfig, option?: NetSendOption): void;
  }
}
