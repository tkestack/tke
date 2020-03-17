import { extend, generateWorkflowActionCreator, OperationTrigger, isSuccessWorkflow, createFFListActions } from '@tencent/ff-redux';
import { RootState, PolicyAssociation, PolicyPlain, PolicyFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import {
  initPolicyAssociationState
} from '../../constants/initState';

type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchPolicyActions = createFFListActions<PolicyPlain, PolicyFilter>({
  actionName: ActionTypes.PolicyPlainList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchPolicyPlainList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().policyPlainList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {

  }
});

/**
 * 列表操作
 */
const fetchPolicyAssociatedActions = createFFListActions<PolicyPlain, PolicyFilter>({
  actionName: ActionTypes.PolicyAssociatedList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchPolicyAssociatedList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().policyAssociatedList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    // 拉取关联列表之后，更新关联策略面板WorkflowDialog会用到的PolicyAssociation状态数据
    dispatch({
      type: ActionTypes.UpdatePolicyAssociation,
      payload: Object.assign({}, getState().policyAssociation, {
        policies: record.data.records,
        originPolicies: record.data.records,
        addPolicies: [],
        removePolicies: [],
      })
    });
  }
});

/**
 * 关联操作
 */
const associatePolicyWorkflow = generateWorkflowActionCreator<PolicyAssociation, PolicyFilter>({
  actionType: ActionTypes.AssociatePolicy,
  workflowStateLocator: (state: RootState) => state.associatePolicyWorkflow,
  operationExecutor: WebAPI.associatePolicy,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      let { associatePolicyWorkflow, policyAssociation, policyFilter } = getState();
      if (isSuccessWorkflow(associatePolicyWorkflow)) {
        /** 回调函数 */
        policyFilter.callback && policyFilter.callback();
        /** 开始解绑工作流 */
        if (policyAssociation.removePolicies.length > 0) {
          dispatch(associateActions.disassociatePolicyWorkflow.start([policyAssociation], policyFilter));
          dispatch(associateActions.disassociatePolicyWorkflow.perform());
        } else {
          /** 重新加载关联数据 */
          dispatch(associateActions.policyAssociatedList.applyFilter(policyFilter));
        }
      }
      /** 结束工作流 */
      dispatch(associateActions.associatePolicyWorkflow.reset());
    }
  }
});

/**
 * 解绑操作
 */
const disassociatePolicyWorkflow = generateWorkflowActionCreator<PolicyAssociation, PolicyFilter>({
  actionType: ActionTypes.DisassociatePolicy,
  workflowStateLocator: (state: RootState) => state.disassociatePolicyWorkflow,
  operationExecutor: WebAPI.disassociatePolicy,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      /** 解绑之后重新加载关联数据 */
      let { disassociatePolicyWorkflow, policyFilter } = getState();
      if (isSuccessWorkflow(disassociatePolicyWorkflow)) {
        /** 回调函数 */
        policyFilter.callback && policyFilter.callback();
        dispatch(associateActions.policyAssociatedList.applyFilter(policyFilter));
      }
      /** 结束工作流 */
      dispatch(associateActions.disassociatePolicyWorkflow.reset());
    }
  }
});

const restActions = {
  associatePolicyWorkflow,
  disassociatePolicyWorkflow,

  /** 设置策略过滤器 */
  setupPolicyFilter: (filter: PolicyFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionTypes.UpdatePolicyFilter,
        payload: Object.assign({}, getState().policyFilter, filter),
      });
    };
  },

  /** 选中策略，根据原始数据计算将添加的策略和将删除的策略 */
  selectPolicy: (policies: PolicyPlain[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      /** 选中关联策略，则更新关联策略面板WorkflowDialog会用到的PolicyAssociation状态数据 */
      /** 比对计算出新增和删除的策略，originPolicies是指原先绑定的策略 */
      const { originPolicies } = getState().policyAssociation;
      const getDifferenceSet = (arr1, arr2) => {
        let a1 = arr1.map(JSON.stringify);
        let a2 = arr2.map(JSON.stringify);
        return a1.concat(a2).filter(v => !a1.includes(v) || !a2.includes(v)).map(JSON.parse);
      };
      let allPolicies = policies.concat(originPolicies);
      let removePolicies = getDifferenceSet(policies, allPolicies);
      let addPolicies = getDifferenceSet(originPolicies, allPolicies);
      dispatch({
        type: ActionTypes.UpdatePolicyAssociation,
        payload: Object.assign({}, getState().policyAssociation, {
          policies: policies,
          addPolicies: addPolicies,
          removePolicies: removePolicies
        })
      });
    };
  },

  /** 清除策略关联状态数据 */
  clearPolicyAssociation: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionTypes.UpdatePolicyAssociation,
        payload: initPolicyAssociationState
      });
    };
  }
};

export const associateActions = extend({},
  {
    policyList: fetchPolicyActions,
    policyAssociatedList: fetchPolicyAssociatedActions
  },
  restActions);
