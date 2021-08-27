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
export enum FFListType {
  Base = 'Base',
  Fetch = 'Fetch',
  Query = 'Query',
  InitValue = 'InitValue',
  Selection = 'Selection',
  InitValues = 'InitValues',
  Selections = 'Selections',
  Clear = 'Clear',
  ClearData = 'ClearData',
  ClearSelection = 'ClearSelection',
  DisplayField = 'DisplayField',
  ValueField = 'ValueField',
  GroupKeyField = 'GroupKeyField'
}
export function getFFListActionType(actionName: string, type: FFListType, id?: string) {
  let ts = id ? [actionName, type, id] : [actionName, type];
  return ts.join('_');
}

export function createFFListActionType(actionName: string, id?: string) {
  return {
    Base: getFFListActionType(actionName, FFListType.Base, id),
    Fetch: getFFListActionType(actionName, FFListType.Fetch, id),
    Query: getFFListActionType(actionName, FFListType.Query, id),
    InitValue: getFFListActionType(actionName, FFListType.InitValue, id),
    Selection: getFFListActionType(actionName, FFListType.Selection, id),
    InitValues: getFFListActionType(actionName, FFListType.InitValues, id),
    Selections: getFFListActionType(actionName, FFListType.Selections, id),
    Clear: getFFListActionType(actionName, FFListType.Clear, id),
    ClearSelection: getFFListActionType(actionName, FFListType.ClearSelection, id),
    DisplayField: getFFListActionType(actionName, FFListType.DisplayField, id),
    ValueField: getFFListActionType(actionName, FFListType.ValueField, id),
    GroupKeyField: getFFListActionType(actionName, FFListType.GroupKeyField, id)
  };
}
