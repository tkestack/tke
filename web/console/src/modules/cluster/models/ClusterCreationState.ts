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
import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface ClusterCreationState extends Identifiable {
  /**链接集群名字 */
  name?: string;
  v_name?: Validation;

  /**apiServer地址 */
  apiServer?: string;
  v_apiServer?: Validation;

  /**证书 */
  certFile?: string;
  v_certFile?: Validation;

  token?: string;
  v_token?: Validation;

  jsonData?: any;

  currentStep?: number;

  clientCert?: string;
  clientKey?: string;

  username?: string;
  as?: string;

  clusternetCertificate?: string;
  clusternetPrivatekey?: string;
}
