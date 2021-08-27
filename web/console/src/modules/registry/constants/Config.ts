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
import { initValidator } from '../../common/models/Validation';

export const REPO_URL = '/apis/registry.tkestack.io/v1/namespaces/';

export const CHART_URL = '/apis/registry.tkestack.io/v1/chartgroups/';

export const InitApiKey = {
  description: '',
  expire: 1,
  v_expire: initValidator,
  unit: 'h'
};

export const InitRepo = {
  displayName: '',
  name: '',
  v_name: initValidator,
  visibility: 'Public'
};

export const InitChart = {
  displayName: '',
  name: '',
  v_name: initValidator,
  visibility: 'Public'
};

export const InitImage = {
  displayName: '',
  name: '',
  v_name: initValidator,
  visibility: 'Public'
};

export const Default_D_URL = 'registry.tkestack.com';
