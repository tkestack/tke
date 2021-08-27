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
export enum FFObjectType {
  Base = 'Base',
  Fetch = 'Fetch',
  Query = 'Query',
  Clear = 'Clear'
}
export function getFFObjectActionType(actionName: string, type: FFObjectType, id?: string) {
  let ts = id ? [actionName, type, id] : [actionName, type];
  return ts.join('_');
}

export function createFFObjectActionType(actionName: string, id?: string) {
  return {
    Base: getFFObjectActionType(actionName, FFObjectType.Fetch, id),
    Fetch: getFFObjectActionType(actionName, FFObjectType.Fetch, id),
    Query: getFFObjectActionType(actionName, FFObjectType.Query, id),
    Clear: getFFObjectActionType(actionName, FFObjectType.Clear, id)
  };
}
