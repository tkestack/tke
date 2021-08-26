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

/**
 * @author techirdliu@tencent.com
 *
 * A workflow abstructor utility for redux operation, with async supported
 *
 */
import { Dispatch } from 'redux';

import { Identifiable, ReduxAction } from '../../../';

export type OperationResult<TTarget> = {
  target: TTarget;
  success: boolean;
  error?: any;
};

export enum OperationState {
  /** indicates an operation is not started yet */
  Pending = 'Pending' as any,

  /** indicates an operation is started, under user interactions */
  Started = 'Started' as any,

  /** indicates the operation is performing after an user action */
  Performing = 'Performing' as any,

  /** indicates the operation is done (success or failed info is in operation result) */
  Done = 'Done' as any
}

export enum OperationTrigger {
  /** `Pending` =====(target)=====> `Started` */
  Start = 'Start' as any,

  /** `Started` =====(params)=====> `Performing` */
  Perform = 'Perform' as any,

  /** `Started` =======> `Pending` */
  Cancel = 'Cancel' as any,

  /** `Performing` =======> `Done` */
  Done = 'Done' as any,

  /** `Done` =======> `Pending` */
  Reset = 'Reset' as any
}

export type WorkflowState<TTarget extends Identifiable, TParam> = {
  /** current operation state */
  operationState: OperationState;

  /** target is specific when started */
  targets?: TTarget[];

  /** params will be stored after perform */
  params?: TParam;

  /** result is specific when done */
  results?: OperationResult<TTarget>[];
};

/** a workflow trigger action */
export type WorkflowPayload<TTarget extends Identifiable, TParam> = {
  /** trigger the work flow */
  trigger: OperationTrigger;

  /** specific operation targers when start, use array for batch operation */
  targets?: TTarget[];

  /** params need during perform trigger */
  params?: TParam;

  /** tell the result when done or failed, will be array for batch operation */
  results?: OperationResult<TTarget>[];
};

export type WorkflowAction<TTarget extends Identifiable, TParam> = ReduxAction<WorkflowPayload<TTarget, TParam>>;

export function generateWorkflowReducer<TTarget extends Identifiable, TParam>({
  actionType
}: {
  actionType: string | number;
}) {
  /** reducer for handling workflow states, takes workflow trigger for the action payload */
  function WorkFlowReducer(
    state: WorkflowState<TTarget, TParam> = {
      operationState: OperationState.Pending
    },
    action: WorkflowAction<TTarget, TParam>
  ): WorkflowState<TTarget, TParam> {
    if (action.type.toString().indexOf(actionType as string) !== 0) {
      return state;
    }

    switch (action.payload.trigger) {
      case OperationTrigger.Start:
        // 执行过程中不允许设置目标
        if (state.operationState !== OperationState.Performing) {
          return {
            operationState: OperationState.Started,
            targets: action.payload.targets,
            params: action.payload.params
          };
        }
        break;

      case OperationTrigger.Perform:
        if (state.operationState === OperationState.Started || state.operationState === OperationState.Done) {
          return {
            operationState: OperationState.Performing,
            targets: state.targets,
            params: action.payload.params
          };
        }
        break;

      case OperationTrigger.Done:
        if (state.operationState === OperationState.Performing) {
          return {
            operationState: OperationState.Done,
            targets: state.targets,
            results: action.payload.results
          };
        }
        break;

      case OperationTrigger.Cancel:
        if (state.operationState === OperationState.Started) {
          return {
            operationState: OperationState.Pending
          };
        }
        break;

      case OperationTrigger.Reset:
        if (state.operationState === OperationState.Done) {
          return {
            operationState: OperationState.Pending
          };
        }
        break;
    }
    return state;
  }
  return WorkFlowReducer;
}

export type WorkflowActionCreator<TTarget extends Identifiable, TParam> = {
  start(targets?: TTarget[], operationParams?: TParam);
  perform(operationParams?: TParam);
  cancel();
  reset();
};

export type OperationHooks = {
  [trigger: number]: (dispatch: Dispatch, getState: () => any) => void;
};

