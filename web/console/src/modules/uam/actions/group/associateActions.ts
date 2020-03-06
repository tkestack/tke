import { extend, createFFListActions, isSuccessWorkflow, generateWorkflowActionCreator, OperationTrigger  } from '@tencent/ff-redux';
import { RootState, GroupAssociation, GroupPlain, GroupFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import {
  initGroupAssociationState
} from '../../constants/initState';

type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchGroupActions = createFFListActions<GroupPlain, GroupFilter>({
  actionName: ActionTypes.GroupPlainList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchGroupPlainList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().groupPlainList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {

  }
});

/**
 * 列表操作
 */
const fetchGroupAssociatedActions = createFFListActions<GroupPlain, GroupFilter>({
  actionName: ActionTypes.GroupAssociatedList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchGroupAssociatedList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().groupAssociatedList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    // 拉取关联列表之后，更新关联用户组面板WorkflowDialog会用到的GroupAssociation状态数据
    dispatch({
      type: ActionTypes.UpdateGroupAssociation,
      payload: Object.assign({}, getState().groupAssociation, {
        groups: record.data.records,
        originGroups: record.data.records,
        addGroups: [],
        removeGroups: [],
      })
    });
  }
});

/**
 * 关联操作
 */
const associateGroupWorkflow = generateWorkflowActionCreator<GroupAssociation, GroupFilter>({
  actionType: ActionTypes.AssociateGroup,
  workflowStateLocator: (state: RootState) => state.associateGroupWorkflow,
  operationExecutor: WebAPI.associateGroup,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      let { associateGroupWorkflow, groupAssociation, groupFilter } = getState();
      if (isSuccessWorkflow(associateGroupWorkflow)) {
        /** 回调函数 */
        groupFilter.callback && groupFilter.callback();
        /** 开始解绑工作流 */
        if (groupAssociation.removeGroups.length > 0) {
          dispatch(associateActions.disassociateGroupWorkflow.start([groupAssociation], groupFilter));
          dispatch(associateActions.disassociateGroupWorkflow.perform());
        } else {
          /** 重新加载关联数据 */
          dispatch(associateActions.groupAssociatedList.applyFilter(groupFilter));
        }
      }
      /** 结束工作流 */
      dispatch(associateActions.associateGroupWorkflow.reset());
    }
  }
});

/**
 * 解绑操作
 */
const disassociateGroupWorkflow = generateWorkflowActionCreator<GroupAssociation, GroupFilter>({
  actionType: ActionTypes.DisassociateGroup,
  workflowStateLocator: (state: RootState) => state.disassociateGroupWorkflow,
  operationExecutor: WebAPI.disassociateGroup,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      let { disassociateGroupWorkflow, groupFilter } = getState();
      if (isSuccessWorkflow(disassociateGroupWorkflow)) {
        /** 回调函数 */
        groupFilter.callback && groupFilter.callback();
        /** 解绑之后重新加载关联数据 */
        dispatch(associateActions.groupAssociatedList.applyFilter(groupFilter));
      }
      /** 结束工作流 */
      dispatch(associateActions.disassociateGroupWorkflow.reset());
    }
  }
});

const restActions = {
  associateGroupWorkflow,
  disassociateGroupWorkflow,

  /** 设置用户组过滤器 */
  setupGroupFilter: (filter: GroupFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionTypes.UpdateGroupFilter,
        payload: Object.assign({}, getState().groupFilter, filter),
      });
    };
  },

  /** 选中用户组，根据原始数据计算将添加的用户组和将删除的用户组 */
  selectGroup: (groups: GroupPlain[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      /** 选中关联用户组，则更新关联用户组面板WorkflowDialog会用到的GroupAssociation状态数据 */
      /** 比对计算出新增和删除的用户组，originGroups是指原先绑定的用户组 */
      const { originGroups } = getState().groupAssociation;
      const getDifferenceSet = (arr1, arr2) => {
        let a1 = arr1.map(JSON.stringify);
        let a2 = arr2.map(JSON.stringify);
        return a1.concat(a2).filter(v => !a1.includes(v) || !a2.includes(v)).map(JSON.parse);
      };
      let allGroups = groups.concat(originGroups);
      let removeGroups = getDifferenceSet(groups, allGroups);
      let addGroups = getDifferenceSet(originGroups, allGroups);
      dispatch({
        type: ActionTypes.UpdateGroupAssociation,
        payload: Object.assign({}, getState().groupAssociation, {
          groups: groups,
          addGroups: addGroups,
          removeGroups: removeGroups
        })
      });
    };
  },

  /** 清除用户组关联状态数据 */
  clearGroupAssociation: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionTypes.UpdateGroupAssociation,
        payload: initGroupAssociationState
      });
    };
  }
};

export const associateActions = extend({},
  {
    groupList: fetchGroupActions,
    groupAssociatedList: fetchGroupAssociatedActions
  },
  restActions);
