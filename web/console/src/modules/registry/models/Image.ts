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

import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface Image extends Identifiable {
  apiVersion?: string;
  kind?: string;
  metadata?: {
    annotations?: any;
    clusterName?: string;
    creationTimestamp?: string;
    generateName?: string;
    generation?: number;
    name?: 'string';
    namespace?: 'string';
    resourceVersion?: 'string';
    selfLink?: 'string';
    uid?: 'string';
  };
  spec?: {
    displayName?: 'string';
    name?: 'string';
    namespaceName?: 'string';
    tenantID?: 'string';
    visibility?: 'Public' | 'Private';
  };
  status?: {
    locked?: boolean;
    pullCount?: number;
    tags?: Tag[];
  };
}

export interface Tag extends Identifiable {
  digest?: string;
  name?: string;
  timeCreated?: string;
}

export interface ImageFilter {
  namespace?: string;
  namespaceName?: string;
}

export interface ImageCreation extends Identifiable {
  // kind: Repository
  // metadata?: {
  namespace?: string;
  // };
  // spec?: {
  displayName?: string;
  name?: string;
  v_name?: Validation;
  namespaceName?: string;
  visibility?: 'Public' | 'Private';
  // };
}