export function generateWorkflowActionCreator<TTarget extends Identifiable, TParam>({
  actionType,
  workflowStateLocator,
  operationExecutor,
  before,
  after
}: {
  actionType: number | string;
  workflowStateLocator: (state: any) => WorkflowState<TTarget, TParam>;
  operationExecutor: (
    targets: TTarget[],
    params: TParam,
    dispatch: Dispatch,
    getState: () => any
  ) => Promise<OperationResult<TTarget>[]>;
  before?: OperationHooks;
  after?: OperationHooks;
}) {
  type ActionType = WorkflowAction<TTarget, TParam>;

  function beforeHook(trigger: OperationTrigger, dispatch: Dispatch, getState: () => any) {
    if (before && typeof before[trigger] === 'function') {
      before[trigger](dispatch, getState);
    }
  }

  function afterHook(trigger: OperationTrigger, dispatch: Dispatch, getState: () => any) {
    if (after && typeof after[trigger] === 'function') {
      after[trigger](dispatch, getState);
    }
  }

  function dispatchWithHook(dispatch: Dispatch, getState: () => any, action: WorkflowAction<TTarget, TParam>) {
    beforeHook(action.payload.trigger, dispatch, getState);
    dispatch(action);
    afterHook(action.payload.trigger, dispatch, getState);
  }

  // create a start action for the operation
  function start(targets: TTarget[], operationParams?: TParam) {
    return (dispatch: Dispatch, getState: any) => {
      const startAction = {
        type: actionType + (OperationTrigger.Start as any),
        payload: {
          trigger: OperationTrigger.Start,
          targets: targets,
          params: operationParams
        }
      };
      dispatchWithHook(dispatch, getState, startAction);
    };
  }

  // create a perform action for the operation
  function perform(operationParams?: TParam) {
    // perform will return a thunk for async support
    return (dispatch: Dispatch, getState: () => any) => {
      const performAction: ActionType = {
        type: actionType + (OperationTrigger.Perform as any),
        payload: {
          trigger: OperationTrigger.Perform,
          params: operationParams || workflowStateLocator(getState()).params
        }
      };
      dispatchWithHook(dispatch, getState, performAction);

      const { params, targets } = workflowStateLocator(getState());

      const resultWithError = (target: TTarget, error: any) => ({
        target,
        success: false,
        error
      });

      operationExecutor(targets, params, dispatch, getState).then(
        results => dispatch(done(results)),
        error => dispatch(done(targets.map(target => resultWithError(target, error))))
      );
    };
  }

  // (private) create a done action for the operation
  function done(results: OperationResult<TTarget>[]) {
    return (dispatch: Dispatch, getState: () => any) => {
      const doneAction = {
        type: actionType + (OperationTrigger.Done as any),
        payload: {
          trigger: OperationTrigger.Done,
          results: results
        }
      };
      dispatchWithHook(dispatch, getState, doneAction);
    };
  }

  // create a cancel action for the operation
  function cancel() {
    return (dispatch: Dispatch, getState: () => any) => {
      const cancelAction = {
        type: actionType + (OperationTrigger.Cancel as any),
        payload: {
          trigger: OperationTrigger.Cancel
        }
      };
      dispatchWithHook(dispatch, getState, cancelAction);
    };
  }

  // create a reset action for the operation
  function reset() {
    return (dispatch: Dispatch, getState: () => any) => {
      const resetAction = {
        type: actionType + (OperationTrigger.Reset as any),
        payload: {
          trigger: OperationTrigger.Reset
        }
      };
      dispatchWithHook(dispatch, getState, resetAction);
    };
  }

  return { start, perform, cancel, reset } as WorkflowActionCreator<TTarget, TParam>;
}

/**
 * generate a reducer for specific workflow triggers
 */
export function reducerForWorkflowTrigger<TState>({
  actionType,
  triggers,
  reducer
}: {
  actionType?: number | string;
  triggers: OperationTrigger[];
  reducer: (state: TState, action: WorkflowAction<any, any>) => TState;
}) {
  return (state: TState, action: WorkflowAction<any, any>) => {
    if (typeof actionType !== 'undefined' && action.type !== actionType) {
      return state;
    }
    if (triggers.indexOf(action.payload.trigger) === -1) {
      return state;
    }
    return reducer(state, action);
  };
}

export interface WorkflowProps<TTarget extends Identifiable, TParam> {
  actions: WorkflowActionCreator<TTarget, TParam>;
  workflow: WorkflowState<TTarget, TParam>;
}

export function getWorkflowStatistics(state: WorkflowState<any, any>) {
  if (state.operationState === OperationState.Done) {
    const total = state.results.length;
    const success = state.results.filter(x => x.success).length;
    const failed = total - success;
    return { total, success, failed };
  }
  throw new RangeError('can only get statistics when a workflow is done.');
}

export function isSuccessWorkflow(state: WorkflowState<any, any>) {
  const { total, success } = getWorkflowStatistics(state);
  return total === success;
}

export function makeBatchExecutor<TTarget, TParam>(
  singleExecutor: (target: TTarget, params?: TParam) => Promise<OperationResult<TTarget>>
): (targets: TTarget[], params?: TParam) => Promise<OperationResult<TTarget>[]> {
  return (targets, params) => singleExecutor(targets[0], params).then(result => [result]);
}
