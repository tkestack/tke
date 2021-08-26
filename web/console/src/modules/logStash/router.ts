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

import { Router } from '../../../helpers/Router';

/**
 * @param sub 二级导航，eg: create、update等
 * @param tab 详情等tab子页
 * 业务侧和平台侧使用不同的URL 路径
 */
const baseURI = window.location.href.includes('/tkestack-project') ? '/tkestack-project/log' : '/tkestack/log';
export const router = new Router(`${baseURI}(/:mode)(/:tab)`, { mode: '', tab: '' });
