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
