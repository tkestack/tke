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
import { Router } from '../../../helpers/Router';

/**
 * @param sub   二级菜单栏的一级导航
 * @param mode  当前的展示内容类型 list | create | update | detail
 * @param type  三级菜单栏所对应的资源 resource | service ……
 * @param resourceName  资源的名称  deployment 等
 * @param tab   tab页面
 */
export const router = new Router('/tkestack/helm(/:sub)(/:tab)', {
  sub: '',
  tab: ''
});
