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
